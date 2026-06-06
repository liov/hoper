import 'dart:async';
import 'dart:math';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/connect_cancel.dart';
import 'package:app/remotebrowse/link_kind.dart';
import 'package:app/remotebrowse/rb_connect_endpoint.dart';
import 'package:app/remotebrowse/rb_connection_profile.dart';
import 'package:app/remotebrowse/rb_file_icon.dart';
import 'package:app/remotebrowse/rb_file_media.dart';
import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:app/remotebrowse/rb_thumb_loader.dart';
import 'package:app/remotebrowse/rb_preview_store.dart';
import 'package:app/remotebrowse/rb_connection_store.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:app/remotebrowse/rb_user_message.dart';
import 'package:app/pages/remotebrowse/rb_grid_tile.dart';
import 'package:app/pages/remotebrowse/rb_media_grid_panel.dart';
import 'package:app/pages/remotebrowse/rb_selectable_grid.dart';
import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/pages/remotebrowse/rb_file_info_panel.dart';
import 'package:app/remotebrowse/rb_file_download.dart';
import 'package:app/pages/remotebrowse/remote_browse_file_preview_page.dart';
import 'package:app/remotebrowse/rb_preview_nav.dart';
import 'package:app/remotebrowse/viewer_session.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/util/nav.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';

class RemoteBrowseFilesView extends StatefulWidget {
  const RemoteBrowseFilesView({super.key, required this.profile});

  final RbConnectionProfile profile;

  @override
  State<RemoteBrowseFilesView> createState() => _RemoteBrowseFilesViewState();
}

enum _BrowseViewMode { list, grid }

class _RemoteBrowseFilesViewState extends State<RemoteBrowseFilesView> {
  late final RbConnectTarget _connectTarget;
  late final String _room;
  /// 起始 list 路径；Windows Agent 首次解析后会更新为 `D:/…` 并持久化。
  late String _listPath;
  List<RbFileEntry> _entries = [];
  _BrowseViewMode _viewMode = _BrowseViewMode.list;
  var _loading = false;
  var _status = '';
  var _errorMsg = '';
  String _infoBanner = '';
  RbLinkKind? _linkKind;
  RbGrpcSession? _grpc;
  RbThumbLoader? _thumbs;
  RbPreviewContentStore? _previewStore;
  RbConnectCancel? _connectCancel;
  String _currentPath = '';
  /// 本次会话远端 Agent 平台（health 探测 / 连接配置），与 rb 启动终端无关。
  String _sessionAgentPlatform = '';

  /// 本次连接首次列目录的 [resolvedRootPath]，再回退到此即回到连接列表。
  String? _browseRootPath;
  int _connectSeq = 0;
  int _reconnectSeq = 0;
  int _listSeq = 0;
  int _listThumbEpoch = 0;
  /// 仅在实际换目录时递增，用作列表/宫格 [State] 的 key（不用路径字符串，避免 resolved 抖动重建滚动）。
  int _browseDirEpoch = 0;
  var _reconnecting = false;

  /// 数据面重连成功时递增，预览页监听后自动刷新 wire 并重试。
  final _wireEpoch = ValueNotifier<int>(0);
  var _showConnectionDetails = false;
  var _selecting = false;
  final _selectedIds = <String>{};
  var _deleting = false;
  Timer? _errorClearTimer;

  @override
  void initState() {
    super.initState();
    final p = widget.profile;
    _connectTarget = RbConnectTarget.fromProfile(p);
    _room = p.room;
    _listPath = rbSanitizeViewerListPath(p.path, agentPlatform: p.agentPlatform);
    _sessionAgentPlatform = p.agentPlatform.trim();
    WidgetsBinding.instance.addPostFrameCallback((_) => _load());
  }

  @override
  void dispose() {
    _connectSeq++;
    _errorClearTimer?.cancel();
    _connectCancel?.cancel();
    _wireEpoch.dispose();
    unawaited(_grpc?.close());
    super.dispose();
  }

  void _bumpWireEpoch() => _wireEpoch.value++;

  void _cancelConnect() => _disconnect(status: '已取消');

  Future<void> _scheduleReconnect(Object cause) async {
    final rid = ++_reconnectSeq;
    if (_grpc == null && !_reconnecting) {
      return;
    }
    final c = _connectCancel;
    _connectCancel = null;
    if (c != null) {
      await c.closeSignal();
    }
    await _grpc?.close();
    _grpc = null;
    _thumbs = null;
    if (!mounted || rid != _reconnectSeq) {
      return;
    }
    setState(() {
      _reconnecting = true;
      _loading = true;
      _errorMsg = '';
      _status = '正在重新连接…';
      _linkKind = null;
    });
    await Future<void>.delayed(const Duration(milliseconds: 800));
    if (!mounted || rid != _reconnectSeq || _grpc != null) {
      return;
    }
    await _load();
  }

  void _showErrorBanner(String msg, {Duration duration = const Duration(seconds: 5)}) {
    final text = msg.trim();
    if (text.isEmpty) {
      return;
    }
    _errorClearTimer?.cancel();
    setState(() => _errorMsg = text);
    _errorClearTimer = Timer(duration, () {
      if (mounted && _errorMsg == text) {
        setState(() => _errorMsg = '');
      }
    });
  }

  /// 未连接：仅在下方主区域展示一次；已连接时断线自动重连。
  void _presentError(Object e) {
    final msg = rbUserMessage(e);
    if (!mounted) {
      return;
    }
    if (_connected && rbIsWireLost(e)) {
      unawaited(_scheduleReconnect(e));
      return;
    }
    setState(() {
      _loading = false;
      if (!_connected) {
        _status = '';
        _linkKind = null;
      }
    });
    _showErrorBanner(msg);
  }

  Future<void> _disconnect({
    String status = '已断开',
    bool keepEntries = false,
  }) async {
    final c = _connectCancel;
    _connectCancel = null;
    if (c != null) {
      await c.closeSignal();
    }
    await _grpc?.close();
    _grpc = null;
    _thumbs = null;
    if (!mounted) {
      return;
    }
    setState(() {
      _loading = false;
      if (!keepEntries) {
        _entries = [];
        _linkKind = null;
        _currentPath = '';
        _browseRootPath = null;
        _reconnecting = false;
      }
      _status = status;
      if (!keepEntries && (status.isEmpty || status == '已断开')) {
        _errorMsg = '';
      }
    });
  }

  /// UI 展示路径：Agent [resolvedRootPath] 已是 Windows `D:/…`，直接展示。
  String _displayLocation() => _activeDirPath;

  String? get _overlayBannerMessage {
    final info = _infoBanner.trim();
    if (info.isNotEmpty) {
      return info;
    }
    if (_reconnecting) {
      return _status.isEmpty ? '正在重新连接…' : _status;
    }
    final err = _errorMsg.trim();
    if (err.isNotEmpty) {
      return err;
    }
    if (_loading) {
      if (_status.isNotEmpty) {
        return _status;
      }
      return _connected ? '加载中…' : '正在连接…';
    }
    return null;
  }

  RbBannerTone get _overlayBannerTone {
    if (_infoBanner.trim().isNotEmpty) {
      return RbBannerTone.info;
    }
    if (_errorMsg.trim().isNotEmpty) {
      return RbBannerTone.warning;
    }
    return RbBannerTone.loading;
  }

