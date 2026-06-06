import 'dart:async';
import 'dart:math';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_file_icon.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:app/remotebrowse/rb_thumb_store.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/scheduler.dart';
import 'package:flutter/widgets.dart';

/// 缩略图按需加载：仅预取可见范围（含前后缓冲），不一次拉全目录。
class RbThumbLoader extends ChangeNotifier {
  RbThumbLoader(
    this._grpc, {
    String? hostKey,
    RbThumbStore? store,
    Duration cacheTtl = const Duration(days: 7),
  }) : _store = store ??
            RbThumbStore(
              hostKey: hostKey ?? _grpc.thumbCacheHostKey(),
              cacheTtl: cacheTtl,
            ) {
    unawaited(_store.purgeExpired());
  }

  final RbGrpcSession _grpc;
  final RbThumbStore _store;
  static const batchSize = 32;
  static const prefetchAhead = 10;
  static const prefetchBehind = 4;
  /// 列表与预览底栏共用，磁盘缓存 key 一致。
  static const sharedThumbEdge = 80;

  /// 首屏/换目录预取张数：可见行（列）× 前后缓冲，并按缩略图边长缩放（大图少拉）。
  static int initialPrefetchLimit({
    required double viewportHeight,
    required double mainAxisExtent,
    int crossAxisCount = 1,
    int aheadRows = prefetchAhead,
    int behindRows = prefetchBehind,
    required int thumbMaxEdge,
  }) {
    final extent = mainAxisExtent > 0 ? mainAxisExtent : 48.0;
    final visibleRows = (viewportHeight / extent).ceil().clamp(1, 999);
    final rows = visibleRows + aheadRows + behindRows;
    var count = rows * crossAxisCount.clamp(1, 99);
    final edge = thumbMaxEdge.clamp(48, 320);
    final scale = (sharedThumbEdge / edge).clamp(0.45, 1.35);
    count = (count * scale).round();
    return count.clamp(8, 128);
  }

  /// 单次 drain 最多发起网络批量的路径数，与首屏预取规模挂钩。
  static int drainTodoCap(int prefetchCount) => (prefetchCount * 2).clamp(24, 160);

  final _data = <String, Uint8List>{};
  final _entryByPath = <String, RbFileEntry>{};
  final _inFlight = <String>{};
  var _prefetchGen = 0;
  final _pendingPrefetch = <String, ({RbFileEntry entry, int maxEdge})>{};
  var _prefetchDraining = false;
  var _prefetchDrainAgain = false;
  Timer? _rangeDebounce;
  final _rangeQueue = <String, ({RbFileEntry entry, int maxEdge})>{};
  int _rangeMaxEdge = sharedThumbEdge;
  var _rangeBypassPause = false;
  Timer? _stripDebounce;
  final _stripQueue = <String, ({RbFileEntry entry, int maxEdge})>{};
  Timer? _mountDebounce;
  var _mountBypassPause = false;
  var _notifyScheduled = false;

  static String _key(String path) => rbNormRemotePath(path);

