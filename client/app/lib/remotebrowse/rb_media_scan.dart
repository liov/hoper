import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_file_icon.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';

bool rbIsMediaGridEntry(RbFileEntry entry) {
  if (entry.isDirectory) {
    return false;
  }
  return entry.isMotionPhoto || rbFileSupportsAgentThumb(entry.name);
}

/// 按需拉取媒体项：优先服务端 DFS 分页，旧 Agent 回退客户端逐目录 listFiles。
class RbMediaDfsScanner {
  RbMediaDfsScanner({required this.wire, required this.rootPath});

  final RbGrpcSession wire;
  final String rootPath;

  var _cancelled = false;
  var _done = false;
  var _useClientDfs = false;
  var _serverTried = false;
  MediaListCursor _cursor = MediaListCursor();
  final _clientStack = <String>[];
  final _clientPending = <RbMediaItem>[];

  bool get isExhausted => _done && _clientPending.isEmpty;

  void cancel() => _cancelled = true;

  /// 拉下一批，最多 [count] 项；无更多时返回空且 [isExhausted] 为 true。
  Future<List<RbMediaItem>> fetchItems({int count = 48}) async {
    if (_cancelled || !wire.isOpen || isExhausted) {
      return const [];
    }
    final limit = count.clamp(1, 128);
    if (!_useClientDfs) {
      try {
        return await _fetchServer(limit);
      } catch (e, st) {
        if (!_serverTried) {
          rbLog.info('listMedia unavailable, fallback client dfs: $e');
          rbLog.fine('listMedia fallback', e, st);
        }
        _useClientDfs = true;
        _done = false;
      }
    }
    return _fetchClient(limit);
  }

  Future<List<RbMediaItem>> _fetchServer(int limit) async {
    _serverTried = true;
    final root = rbNormRemotePath(rootPath);
    for (var round = 0; round < 64; round++) {
      if (_cancelled || !wire.isOpen || isExhausted) {
        return const [];
      }
      final page = await wire.listMedia(root, cursor: _cursor, pageSize: limit);
      _cursor = page.nextCursor;
      _done = page.done || (page.items.isEmpty && rbMediaCursorEmpty(page.nextCursor));
      if (page.items.isNotEmpty) {
        return page.items;
      }
      if (_done) {
        return const [];
      }
    }
    return const [];
  }

  Future<List<RbMediaItem>> _fetchClient(int limit) async {
    final root = rbNormRemotePath(rootPath);
    if (_clientStack.isEmpty && _clientPending.isEmpty && !_done) {
      _clientStack.add(root);
    }
    while (_clientPending.length < limit && _clientStack.isNotEmpty && !_cancelled && wire.isOpen) {
      final dir = _clientStack.removeLast();
      List<RbFileEntry> entries;
      String resolved;
      try {
        final list = await wire.listFiles(dir);
        resolved = list.resolvedRootPath.isNotEmpty ? rbNormResolvedRootPath(list.resolvedRootPath) : dir;
        entries = list.entries;
      } catch (_) {
        continue;
      }
      final subdirs = entries.where((e) => e.isDirectory).toList()..sort((a, b) => a.name.compareTo(b.name));
      for (var i = subdirs.length - 1; i >= 0; i--) {
        _clientStack.add(rbJoinRemotePath(resolved, subdirs[i].name));
      }
      final media = entries.where(rbIsMediaGridEntry);
      final need = limit - _clientPending.length;
      if (media.length > need + 32) {
        for (final e in media.take(need)) {
          _clientPending.add(RbMediaItem(relPath: rbJoinRemotePath(resolved, e.name), entry: e));
        }
        break;
      }
      final sorted = media.toList();
      rbSortEntriesByMtimeDesc(sorted);
      for (final e in sorted) {
        _clientPending.add(RbMediaItem(relPath: rbJoinRemotePath(resolved, e.name), entry: e));
        if (_clientPending.length >= limit) {
          break;
        }
      }
    }
    if (_clientStack.isEmpty && _clientPending.isEmpty) {
      _done = true;
    }
    final n = limit.clamp(0, _clientPending.length);
    if (n == 0) {
      return const [];
    }
    final chunk = _clientPending.sublist(0, n);
    _clientPending.removeRange(0, n);
    return chunk;
  }
}