  void _showInfoBanner(String msg) {
    final text = msg.trim();
    if (text.isEmpty) {
      return;
    }
    setState(() => _infoBanner = text);
    Future<void>.delayed(const Duration(seconds: 3), () {
      if (mounted && _infoBanner == text) {
        setState(() => _infoBanner = '');
      }
    });
  }

  bool get _isWindowsRemoteAgent =>
      rbParseRemoteOs(_sessionAgentPlatform.isNotEmpty ? _sessionAgentPlatform : widget.profile.agentPlatform) ==
      RbRemoteOs.windows;

  /// 当前目录绝对路径（Agent resolved_root_path，wire 形态）。
  String get _activeDirPath {
    final p = _currentPath.trim();
    if (p.isEmpty) {
      return '';
    }
    return rbNormRemotePath(p);
  }

  /// 按绝对路径分段，每段可点击跳转。
  List<({String label, String path})> _pathSegments() =>
      rbRemotePathSegments(_activeDirPath);

  String _entrySummary() {
    final dirs = _entries.where((e) => e.isDirectory).length;
    final files = _entries.length - dirs;
    if (dirs == 0 && files == 0) {
      return '';
    }
    final parts = <String>[];
    if (dirs > 0) {
      parts.add('$dirs 个文件夹');
    }
    if (files > 0) {
      parts.add('$files 个文件');
    }
    return parts.join(' · ');
  }

  bool get _connected => _grpc != null;

  void _toggleViewMode() {
    if (_selecting) {
      _exitSelection();
    }
    setState(() {
      _viewMode = _viewMode == _BrowseViewMode.list
          ? _BrowseViewMode.grid
          : _BrowseViewMode.list;
    });
  }

  void _exitSelection() {
    setState(() {
      _selecting = false;
      _selectedIds.clear();
    });
  }

  void _enterSelection(RbFileEntry e) {
    if (!_connected) {
      return;
    }
    setState(() {
      _selecting = true;
      _selectedIds
        ..clear()
        ..add(e.id);
    });
  }

  void _toggleSelection(RbFileEntry e) {
    setState(() {
      if (!_selecting) {
        _selecting = true;
        _selectedIds.add(e.id);
        return;
      }
      if (_selectedIds.contains(e.id)) {
        _selectedIds.remove(e.id);
        if (_selectedIds.isEmpty) {
          _selecting = false;
        }
      } else {
        _selectedIds.add(e.id);
      }
    });
  }

  void _applySlideSelection(Set<String> ids) {
    if (!_selecting) {
      return;
    }
    if (_selectedIds.length == ids.length && _selectedIds.containsAll(ids)) {
      return;
    }
    setState(() {
      _selectedIds
        ..clear()
        ..addAll(ids);
      if (_selectedIds.isEmpty) {
        _selecting = false;
      }
    });
  }

  void _onFileTap(RbFileEntry e) {
    if (_selecting) {
      _toggleSelection(e);
      return;
    }
    _openFile(e);
  }

  String _deleteDialogTitle(List<RbFileEntry> targets) {
    final dirs = targets.where((e) => e.isDirectory).length;
    final files = targets.length - dirs;
    if (dirs == 0) {
      return '删除 ${targets.length} 个文件';
    }
    if (files == 0) {
      return targets.length == 1 ? '删除文件夹' : '删除 $dirs 个文件夹';
    }
    return '删除 ${targets.length} 项';
  }

  String _deleteDialogBody(List<RbFileEntry> targets) {
    final dirs = targets.where((e) => e.isDirectory).length;
    final files = targets.length - dirs;
    if (targets.length == 1) {
      final e = targets.first;
      if (e.isDirectory) {
        return '确定删除文件夹「${e.name}」？\n文件夹内全部内容将被删除。';
      }
      return e.name;
    }
    final parts = <String>[];
    if (files > 0) {
      parts.add('$files 个文件');
    }
    if (dirs > 0) {
      parts.add('$dirs 个文件夹（含其内全部内容）');
    }
    return '确定删除 ${parts.join('、')}？';
  }