  void _notifyListenersSafe() {
    if (_notifyScheduled) {
      return;
    }
    _notifyScheduled = true;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _notifyScheduled = false;
      notifyListeners();
    });
  }

  /// 异步拉取/读盘缓存后立刻刷新缩略图；仅在 build 阶段延后到下一帧。
  void _notifyThumbDataChanged() {
    if (SchedulerBinding.instance.schedulerPhase == SchedulerPhase.persistentCallbacks) {
      _notifyListenersSafe();
      return;
    }
    notifyListeners();
  }

  Uint8List? bytesFor(String path) => _data[_key(path)] ?? _data[path];

  bool isLoading(String path) => _inFlight.contains(_key(path)) || _inFlight.contains(path);

  void reset() {
    _prefetchGen++;
    _rangeDebounce?.cancel();
    _rangeDebounce = null;
    _rangeQueue.clear();
    _stripDebounce?.cancel();
    _stripDebounce = null;
    _stripQueue.clear();
    _mountDebounce?.cancel();
    _mountDebounce = null;
    _mountBypassPause = false;
    _pendingPrefetch.clear();
    _prefetchDraining = false;
    _prefetchDrainAgain = false;
    _data.clear();
    _entryByPath.clear();
    _inFlight.clear();
    _notifyListenersSafe();
  }

  /// 目录列表刷新后与 Agent 对齐：删掉的清本地；mtime/size/hash 变了清本地再拉。
  void syncWithListing({
    required List<RbFileEntry> previous,
    required List<RbFileEntry> next,
    required String Function(RbFileEntry) relPath,
  }) {
    final oldFiles = previous.where((e) => !e.isDirectory);
    final newFiles = next.where((e) => !e.isDirectory);
    final oldRel = {for (final e in oldFiles) rbNormRemotePath(relPath(e)): e};
    final newRel = {for (final e in newFiles) rbNormRemotePath(relPath(e)): e};
    for (final p in oldRel.keys) {
      if (!newRel.containsKey(p)) {
        evict(p, entry: oldRel[p]);
      }
    }
    for (final p in newRel.keys) {
      final o = oldRel[p];
      if (o == null) {
        continue;
      }
      final n = newRel[p]!;
      if (o.mtimeUnixMs != n.mtimeUnixMs || o.size != n.size || o.thumbHash != n.thumbHash) {
        evict(p, entry: n);
      }
    }
  }

  void evict(String relPath, {RbFileEntry? entry}) {
    final k = _key(relPath);
    final hash = entry?.thumbHash ?? _entryByPath[k]?.thumbHash ?? '';
    _data.remove(k);
    _entryByPath.remove(k);
    _inFlight.remove(k);
    _rangeQueue.remove(k);
    _stripQueue.remove(k);
    unawaited(_store.removeForRelPath(relPath));
    if (hash.isNotEmpty) {
      unawaited(_store.removeForThumbHash(hash));
    }
    _notifyListenersSafe();
  }

  var _paused = false;

  /// 视频全屏播放时暂停列表缩略图批量，避免与 Range 读盘争抢。
  void setPaused(bool paused) {
    if (_paused == paused) {
      return;
    }
    _paused = paused;
    if (!paused) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        notifyListeners();
      });
    }
  }

  /// 可见区变化时调用（列表/宫格滚动）；合并防抖后批量请求。
  void prefetchVisibleRange(List<({String path, RbFileEntry entry})> items, {int maxEdge = sharedThumbEdge, bool bypassPause = false}) {
    if (_paused && !bypassPause) {
      return;
    }
    if (items.isEmpty) {
      return;
    }
    for (final it in items) {
      if (!rbFileSupportsAgentThumb(it.entry.name)) {
        continue;
      }
      final k = _key(it.path);
      _entryByPath[k] = it.entry;
      _rangeQueue[k] = (entry: it.entry, maxEdge: maxEdge);
    }
    _rangeMaxEdge = maxEdge;
    _rangeBypassPause = _rangeBypassPause || bypassPause;
    _rangeDebounce?.cancel();
    _rangeDebounce = Timer(const Duration(milliseconds: 100), () {
      _rangeDebounce = null;
      final pending = _rangeQueue.entries.toList();
      _rangeQueue.clear();
      if (pending.isEmpty) {
        _rangeBypassPause = false;
        return;
      }
      final items = pending.map((e) => (path: e.key, entry: e.value.entry)).toList();
      final bp = _rangeBypassPause;
      _rangeBypassPause = false;
      unawaited(prefetch(items, maxEdge: _rangeMaxEdge, bypassPause: bp));
    });
  }

  /// 预览底栏横向条：独立队列，不受 [setPaused] 影响。
  void prefetchStripVisibleRange(List<({String path, RbFileEntry entry})> items, {int maxEdge = sharedThumbEdge}) {
    if (items.isEmpty) {
      return;
    }
    for (final it in items) {
      if (!rbFileSupportsAgentThumb(it.entry.name)) {
        continue;
      }
      final k = _key(it.path);
      _entryByPath[k] = it.entry;
      _stripQueue[k] = (entry: it.entry, maxEdge: maxEdge);
    }
    _stripDebounce?.cancel();
    _stripDebounce = Timer(const Duration(milliseconds: 80), () {
      _stripDebounce = null;
      final pending = _stripQueue.entries.toList();
      _stripQueue.clear();
      if (pending.isEmpty) {
        return;
      }
      final edge = pending.map((e) => e.value.maxEdge).reduce(max);
      final batch = pending.map((e) => (path: e.key, entry: e.value.entry)).toList();
      unawaited(prefetch(batch, maxEdge: edge, bypassPause: true));
    });
  }

  /// 单格挂载：并入待拉队列并防抖合并，避免每格一次网络批。
  void requestThumb(String path, RbFileEntry entry, {int maxEdge = sharedThumbEdge, bool bypassPause = false}) {
    if (!rbFileSupportsAgentThumb(entry.name)) {
      return;
    }
    final k = _key(path);
    _entryByPath[k] = entry;
    if (_data.containsKey(k) || _inFlight.contains(k)) {
      return;
    }
    _pendingPrefetch[k] = (entry: entry, maxEdge: maxEdge);
    _scheduleMountDrain(bypassPause: bypassPause);
  }

  void _scheduleMountDrain({required bool bypassPause}) {
    _mountBypassPause = _mountBypassPause || bypassPause;
    _mountDebounce?.cancel();
    _mountDebounce = Timer(const Duration(milliseconds: 80), () {
      _mountDebounce = null;
      final bp = _mountBypassPause;
      _mountBypassPause = false;
      unawaited(_drainPrefetch(bypassPause: bp, todoCap: drainTodoCap(_pendingPrefetch.length)));
    });
  }

  /// 底栏单格：并入 [prefetchStripVisibleRange] 批量队列，勿单张 [ensure]。
  void requestStripThumb(String path, RbFileEntry entry, {int maxEdge = sharedThumbEdge}) {
    prefetchStripVisibleRange([(path: path, entry: entry)], maxEdge: maxEdge);
  }

  void _enqueuePrefetch(List<({String path, RbFileEntry entry})> items, {required int maxEdge}) {
    for (final it in items) {
      if (!rbFileSupportsAgentThumb(it.entry.name)) {
        continue;
      }
      final k = _key(it.path);
      _entryByPath[k] = it.entry;
      _pendingPrefetch[k] = (entry: it.entry, maxEdge: maxEdge);
    }
  }

  Future<void> prefetch(
    List<({String path, RbFileEntry entry})> items, {
    int maxEdge = sharedThumbEdge,
    bool bypassPause = false,
    int? todoCap,
  }) async {
    if (_paused && !bypassPause) {
      return;
    }
    if (items.isEmpty) {
      return;
    }
    _enqueuePrefetch(items, maxEdge: maxEdge);
    await _drainPrefetch(bypassPause: bypassPause, todoCap: todoCap ?? drainTodoCap(items.length));
  }

  Future<void> _fetchThumbChunk(List<String> chunk, {required int edge, required int gen}) async {
    if (gen != _prefetchGen || chunk.isEmpty) {
      return;
    }
    for (final p in chunk) {
      _inFlight.add(_key(p));
    }
    _notifyThumbDataChanged();
    try {
      final map = await _grpc.fetchThumbsBatch(chunk, maxEdge: edge);
      if (gen != _prefetchGen) {
        return;
      }
      for (final p in chunk) {
        final k = _key(p);
        final bytes = map[p] ?? map[k];
        if (bytes == null || bytes.isEmpty) {
          continue;
        }
        _data[k] = bytes;
        final entry = _entryByPath[k];
        if (entry != null) {
          unawaited(_store.put(entry, p, edge, bytes));
        }
      }
      if (map.isEmpty && chunk.isNotEmpty) {
        rbLog.warning('thumb batch empty paths=${chunk.length} sample=${chunk.first}');
      }
    } catch (e, st) {
      rbLog.warning('thumb batch fail', e, st);
    } finally {
      for (final p in chunk) {
        _inFlight.remove(_key(p));
      }
      if (gen == _prefetchGen) {
        _notifyThumbDataChanged();
      }
    }
  }

  Future<void> _drainPrefetch({required bool bypassPause, int todoCap = 72}) async {
    if (_prefetchDraining) {
      _prefetchDrainAgain = true;
      return;
    }
    _prefetchDraining = true;
    final gen = _prefetchGen;
    try {
      while (_pendingPrefetch.isNotEmpty && gen == _prefetchGen) {
        if (_paused && !bypassPause) {
          _pendingPrefetch.clear();
          return;
        }
        final snap = Map<String, ({RbFileEntry entry, int maxEdge})>.from(_pendingPrefetch);
        _pendingPrefetch.clear();
        final edge = snap.values.map((v) => v.maxEdge).reduce(max);
        final cacheHits = await Future.wait(snap.entries.map((e) async {
          final k = e.key;
          if (_data.containsKey(k) || _inFlight.contains(k)) {
            return (k: k, cached: null as Uint8List?, defer: e.value);
          }
          final cached = await _store.load(e.value.entry, k, edge);
          return (k: k, cached: cached, defer: e.value);
        }));
        final todo = <String>[];
        var cacheHitCount = 0;
        for (final hit in cacheHits) {
          final cached = hit.cached;
          if (cached != null && cached.isNotEmpty) {
            _data[hit.k] = cached;
            cacheHitCount++;
            continue;
          }
          if (_data.containsKey(hit.k) || _inFlight.contains(hit.k)) {
            continue;
          }
          if (todo.length >= todoCap) {
            _pendingPrefetch[hit.k] = hit.defer;
            continue;
          }
          todo.add(rbNormRemotePath(hit.k));
        }
        if (cacheHitCount > 0 && gen == _prefetchGen) {
          _notifyThumbDataChanged();
        }
        if (gen != _prefetchGen || todo.isEmpty) {
          continue;
        }
        final chunks = <List<String>>[];
        for (var i = 0; i < todo.length; i += batchSize) {
          chunks.add(todo.sublist(i, min(i + batchSize, todo.length)));
        }
        const maxParallel = 2;
        for (var i = 0; i < chunks.length; i += maxParallel) {
          final batch = chunks.skip(i).take(maxParallel);
          await Future.wait(batch.map((chunk) => _fetchThumbChunk(chunk, edge: edge, gen: gen)));
        }
      }
    } finally {
      _prefetchDraining = false;
      final again = _prefetchDrainAgain || (_pendingPrefetch.isNotEmpty && gen == _prefetchGen);
      _prefetchDrainAgain = false;
      if (again) {
        unawaited(_drainPrefetch(bypassPause: bypassPause, todoCap: drainTodoCap(_pendingPrefetch.length)));
      }
    }
  }

  Future<Uint8List?> ensure(String path, RbFileEntry entry, {int maxEdge = 256}) async {
    if (!rbFileSupportsAgentThumb(entry.name)) {
      return null;
    }
    final k = _key(path);
    _entryByPath[k] = entry;
    final hit = _data[k];
    if (hit != null && hit.isNotEmpty) {
      return hit;
    }
    final cached = await _store.load(entry, path, maxEdge);
    if (cached != null && cached.isNotEmpty) {
      _data[k] = cached;
      _notifyListenersSafe();
      return cached;
    }
    try {
      final map = await _grpc.fetchThumbsBatch([path], maxEdge: maxEdge);
      final bytes = map[path] ?? map[k];
      if (bytes != null && bytes.isNotEmpty) {
        _data[k] = bytes;
        unawaited(_store.put(entry, path, maxEdge, bytes));
        _notifyListenersSafe();
      }
      return bytes;
    } catch (e, st) {
      rbLog.fine('thumb ensure fail: $e', e, st);
      return null;
    }
  }
}