  Future<void> _deleteSelected() async {
    final wire = _grpc;
    if (wire == null || _selectedIds.isEmpty || _deleting) {
      return;
    }
    final targets = _entries.where((e) => _selectedIds.contains(e.id)).toList();
    if (targets.isEmpty) {
      _exitSelection();
      return;
    }
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(_deleteDialogTitle(targets)),
        content: Text(_deleteDialogBody(targets)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('取消'),
          ),
          FilledButton(
            style: FilledButton.styleFrom(
              backgroundColor: Theme.of(ctx).colorScheme.error,
            ),
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('删除'),
          ),
        ],
      ),
    );
    if (ok != true || !mounted) {
      return;
    }
    setState(() {
      _deleting = true;
      _loading = true;
      _errorMsg = '';
      _infoBanner = '';
      _status = '删除中…';
    });
    var failed = 0;
    final deletedIds = <String>{};
    try {
      const batchSize = 4;
      for (var i = 0; i < targets.length; i += batchSize) {
        final batch = targets.sublist(i, min(i + batchSize, targets.length));
        await Future.wait(batch.map((e) async {
          try {
            await wire.deleteFile(_entryRelPath(e));
            _thumbs?.evict(_entryRelPath(e), entry: e);
            deletedIds.add(e.id);
          } catch (_) {
            failed++;
          }
        }));
      }
      if (!mounted) {
        return;
      }
      setState(() {
        _purgeEntriesByIds(deletedIds);
        _deleting = false;
        _loading = deletedIds.isNotEmpty;
        if (deletedIds.isNotEmpty) {
          _status = '同步列表…';
        }
      });
      if (deletedIds.isNotEmpty) {
        await _refreshList();
      }
      if (!mounted) {
        return;
      }
      if (failed > 0) {
        _showErrorBanner('有 $failed 项删除失败');
      } else if (deletedIds.isNotEmpty) {
        _showInfoBanner('已删除 ${deletedIds.length} 项');
      }
    } catch (e) {
      _presentError(e);
    } finally {
      if (mounted) {
        setState(() {
          _deleting = false;
          if (deletedIds.isEmpty) {
            _loading = false;
          }
        });
      }
    }
  }

  /// 删除成功后先从本地列表摘掉，避免仅依赖刷新仍看到旧项。
  void _purgeEntriesByIds(Set<String> ids) {
    if (ids.isEmpty) {
      return;
    }
    final prev = List<RbFileEntry>.from(_entries);
    _entries = _entries.where((e) => !ids.contains(e.id)).toList();
    _selectedIds.removeWhere(ids.contains);
    if (_selectedIds.isEmpty) {
      _selecting = false;
    }
    final dirs = _entries.where((e) => e.isDirectory).length;
    final files = _entries.length - dirs;
    _status = dirs == 0 ? '$files 个文件' : '$dirs 个目录，$files 个文件';
    _thumbs?.syncWithListing(
      previous: prev,
      next: List<RbFileEntry>.from(_entries),
      relPath: _entryRelPath,
    );
  }

  Future<void> _refreshList() async {
    final wire = _grpc;
    if (wire == null) {
      if (mounted) {
        _showErrorBanner('未连接，请等待连接完成');
      }
      return;
    }
    final seq = ++_listSeq;
    if (mounted) {
      setState(() {
        _loading = true;
        _errorMsg = '';
        _status = '刷新列表…';
      });
    }
    try {
      final list = await wire.listFiles(_apiListPath(_activeDirPath));
      if (!mounted || seq != _listSeq) {
        return;
      }
      setState(() => _applyListResult(list, dirChanged: false));
    } catch (e) {
      if (mounted && seq == _listSeq) {
        _presentError(e);
      }
    } finally {
      if (mounted && seq == _listSeq) {
        setState(() => _loading = false);
      }
    }
  }

  Future<void> _load() async {
    final seq = ++_connectSeq;
    _listSeq++;
    _reconnectSeq++;
    final keepEntries = _reconnecting;
    await _disconnect(status: '', keepEntries: keepEntries);
    if (seq != _connectSeq || !mounted) {
      return;
    }
    final cancel = RbConnectCancel();
    _connectCancel = cancel;
    setState(() {
      _loading = true;
      _status = '连接文件端…';
      _errorMsg = '';
      _linkKind = null;
    });
    try {
      if (_connectTarget.isDirect) {
        _grpc = await RbViewerSession.connectDirect(
          _connectTarget.directHost!,
          _connectTarget.directPort,
          room: _room.trim(),
        );
      } else {
        _grpc = await RbViewerSession.connect(
          _connectTarget.signalUri!,
          _room.trim(),
          cancel: cancel,
        );
      }
      cancel.check();
      if (seq != _connectSeq) {
        throw const RbConnectCancelled();
      }
      await _probeSessionAgentPlatform();
      final list = await _grpc!.listFiles(_initialBrowsePath());
      cancel.check();
      if (seq != _connectSeq) {
        throw const RbConnectCancelled();
      }
      setState(() {
        final hostKey = _grpc!.thumbCacheHostKey(
          fallback: _connectTarget.isDirect
              ? _connectTarget.directStorage
              : (_connectTarget.signalUri?.host ?? widget.profile.id),
        );
        _thumbs = RbThumbLoader(_grpc!, hostKey: hostKey);
        if (!kIsWeb) {
          _previewStore ??= RbPreviewContentStore(hostKey: hostKey);
          unawaited(_previewStore!.purgeExpired());
        }
        _applyListResult(list);
        if (_sessionAgentPlatform.isEmpty && rbIsWindowsRemotePath(_currentPath)) {
          _sessionAgentPlatform = 'windows';
        }
        _linkKind = _grpc!.linkKind;
        _errorMsg = '';
        _reconnecting = false;
      });
      _bumpWireEpoch();
    } on RbConnectCancelled {
      await _grpc?.close();
      _grpc = null;
      if (mounted) {
        setState(() => _status = '已取消');
      }
    } on ArgumentError catch (e) {
      if (seq == _connectSeq) {
        _presentError(e);
      }
      final w = _grpc;
      _grpc = null;
      unawaited(w?.close());
    } catch (e) {
      if (seq == _connectSeq) {
        _presentError(e);
      }
      final w = _grpc;
      _grpc = null;
      unawaited(w?.close());
    } finally {
      if (_connectCancel == cancel) {
        _connectCancel = null;
      }
      if (mounted && seq == _connectSeq) {
        setState(() {
          _loading = false;
          if (_grpc == null) {
            _reconnecting = false;
          }
        });
      }
    }
  }

  String _initialBrowsePath() => _listPath.isEmpty ? '.' : _listPath;

  String _apiListPath(String path) {
    final p = path.trim();
    if (p.isEmpty || p == '.') {
      final active = _activeDirPath;
      if (active.isNotEmpty && active != '/') {
        return active;
      }
      return '.';
    }
    final norm = rbNormRemotePath(p);
    if ((norm == '/' || norm.isEmpty) && (_isWindowsRemoteAgent || rbIsWindowsRemotePath(_browseRootPath ?? ''))) {
      return '.';
    }
    return norm;
  }

  static String _normBrowsePath(String path) => rbNormRemotePath(path);

  bool get _canGoUp {
    final cur = _normBrowsePath(_activeDirPath);
    final root = _browseRootPath;
    if (root != null && root.isNotEmpty) {
      return !rbRemotePathsEqual(cur, root);
    }
    if (cur == '/' || cur.isEmpty) {
      return false;
    }
    return _normBrowsePath(_parentDirPath()) != cur;
  }

  void _applyListResult(RbListFilesResult list, {bool dirChanged = false}) {
    final prev = List<RbFileEntry>.from(_entries);
    final resolved = rbNormResolvedRootPath(list.resolvedRootPath);
    if (_sessionAgentPlatform.isEmpty && rbIsWindowsRemotePath(resolved)) {
      _sessionAgentPlatform = 'windows';
    }
    var dirForPaths = resolved.isNotEmpty ? resolved : _activeDirPath;
    if (dirForPaths.isEmpty && (_browseRootPath ?? '').isNotEmpty) {
      dirForPaths = _browseRootPath!;
    }
    final sorted = List<RbFileEntry>.from(list.entries);
    rbSortEntriesByMtimeDesc(sorted);
    if (dirChanged || prev.isEmpty) {
      _thumbs?.reset();
    } else {
      _thumbs?.syncWithListing(
        previous: prev,
        next: sorted,
        relPath: (e) => rbJoinRemotePath(dirForPaths, e.name),
      );
    }
    final dirs = sorted.where((e) => e.isDirectory).length;
    final files = sorted.length - dirs;
    if (resolved.isNotEmpty) {
      final firstRoot = _browseRootPath == null;
      _currentPath = resolved;
      _browseRootPath ??= resolved;
      _syncListPathAfterResolve(resolved, firstRoot: firstRoot);
    }
    _entries = sorted;
    _listThumbEpoch++;
    _errorMsg = '';
    _status = dirs == 0 ? '$files 个文件' : '$dirs 个目录，$files 个文件';
    if (_selecting) {
      final liveIds = {for (final e in sorted) e.id};
      _selectedIds.removeWhere((id) => !liveIds.contains(id));
      if (_selectedIds.isEmpty) {
        _selecting = false;
      }
    }
    unawaited(_maybePersistSessionProfile());
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) {
        _prefetchListingThumbs(sorted, dirForPaths);
      }
    });
  }

  /// 列表/宫格首屏缩略图：不依赖子列表 ScrollController 是否已 attach。
  void _prefetchListingThumbs(List<RbFileEntry> entries, String dirPath) {
    final loader = _thumbs;
    if (loader == null || entries.isEmpty) {
      return;
    }
    final mq = MediaQuery.of(context);
    final bodyH = max(320.0, mq.size.height - mq.padding.top - mq.padding.bottom - 160);
    final bodyW = mq.size.width;
    late final int thumbEdge;
    late final int cap;
    if (_viewMode == _BrowseViewMode.grid) {
      const spacing = 1.0;
      final cross = (bodyW / 76.0).floor().clamp(3, 10);
      final cellW = (bodyW - spacing * (cross - 1)) / cross;
      final rowH = cellW + spacing;
      thumbEdge = (cellW * 1.05).round().clamp(72, 320);
      cap = RbThumbLoader.initialPrefetchLimit(
        viewportHeight: bodyH,
        mainAxisExtent: rowH,
        crossAxisCount: cross,
        thumbMaxEdge: thumbEdge,
      );
    } else {
      thumbEdge = RbThumbLoader.sharedThumbEdge;
      cap = RbThumbLoader.initialPrefetchLimit(
        viewportHeight: bodyH,
        mainAxisExtent: 48,
        thumbMaxEdge: thumbEdge,
      );
    }
    final slice = <({String path, RbFileEntry entry})>[];
    for (final e in entries) {
      if (e.isDirectory || !rbFileSupportsAgentThumb(e.name)) {
        continue;
      }
      slice.add((path: rbJoinRemotePath(dirPath, e.name), entry: e));
      if (slice.length >= cap) {
        break;
      }
    }
    if (slice.isEmpty) {
      return;
    }
    unawaited(
      loader.prefetch(
        slice,
        maxEdge: thumbEdge,
        bypassPause: true,
        todoCap: RbThumbLoader.drainTodoCap(slice.length),
      ),
    );
  }

  /// Agent 解析后同步 Windows `D:/…` 起始路径供后续 API 使用。
  void _syncListPathAfterResolve(String resolved, {required bool firstRoot}) {
    if (!_isWindowsRemoteAgent || !rbIsWindowsRemotePath(resolved)) {
      return;
    }
    final configured = _listPath.trim();
    if (firstRoot && (configured.isEmpty || configured == '.' || !rbIsWindowsRemotePath(configured))) {
      _listPath = resolved;
    }
  }

  Future<void> _maybePersistSessionProfile() async {
    final plat = _sessionAgentPlatform.trim();
    final list = await RbConnectionStore.loadAll();
    final i = list.indexWhere((p) => p.id == widget.profile.id);
    if (i < 0) {
      return;
    }
    final old = list[i];
    var newPath = old.path;
    if (_isWindowsRemoteAgent) {
      final root = (_browseRootPath ?? '').trim();
      if (root.isNotEmpty && rbIsWindowsRemotePath(root)) {
        final configured = old.path.trim();
        if (configured.isEmpty || configured == '.' || !rbIsWindowsRemotePath(configured)) {
          newPath = root;
        }
      }
    }
    if (newPath == old.path && (plat.isEmpty || old.agentPlatform == plat)) {
      return;
    }
    list[i] = RbConnectionProfile(
      id: old.id,
      name: old.name,
      signalUrl: old.signalUrl,
      direct: old.direct,
      room: old.room,
      path: newPath,
      agentPlatform: plat.isNotEmpty ? plat : old.agentPlatform,
    );
    await RbConnectionStore.persist(list, lastId: old.id);
  }

  void _openFile(RbFileEntry e) {
    final wire = _grpc;
    if (wire == null || e.isDirectory) {
      return;
    }
    if (!rbIsPreviewableEntry(e)) {
      unawaited(
        RbFileInfoBottomPanel.showFromList(
          context,
          entry: e,
          remoteCurrentDir: _activeDirPath,
          remoteBrowseRoot: _browseRootPath ?? _activeDirPath,
          relPath: _entryRelPath(e),
          onNavigateToDir: (dir) => unawaited(_listAt(dir)),
          onDownload: () => rbRunFileInfoDownload(
            context,
            wire: wire,
            relPath: _entryRelPath(e),
            entry: e,
          ),
        ),
      );
      return;
    }
    final fileEntries = _entries
        .where((x) => !x.isDirectory && rbIsPreviewableEntry(x))
        .toList();
    if (fileEntries.isEmpty) {
      return;
    }
    var idx = fileEntries.indexWhere((x) => x.id == e.id);
    if (idx < 0) {
      idx = 0;
    }
    unawaited(() async {
      Object? pop;
      _thumbs?.setPaused(true);
      try {
        pop = await AppNavigator.push<Object?>(
          RemoteBrowseFilePreviewPage(
            nav: RbPreviewNav.initial(
              wire: wire,
              thumbs: _thumbs,
              browsePath: _activeDirPath,
              parentBrowsePath: RbPreviewNav.parentOfBrowsePath(_activeDirPath),
              fileEntries: fileEntries,
              fileIndex: idx,
              entryRelPath: _entryRelPath,
              remoteCurrentDir: _activeDirPath,
              remoteBrowseRoot: _browseRootPath ?? _activeDirPath,
            ),
            previewStore: _previewStore,
            wireResolver: () => _grpc,
            thumbsResolver: () => _thumbs,
            wireEpoch: _wireEpoch,
            onWireLost: () {
              if (mounted && !_reconnecting) {
                unawaited(_scheduleReconnect(StateError('连接已断开')));
              }
            },
            onListChanged: () {
              if (mounted) {
                unawaited(_refreshList());
              }
            },
          ),
        );
      } catch (e, st) {
        if (!mounted) {
          return;
        }
        rbLog.warning('open preview failed', e, st);
        _showErrorBanner('预览失败：${rbUserMessage(e)}');
        return;
      } finally {
        _thumbs?.setPaused(false);
      }
      if (!mounted) {
        return;
      }
      if (pop is RbPreviewNavigateToDir) {
        await _listAt(pop.absoluteDirPath);
      } else if (pop == true) {
        await _refreshList();
      } else if (pop == rbPreviewResultWireLost) {
        unawaited(_scheduleReconnect(StateError('连接已断开')));
      }
    }());
  }

  String _entryRelPath(RbFileEntry e) =>
      rbJoinRemotePath(_activeDirPath, e.name);

  Future<void> _openMediaGrid() async {
    final wire = _grpc;
    if (wire == null || _selecting) {
      return;
    }
    final pop = await AppNavigator.push<Object?>(
      RbMediaGridPanel(
        wire: wire,
        thumbs: _thumbs,
        rootPath: _activeDirPath,
        browseRootPath: _browseRootPath ?? _activeDirPath,
        previewStore: _previewStore,
        wireEpoch: _wireEpoch,
        onWireLost: () {
          if (mounted && !_reconnecting) {
            unawaited(_scheduleReconnect(StateError('连接已断开')));
          }
        },
      ),
    );
    if (!mounted) {
      return;
    }
    if (pop is RbPreviewNavigateToDir) {
      await _listAt(pop.absoluteDirPath);
    }
  }

  Future<void> _listAt(String path) async {
    final wire = _grpc;
    if (wire == null) {
      return;
    }
    final normPath = rbNormRemotePath(path.trim().isEmpty ? '.' : path.trim());
    final dirChanged = !rbRemotePathsEqual(normPath, _activeDirPath);
    if (dirChanged) {
      _browseDirEpoch++;
    }
    final target = _apiListPath(_normalizeWirePath(path));
    final seq = ++_listSeq;
    setState(() {
      _loading = true;
      _errorMsg = '';
      _status = '加载…';
    });
    try {
      final list = await wire.listFiles(target);
      if (!mounted || seq != _listSeq) {
        return;
      }
      setState(() => _applyListResult(list, dirChanged: dirChanged));
    } catch (e) {
      if (mounted && seq == _listSeq) {
        _presentError(e);
      }
    } finally {
      if (mounted && seq == _listSeq) {
        setState(() => _loading = false);
      }
    }
  }

  Future<void> _enterDirectory(RbFileEntry e) async {
    if (!e.isDirectory) {
      return;
    }
    await _listAt(_joinBrowsePath(e.name));
  }

  Future<void> _goUp() async {
    if (!_canGoUp) {
      return;
    }
    await _listAt(_parentDirPath());
  }

  /// AppBar / 系统返回：子目录先上级；已在浏览根目录则断开并回到连接列表。
  void _onNavigateBack() {
    if (_selecting) {
      _exitSelection();
      return;
    }
    if (_canGoUp) {
      unawaited(_goUp());
      return;
    }
    unawaited(_disconnect());
    AppNavigator.pop();
  }

  String _joinBrowsePath(String name) => rbJoinRemotePath(_activeDirPath, name);

  String _parentDirPath() => rbRemotePathParent(_activeDirPath) ?? _activeDirPath;

  String _normalizeWirePath(String raw) {
    if (raw.trim().isEmpty) {
      return _activeDirPath;
    }
    return rbNormRemotePath(raw);
  }

  Future<void> _probeSessionAgentPlatform() async {
    var plat = widget.profile.agentPlatform.trim();
    if (_connectTarget.isSignal) {
      try {
        final probed = await probeAgentPlatform(_connectTarget.signalUri!, _room.trim());
        if (probed != null && probed.isNotEmpty) {
          plat = probed;
        }
      } catch (_) {}
    }
    if (plat.isEmpty && rbIsWindowsRemotePath(_currentPath)) {
      plat = 'windows';
    }
    _sessionAgentPlatform = plat;
  }

  Future<void> _promptGotoPath() async {
    final ctrl = TextEditingController(text: _activeDirPath);
    final target = await showDialog<String>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('跳转到目录'),
        content: RbClearTextField(
          controller: ctrl,
          autofocus: true,
          decoration: InputDecoration(
            hintText: _isWindowsRemoteAgent ? '绝对路径，如 D:/Users/…/photos' : '绝对路径，如 /Users/…/photos',
            border: const OutlineInputBorder(),
          ),
          onSubmitted: (v) => Navigator.pop(ctx, v),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('取消'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(ctx, ctrl.text),
            child: const Text('跳转'),
          ),
        ],
      ),
    );
    ctrl.dispose();
    if (target == null || !mounted || _loading) {
      return;
    }
    final path = _normalizeWirePath(target);
    if (rbRemotePathsEqual(path, _activeDirPath)) {
      return;
    }
    await _listAt(path);
  }

  Widget _buildPathBar(BuildContext context) {
    final segments = _pathSegments();
    final cs = Theme.of(context).colorScheme;
    final titleStyle = Theme.of(context).textTheme.titleSmall;
    final last = segments.length - 1;
    final pathChildren = <Widget>[
      for (var i = 0; i < segments.length; i++) ...[
        if (i > 0)
          Text(
            ' / ',
            style: TextStyle(fontSize: 14, color: cs.onSurfaceVariant),
          ),
        InkWell(
          borderRadius: BorderRadius.circular(4),
          onTap: _loading
              ? null
              : () {
                  if (!rbRemotePathsEqual(segments[i].path, _activeDirPath)) {
                    unawaited(_listAt(segments[i].path));
                  }
                },
          child: Text(
            segments[i].label,
            style: titleStyle?.copyWith(
              fontWeight: i == last ? FontWeight.w600 : FontWeight.w500,
              color: i == last ? cs.onSurface : cs.primary,
            ),
          ),
        ),
      ],
    ];
    return Row(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Expanded(
          child: SingleChildScrollView(
            scrollDirection: Axis.horizontal,
            child: Row(children: pathChildren),
          ),
        ),
        IconButton(
          visualDensity: VisualDensity.compact,
          padding: EdgeInsets.zero,
          constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
          tooltip: '跳转到目录',
          onPressed: _loading ? null : () => unawaited(_promptGotoPath()),
          icon: Icon(Icons.edit_outlined, size: 18, color: cs.onSurfaceVariant),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return PopScope(
      canPop: !_selecting && !_canGoUp,
      onPopInvokedWithResult: (didPop, _) {
        if (didPop) {
          return;
        }
        _onNavigateBack();
      },
      child: Scaffold(
        appBar: AppBar(
          title: Text(
            _selecting ? '已选 ${_selectedIds.length}' : widget.profile.name,
          ),
          leading: IconButton(
            icon: Icon(_selecting ? Icons.close : Icons.arrow_back),
            onPressed: _onNavigateBack,
          ),
          actions: [
            if (_selecting) ...[
              if (_selectedIds.length < _entries.length)
                TextButton(
                  onPressed: _loading
                      ? null
                      : () => setState(() {
                          _selectedIds
                            ..clear()
                            ..addAll(_entries.map((e) => e.id));
                        }),
                  child: const Text('全选'),
                ),
            ] else if (_connected) ...[
              IconButton(
                icon: const Icon(Icons.perm_media_outlined),
                tooltip: '媒体墙',
                onPressed: _loading ? null : () => unawaited(_openMediaGrid()),
              ),
              IconButton(
                icon: Icon(
                  _viewMode == _BrowseViewMode.list
                      ? Icons.grid_view
                      : Icons.view_list,
                ),
                tooltip: _viewMode == _BrowseViewMode.list ? '宫格' : '列表',
                onPressed: _toggleViewMode,
              ),
              if (_canGoUp)
                IconButton(
                  icon: const Icon(Icons.arrow_upward),
                  tooltip: '上级目录',
                  onPressed: _loading ? null : () => unawaited(_goUp()),
                ),
            ],
          ],
        ),
        body: RbStatusOverlayHost(
          message: _overlayBannerMessage,
          tone: _overlayBannerTone,
          child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(12, 8, 12, 4),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  if (_connected) ...[
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              _buildPathBar(context),
                              if (_entrySummary().isNotEmpty)
                                Padding(
                                  padding: const EdgeInsets.only(top: 2),
                                  child: Text(
                                    _entrySummary(),
                                    style: Theme.of(
                                      context,
                                    ).textTheme.bodySmall,
                                  ),
                                ),
                            ],
                          ),
                        ),
                        IconButton(
                          visualDensity: VisualDensity.compact,
                          icon: Icon(
                            _showConnectionDetails
                                ? Icons.expand_less
                                : Icons.info_outline,
                            size: 20,
                          ),
                          tooltip: '连接信息',
                          onPressed: () => setState(
                            () => _showConnectionDetails =
                                !_showConnectionDetails,
                          ),
                        ),
                      ],
                    ),
                    if (_showConnectionDetails) ...[
                      const SizedBox(height: 6),
                      Text(
                        '房间码：${_room.isEmpty ? "—" : _room}',
                        style: Theme.of(context).textTheme.labelSmall,
                      ),
                      Text(
                        _connectTarget.isDirect
                            ? '直连：${_connectTarget.directStorage}'
                            : '信令：${RbConnectTarget.endpointStorage(_connectTarget)}',
                        maxLines: 2,
                        overflow: TextOverflow.ellipsis,
                        style: Theme.of(context).textTheme.labelSmall,
                      ),
                    ],
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        if (_linkKind != null) _LinkChip(kind: _linkKind!),
                        const Spacer(),
                        TextButton(
                          style: TextButton.styleFrom(
                            visualDensity: VisualDensity.compact,
                            padding: const EdgeInsets.symmetric(horizontal: 8),
                          ),
                          onPressed: _loading
                              ? null
                              : () => unawaited(_disconnect()),
                          child: const Text('断开'),
                        ),
                      ],
                    ),
                  ] else ...[
                    Text(
                      '房间码 ${_room.isEmpty ? "—" : _room}',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Theme.of(context).colorScheme.onSurfaceVariant,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Expanded(
                          child: FilledButton(
                            onPressed: _loading ? null : _load,
                            child: Text(_loading ? '连接中…' : '连接'),
                          ),
                        ),
                        if (_loading) ...[
                          const SizedBox(width: 8),
                          OutlinedButton(
                            onPressed: _cancelConnect,
                            child: const Text('取消'),
                          ),
                        ],
                      ],
                    ),
                  ],
                ],
              ),
            ),
            Expanded(
              child: _grpc == null && !_reconnecting
                  ? Center(
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: _errorMsg.isNotEmpty
                            ? _ConnectErrorPanel(onRetry: _load)
                            : const RbEmptyState(
                                icon: Icons.link_off_outlined,
                                title: '尚未连接',
                                subtitle: '点「连接」开始浏览远程相册',
                              ),
                      ),
                    )
                  : _entries.isEmpty && !_loading && !_canGoUp
                  ? Center(
                      child: RbEmptyState(
                        icon: Icons.folder_open_outlined,
                        title: '空文件夹',
                        subtitle: _displayLocation(),
                      ),
                    )
                  : RefreshIndicator(
                      onRefresh: _grpc != null ? _refreshList : () async {},
                      child: _BrowseEntries(
                        key: ValueKey('${_viewMode.name}-$_browseDirEpoch'),
                        mode: _viewMode,
                        thumbs: _thumbs,
                        listThumbEpoch: _listThumbEpoch,
                        entries: _entries,
                        canGoUp: _canGoUp,
                        entryRelPath: _entryRelPath,
                        onGoUp: _goUp,
                        onEnterDir: _enterDirectory,
                        onTapFile: _onFileTap,
                        onLongPressFile: _enterSelection,
                        onToggleFile: _toggleSelection,
                        onSlideSelectApply: _applySlideSelection,
                        selecting: _selecting,
                        selectedIds: _selectedIds,
                        reconnecting: _reconnecting && _thumbs == null,
                      ),
                    ),
            ),
          ],
        ),
        ),
        bottomNavigationBar: _selecting && _selectedIds.isNotEmpty
            ? Material(
                elevation: 8,
                color: Theme.of(context).colorScheme.surfaceContainerHigh,
                child: SafeArea(
                  top: false,
                  child: SizedBox(
                    height: 56,
                    child: Row(
                      children: [
                        const SizedBox(width: 16),
                        Text(
                          '${_selectedIds.length} 项',
                          style: Theme.of(context).textTheme.titleSmall,
                        ),
                        const Spacer(),
                        FilledButton.tonalIcon(
                          onPressed: _loading || _deleting
                              ? null
                              : () => unawaited(_deleteSelected()),
                          icon: _deleting
                              ? const SizedBox(
                                  width: 18,
                                  height: 18,
                                  child: CircularProgressIndicator(
                                    strokeWidth: 2,
                                  ),
                                )
                              : Icon(
                                  Icons.delete_outline,
                                  color: Theme.of(context).colorScheme.error,
                                ),
                          label: Text(
                            _deleting ? '删除中…' : '删除',
                            style: TextStyle(
                              color: Theme.of(context).colorScheme.error,
                            ),
                          ),
                        ),
                        const SizedBox(width: 12),
                      ],
                    ),
                  ),
                ),
              )
            : null,
      ),
    );
  }
}

class _BrowseEntries extends StatelessWidget {
  const _BrowseEntries({
    super.key,
    required this.mode,
    required this.thumbs,
    required this.entries,
    required this.listThumbEpoch,
    required this.canGoUp,
    required this.entryRelPath,
    required this.onGoUp,
    required this.onEnterDir,
    required this.onTapFile,
    required this.onLongPressFile,
    required this.onToggleFile,
    required this.onSlideSelectApply,
    required this.selecting,
    required this.selectedIds,
    this.reconnecting = false,
  });

  final _BrowseViewMode mode;
  final RbThumbLoader? thumbs;
  final int listThumbEpoch;
  final List<RbFileEntry> entries;
  final bool canGoUp;
  final String Function(RbFileEntry) entryRelPath;
  final Future<void> Function() onGoUp;
  final Future<void> Function(RbFileEntry) onEnterDir;
  final void Function(RbFileEntry) onTapFile;
  final void Function(RbFileEntry) onLongPressFile;
  final void Function(RbFileEntry) onToggleFile;
  final void Function(Set<String> ids) onSlideSelectApply;
  final bool selecting;
  final Set<String> selectedIds;
  final bool reconnecting;

  Widget _buildBrowseBody() {
    if (mode == _BrowseViewMode.list) {
      return _BrowseEntryList(
        thumbs: thumbs,
        listThumbEpoch: listThumbEpoch,
        entries: entries,
        canGoUp: canGoUp,
        entryRelPath: entryRelPath,
        onGoUp: onGoUp,
        onEnterDir: onEnterDir,
        onTapFile: onTapFile,
        onLongPressFile: onLongPressFile,
        onToggleFile: onToggleFile,
        selecting: selecting,
        selectedIds: selectedIds,
      );
    }
    return _BrowseEntryGrid(
      thumbs: thumbs,
      listThumbEpoch: listThumbEpoch,
      entries: entries,
      canGoUp: canGoUp,
      entryRelPath: entryRelPath,
      onGoUp: onGoUp,
      onEnterDir: onEnterDir,
      onTapFile: onTapFile,
      onLongPressFile: onLongPressFile,
      onToggleFile: onToggleFile,
      onSlideSelectApply: onSlideSelectApply,
      selecting: selecting,
      selectedIds: selectedIds,
    );
  }

  @override
  Widget build(BuildContext context) {
    final child = _buildBrowseBody();
    if (!reconnecting) {
      return child;
    }
    return Stack(
      children: [
        Opacity(opacity: 0.55, child: child),
        const Positioned.fill(child: IgnorePointer()),
      ],
    );
  }
}

enum _BrowseRowKind { up, dir, divider, file }

class _BrowseRow {
  const _BrowseRow(this.kind, [this.entry]);

  final _BrowseRowKind kind;
  final RbFileEntry? entry;
}

class _BrowseEntryList extends StatefulWidget {
  const _BrowseEntryList({
    required this.thumbs,
    required this.listThumbEpoch,
    required this.entries,
    required this.canGoUp,
    required this.entryRelPath,
    required this.onGoUp,
    required this.onEnterDir,
    required this.onTapFile,
    required this.onLongPressFile,
    required this.onToggleFile,
    required this.selecting,
    required this.selectedIds,
  });

  final RbThumbLoader? thumbs;
  final int listThumbEpoch;
  final List<RbFileEntry> entries;
  final bool canGoUp;
  final String Function(RbFileEntry) entryRelPath;
  final Future<void> Function() onGoUp;
  final Future<void> Function(RbFileEntry) onEnterDir;
  final void Function(RbFileEntry) onTapFile;
  final void Function(RbFileEntry) onLongPressFile;
  final void Function(RbFileEntry) onToggleFile;
  final bool selecting;
  final Set<String> selectedIds;

  @override
  State<_BrowseEntryList> createState() => _BrowseEntryListState();
}

class _BrowseEntryListState extends State<_BrowseEntryList> {
  final _scroll = ScrollController();
  late List<_BrowseRow> _rows;
  late List<RbFileEntry> _fileEntries;
  var _filesOffsetY = 0.0;
  Timer? _thumbScrollDebounce;
  static const _listThumbEdge = RbThumbLoader.sharedThumbEdge;
  static const _listRowH = 48.0;

  @override
  void initState() {
    super.initState();
    _rows = _buildRows();
    _scroll.addListener(_onScroll);
    WidgetsBinding.instance.addPostFrameCallback(
      (_) => _scheduleInitialThumbPass(),
    );
  }

  void _scheduleInitialThumbPass() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) {
        return;
      }
      _prefetchListThumbs();
      _syncVisibleThumbs();
    });
  }

  void _prefetchListThumbs() {
    final loader = widget.thumbs;
    if (loader == null) {
      return;
    }
    final slice = _collectListThumbSlice(fallbackFirstScreen: true);
    if (slice.isNotEmpty) {
      unawaited(
        loader.prefetch(
          slice,
          maxEdge: _listThumbEdge,
          bypassPause: true,
          todoCap: RbThumbLoader.drainTodoCap(slice.length),
        ),
      );
    }
  }

  @override
  void didUpdateWidget(covariant _BrowseEntryList oldWidget) {
    super.didUpdateWidget(oldWidget);
    final listChanged = oldWidget.entries != widget.entries ||
        oldWidget.canGoUp != widget.canGoUp ||
        oldWidget.listThumbEpoch != widget.listThumbEpoch;
    if (listChanged) {
      _rows = _buildRows();
      WidgetsBinding.instance.addPostFrameCallback(
        (_) => _scheduleInitialThumbPass(),
      );
    } else if (oldWidget.thumbs == null && widget.thumbs != null) {
      WidgetsBinding.instance.addPostFrameCallback(
        (_) => _scheduleInitialThumbPass(),
      );
    }
  }

  @override
  void dispose() {
    _thumbScrollDebounce?.cancel();
    _scroll.removeListener(_onScroll);
    _scroll.dispose();
    super.dispose();
  }

  List<_BrowseRow> _buildRows() {
    final dirs = widget.entries.where((e) => e.isDirectory).toList();
    final files = widget.entries.where((e) => !e.isDirectory).toList();
    final rows = <_BrowseRow>[];
    var filesOffsetY = 0.0;
    if (widget.canGoUp) {
      rows.add(const _BrowseRow(_BrowseRowKind.up));
      filesOffsetY += _listRowH;
    }
    for (final e in dirs) {
      rows.add(_BrowseRow(_BrowseRowKind.dir, e));
      filesOffsetY += _listRowH;
    }
    if (dirs.isNotEmpty && files.isNotEmpty) {
      rows.add(const _BrowseRow(_BrowseRowKind.divider));
      filesOffsetY += 1;
    }
    _fileEntries = files;
    _filesOffsetY = filesOffsetY;
    for (final e in files) {
      rows.add(_BrowseRow(_BrowseRowKind.file, e));
    }
    return rows;
  }

  void _onScroll() {
    _thumbScrollDebounce?.cancel();
    _thumbScrollDebounce = Timer(const Duration(milliseconds: 150), () {
      if (mounted && _scroll.hasClients) {
        _syncVisibleThumbs();
      }
    });
  }

  double _listViewportHeight() {
    if (_scroll.hasClients) {
      final h = _scroll.position.viewportDimension;
      if (h > 1) {
        return h;
      }
    }
    return 480;
  }

  List<({String path, RbFileEntry entry})> _collectListThumbSlice({
    bool fallbackFirstScreen = false,
  }) {
    final slice = <({String path, RbFileEntry entry})>[];
    if (_fileEntries.isEmpty) {
      return slice;
    }
    final offset = _scroll.hasClients ? _scroll.offset : 0.0;
    final viewBottom = offset + _listViewportHeight();
    if (viewBottom >= _filesOffsetY) {
      final fileTop = max(0.0, offset - _filesOffsetY);
      final fileBottom = viewBottom - _filesOffsetY;
      final first = max(0, (fileTop / _listRowH).floor() - RbThumbLoader.prefetchBehind);
      final last = min(
        _fileEntries.length - 1,
        (fileBottom / _listRowH).ceil() + RbThumbLoader.prefetchAhead,
      );
      for (var i = first; i <= last; i++) {
        final e = _fileEntries[i];
        slice.add((path: widget.entryRelPath(e), entry: e));
      }
    }
    if (slice.isNotEmpty || !fallbackFirstScreen) {
      return slice;
    }
    final cap = RbThumbLoader.initialPrefetchLimit(
      viewportHeight: _listViewportHeight(),
      mainAxisExtent: _listRowH,
      thumbMaxEdge: _listThumbEdge,
    );
    for (var i = 0; i < min(cap, _fileEntries.length); i++) {
      final e = _fileEntries[i];
      slice.add((path: widget.entryRelPath(e), entry: e));
    }
    return slice;
  }

  void _syncVisibleThumbs() {
    final loader = widget.thumbs;
    if (loader == null || _rows.isEmpty) {
      return;
    }
    final slice = _collectListThumbSlice(fallbackFirstScreen: true);
    if (slice.isEmpty) {
      return;
    }
    loader.prefetchVisibleRange(slice, maxEdge: _listThumbEdge);
  }

  @override
  Widget build(BuildContext context) {
    const tilePad = EdgeInsets.symmetric(horizontal: 10, vertical: 0);
    const dense = VisualDensity.compact;
    return ListView.builder(
      scrollCacheExtent: ScrollCacheExtent.pixels(640), primary: false,
      controller: _scroll,
      physics: const AlwaysScrollableScrollPhysics(),
      padding: EdgeInsets.zero,
      itemCount: _rows.length,
      itemBuilder: (ctx, i) {
        final row = _rows[i];
        return switch (row.kind) {
          _BrowseRowKind.up => ListTile(
            dense: true,
            visualDensity: dense,
            contentPadding: tilePad,
            minVerticalPadding: 0,
            leading: const Icon(Icons.arrow_upward, size: 22),
            title: const Text('..', style: TextStyle(fontSize: 14)),
            onTap: () => unawaited(widget.onGoUp()),
          ),
          _BrowseRowKind.dir => ListTile(
            dense: true,
            visualDensity: dense,
            contentPadding: tilePad,
            minVerticalPadding: 0,
            selected:
                widget.selecting && widget.selectedIds.contains(row.entry!.id),
            selectedTileColor: Theme.of(
              context,
            ).colorScheme.primaryContainer.withValues(alpha: 0.35),
            leading: Icon(Icons.folder, size: 26, color: Colors.amber[800]),
            title: Text(
              row.entry!.name,
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontSize: 14),
            ),
            trailing: widget.selecting
                ? Icon(
                    widget.selectedIds.contains(row.entry!.id)
                        ? Icons.check_circle
                        : Icons.circle_outlined,
                    size: 20,
                    color: widget.selectedIds.contains(row.entry!.id)
                        ? Theme.of(context).colorScheme.primary
                        : null,
                  )
                : const Icon(Icons.chevron_right, size: 20),
            onTap: () {
              if (widget.selecting) {
                widget.onToggleFile(row.entry!);
              } else {
                unawaited(widget.onEnterDir(row.entry!));
              }
            },
            onLongPress: () => widget.selecting
                ? widget.onToggleFile(row.entry!)
                : widget.onLongPressFile(row.entry!),
          ),
          _BrowseRowKind.divider => const Divider(height: 1),
          _BrowseRowKind.file => _SelectableFileListTile(
            entry: row.entry!,
            filePath: widget.entryRelPath(row.entry!),
            thumbs: widget.thumbs,
            selecting: widget.selecting,
            selected: widget.selectedIds.contains(row.entry!.id),
            onTap: () => widget.onTapFile(row.entry!),
            onLongPress: () => widget.selecting
                ? widget.onToggleFile(row.entry!)
                : widget.onLongPressFile(row.entry!),
          ),
        };
      },
    );
  }
}

/// 目录宫格：与 [RbSelectableGridView] 共用框选、双指缩放、长按多选。
class _BrowseEntryGrid extends StatefulWidget {
  const _BrowseEntryGrid({
    required this.thumbs,
    required this.listThumbEpoch,
    required this.entries,
    required this.canGoUp,
    required this.entryRelPath,
    required this.onGoUp,
    required this.onEnterDir,
    required this.onTapFile,
    required this.onLongPressFile,
    required this.onToggleFile,
    required this.onSlideSelectApply,
    required this.selecting,
    required this.selectedIds,
  });

  final RbThumbLoader? thumbs;
  final int listThumbEpoch;
  final List<RbFileEntry> entries;
  final bool canGoUp;
  final String Function(RbFileEntry) entryRelPath;
  final Future<void> Function() onGoUp;
  final Future<void> Function(RbFileEntry) onEnterDir;
  final void Function(RbFileEntry) onTapFile;
  final void Function(RbFileEntry) onLongPressFile;
  final void Function(RbFileEntry) onToggleFile;
  final void Function(Set<String> ids) onSlideSelectApply;
  final bool selecting;
  final Set<String> selectedIds;

  @override
  State<_BrowseEntryGrid> createState() => _BrowseEntryGridState();
}

class _BrowseEntryGridState extends State<_BrowseEntryGrid> {
  late List<RbGridCell> _cells;

  @override
  void initState() {
    super.initState();
    _cells = _buildCells();
  }

  @override
  void didUpdateWidget(covariant _BrowseEntryGrid oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.entries != widget.entries ||
        oldWidget.canGoUp != widget.canGoUp ||
        oldWidget.listThumbEpoch != widget.listThumbEpoch) {
      _cells = _buildCells();
    }
  }

  List<RbGridCell> _buildCells() {
    final dirs = widget.entries.where((e) => e.isDirectory).toList();
    final files = widget.entries.where((e) => !e.isDirectory).toList();
    final tiles = <RbGridCell>[];
    if (widget.canGoUp) {
      tiles.add(const RbGridCell.up());
    }
    for (final e in dirs) {
      tiles.add(RbGridCell.dir(e, widget.entryRelPath(e)));
    }
    for (final e in files) {
      tiles.add(RbGridCell.file(e, widget.entryRelPath(e)));
    }
    return tiles;
  }

  @override
  Widget build(BuildContext context) {
    return RbSelectableGridView(
      key: ValueKey('grid-${widget.listThumbEpoch}-${widget.entries.length}-${widget.canGoUp}'),
      cells: _cells,
      pathForEntry: widget.entryRelPath,
      thumbs: widget.thumbs,
      selecting: widget.selecting,
      selectedIds: widget.selectedIds,
      onTapFile: widget.onTapFile,
      onLongPressFile: widget.onLongPressFile,
      onToggleFile: widget.onToggleFile,
      onSlideSelectApply: widget.onSlideSelectApply,
      onGoUp: widget.onGoUp,
      onEnterDir: widget.onEnterDir,
    );
  }
}

class _SelectableFileListTile extends StatelessWidget {
  const _SelectableFileListTile({
    required this.entry,
    required this.filePath,
    required this.thumbs,
    required this.selecting,
    required this.selected,
    required this.onTap,
    required this.onLongPress,
  });

  final RbFileEntry entry;
  final String filePath;
  final RbThumbLoader? thumbs;
  final bool selecting;
  final bool selected;
  final VoidCallback onTap;
  final VoidCallback onLongPress;

  @override
  Widget build(BuildContext context) {
    const tilePad = EdgeInsets.symmetric(horizontal: 10, vertical: 0);
    const dense = VisualDensity.compact;
    final cs = Theme.of(context).colorScheme;
    return ListTile(
      dense: true,
      visualDensity: dense,
      contentPadding: tilePad,
      minVerticalPadding: 0,
      selected: selected,
      selectedTileColor: cs.primaryContainer.withValues(alpha: 0.35),
      leading: SizedBox(
        width: 48,
        height: 48,
        child: Stack(
          fit: StackFit.expand,
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: RbEntryThumb(
                loader: thumbs,
                entry: entry,
                filePath: filePath,
                maxEdge: RbThumbLoader.sharedThumbEdge,
              ),
            ),
            if (entry.isMotionPhoto)
              const Positioned(right: 2, top: 2, child: RbGridLiveBadge()),
            if ((rbFileIsVideo(entry.name) || entry.isMotionPhoto) &&
                entry.durationMs > 0)
              Positioned(
                left: 2,
                bottom: 2,
                child: RbGridVideoDurationLabel(ms: entry.durationMs),
              ),
            if (selecting)
              Positioned(
                right: 2,
                bottom: 2,
                child: Icon(
                  selected ? Icons.check_circle : Icons.circle_outlined,
                  size: 16,
                  color: selected ? cs.primary : Colors.white70,
                ),
              ),
          ],
        ),
      ),
      title: Text(
        entry.name,
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
        style: const TextStyle(fontSize: 14),
      ),
      subtitle: Text(
        entry.durationMs > 0 && rbFileIsVideo(entry.name)
            ? '${rbFormatFileSize(entry.size)} · ${rbFormatDurationMs(entry.durationMs)}'
            : rbFormatFileSize(entry.size),
        style: const TextStyle(fontSize: 11),
      ),
      onTap: onTap,
      onLongPress: onLongPress,
    );
  }
}

class _ConnectErrorPanel extends StatelessWidget {
  const _ConnectErrorPanel({required this.onRetry});

  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return RbEmptyState(
      icon: Icons.error_outline,
      title: '连接失败',
      subtitle: '详情见顶部提示',
      action: FilledButton(onPressed: onRetry, child: const Text('重试')),
    );
  }
}

class _LinkChip extends StatelessWidget {
  const _LinkChip({required this.kind});

  final RbLinkKind kind;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final (bg, fg) = switch (kind) {
      RbLinkKind.directTcp => (cs.tertiaryContainer, cs.onTertiaryContainer),
      RbLinkKind.ice => (cs.primaryContainer, cs.onPrimaryContainer),
      RbLinkKind.relayTcp => (cs.secondaryContainer, cs.onSecondaryContainer),
    };
    return Tooltip(
      message: kind.detail,
      child: Chip(
        label: Text(
          kind.label,
          style: TextStyle(
            color: fg,
            fontWeight: FontWeight.w600,
            fontSize: 12,
          ),
        ),
        backgroundColor: bg,
        visualDensity: VisualDensity.compact,
        padding: EdgeInsets.zero,
        materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
      ),
    );
  }
}
