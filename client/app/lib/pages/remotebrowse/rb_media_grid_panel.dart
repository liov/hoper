import 'dart:async';
import 'dart:math';

import 'package:app/pages/remotebrowse/rb_file_info_panel.dart';
import 'package:app/pages/remotebrowse/rb_selectable_grid.dart';
import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/pages/remotebrowse/remote_browse_file_preview_page.dart';
import 'package:app/remotebrowse/rb_file_download.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/rb_media_scan.dart';
import 'package:app/remotebrowse/rb_preview_nav.dart';
import 'package:app/remotebrowse/rb_preview_store.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:app/remotebrowse/rb_thumb_loader.dart';
import 'package:app/remotebrowse/rb_user_message.dart';
import 'package:app/util/nav.dart';
import 'package:flutter/material.dart';

/// 当前目录及子目录媒体宫格：滚动接近底部时按 [listMedia] cursor 续拉，缩略图仅可见区加载。
class RbMediaGridPanel extends StatefulWidget {
  const RbMediaGridPanel({
    super.key,
    required this.wire,
    required this.thumbs,
    required this.rootPath,
    required this.browseRootPath,
    this.previewStore,
    this.wireEpoch,
    this.onWireLost,
  });

  final RbGrpcSession wire;
  final RbThumbLoader? thumbs;
  final String rootPath;
  final String browseRootPath;
  final RbPreviewContentStore? previewStore;
  final ValueNotifier<int>? wireEpoch;
  final VoidCallback? onWireLost;

  @override
  State<RbMediaGridPanel> createState() => _RbMediaGridPanelState();
}

class _RbMediaGridPanelState extends State<RbMediaGridPanel> with RbTopNoticeMixin {
  static const _gridSpacing = 1.0;

  final _items = <RbMediaItem>[];
  final _scroll = ScrollController();
  final _gridKey = GlobalKey<RbSelectableGridViewState>();
  final _selectedIds = <String>{};
  var _cells = <RbGridCell>[];
  RbMediaDfsScanner? _scanner;
  var _loading = false;
  var _initialLoaded = false;
  var _exhausted = false;
  var _selecting = false;
  var _deleting = false;
  Timer? _thumbScrollDebounce;
  var _loadQueued = false;
  var _stayAtScrollEnd = false;
  var _autoLoadRounds = 0;
  static const _maxAutoLoadRounds = 2;

  @override
  void initState() {
    super.initState();
    _scanner = RbMediaDfsScanner(wire: widget.wire, rootPath: widget.rootPath);
    WidgetsBinding.instance.addPostFrameCallback((_) => _requestLoadMore(force: true));
  }

  int get _crossCount => _gridKey.currentState?.crossAxisCount ?? 4;

  double get _gridWidth => _gridKey.currentState?.layoutWidth ?? 320;

  void _exitSelection() {
    setState(() {
      _selecting = false;
      _selectedIds.clear();
    });
  }

  void _enterSelection(RbFileEntry e) {
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
    setState(() {
      _selectedIds
        ..clear()
        ..addAll(ids);
      if (_selectedIds.isEmpty) {
        _selecting = false;
      }
    });
  }

  void _pinScrollOffset(double anchor) {
    void apply() {
      if (!mounted || !_scroll.hasClients) {
        return;
      }
      final max = _scroll.position.maxScrollExtent;
      final target = anchor.clamp(0.0, max);
      if ((_scroll.offset - target).abs() > 0.5) {
        _scroll.jumpTo(target);
      }
    }

    WidgetsBinding.instance.addPostFrameCallback((_) {
      WidgetsBinding.instance.addPostFrameCallback((_) => apply());
    });
  }

  @override
  void dispose() {
    _scanner?.cancel();
    _thumbScrollDebounce?.cancel();
    _scroll.dispose();
    super.dispose();
  }

  int _batchSize() => rbGridMediaFetchBatchSize(_crossCount);

  double _rowStride() {
    final cross = _crossCount;
    final w = _gridWidth;
    final rowH = (w - _gridSpacing * (cross - 1)) / cross + _gridSpacing;
    return rowH > 0 ? rowH : 80;
  }

  double _gridViewportHeight() {
    if (!_scroll.hasClients) {
      return 480;
    }
    final h = _scroll.position.viewportDimension;
    return h > 1 ? h : 480;
  }

  int _thumbMaxEdge() =>
      _gridKey.currentState?.thumbMaxEdgeForLayout() ??
      (RbThumbLoader.sharedThumbEdge * 0.75).round();

  bool _atScrollEnd() {
    if (!_scroll.hasClients) {
      return _items.isEmpty;
    }
    final pos = _scroll.position;
    if (!pos.hasContentDimensions) {
      return false;
    }
    final slack = max(1.0, _rowStride());
    return pos.pixels >= pos.maxScrollExtent - slack;
  }

  bool _needsMoreItems({bool includeScrollEnd = true}) {
    if (_items.isEmpty) {
      return true;
    }
    if (includeScrollEnd && _atScrollEnd()) {
      return true;
    }
    return _shouldPrefetchMore();
  }

  void _onScroll() {
    _autoLoadRounds = 0;
    _thumbScrollDebounce?.cancel();
    _thumbScrollDebounce = Timer(const Duration(milliseconds: 120), () {
      if (mounted) {
        _syncVisibleThumbs();
      }
    });
    if (_exhausted || _selecting) {
      return;
    }
    if (_needsMoreItems()) {
      _requestLoadMore();
    }
  }

  /// 距列表底部还有多少行时开始用 cursor 拉下一页（约 2 屏，避免滑到底才等网络）。
  int _prefetchRowBuffer() {
    final rowH = _rowStride();
    if (rowH <= 0) {
      return 16;
    }
    final viewRows = (_gridViewportHeight() / rowH).ceil().clamp(4, 24);
    return max(12, viewRows * 2);
  }

  bool _shouldPrefetchMore() {
    if (_items.isEmpty) {
      return true;
    }
    if (!_scroll.hasClients) {
      return false;
    }
    final pos = _scroll.position;
    if (!pos.hasPixels || !pos.hasContentDimensions) {
      return false;
    }
    final rowH = _rowStride();
    final lastVisible = ((pos.pixels + pos.viewportDimension) / rowH).ceil() * _crossCount;
    final buffer = _crossCount * _prefetchRowBuffer();
    return lastVisible >= _items.length - buffer;
  }

  void _requestLoadMore({bool force = false}) {
    if (_exhausted || _selecting) {
      return;
    }
    if (_loading) {
      if (force || _needsMoreItems()) {
        _loadQueued = true;
      }
      return;
    }
    unawaited(_loadMore(force: force));
  }

  void _scheduleFillViewport() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted || _loading || _exhausted || !_scroll.hasClients) {
        return;
      }
      final pos = _scroll.position;
      if (!pos.hasPixels || pos.maxScrollExtent > 80) {
        return;
      }
      _requestLoadMore(force: true);
    });
  }

  void _scheduleContinueLoading({required List<RbMediaItem> batch, required bool stayAtScrollEnd}) {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted || _loading || _exhausted) {
        return;
      }
      final wantsMore = _loadQueued || batch.isEmpty || stayAtScrollEnd || _needsMoreItems();
      if (!wantsMore) {
        _loadQueued = false;
        return;
      }
      if (!stayAtScrollEnd && _autoLoadRounds >= _maxAutoLoadRounds) {
        _loadQueued = false;
        return;
      }
      _autoLoadRounds++;
      _loadQueued = false;
      _requestLoadMore(force: batch.isEmpty);
    });
  }

  RbGridCell _cellForItem(RbMediaItem it) {
    return RbGridCell.file(
      RbFileEntry(
        id: _mediaStableId(it),
        name: it.entry.name,
        size: it.entry.size,
        mtimeUnixMs: it.entry.mtimeUnixMs,
        thumbHash: it.entry.thumbHash,
        isDirectory: it.entry.isDirectory,
        durationMs: it.entry.durationMs,
        isMotionPhoto: it.entry.isMotionPhoto,
        motionOffset: it.entry.motionOffset,
        motionLength: it.entry.motionLength,
        motionCompanion: it.entry.motionCompanion,
      ),
      it.relPath,
    );
  }

  void _rebuildCells() {
    _cells = [for (final it in _items) _cellForItem(it)];
  }

  void _appendItems(List<RbMediaItem> batch, {required bool exhausted}) {
    final anchor = _scroll.hasClients ? _scroll.offset : 0.0;
    ScrollHoldController? hold;
    if (_scroll.hasClients) {
      hold = _scroll.position.hold(() {});
    }
    setState(() {
      _items.addAll(batch);
      _exhausted = exhausted;
      _initialLoaded = true;
      _loading = false;
      _cells.addAll([for (final it in batch) _cellForItem(it)]);
    });
    rbClearLoadingNotice();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      hold?.cancel();
      _pinScrollOffset(anchor);
    });
  }

  Future<void> _loadMore({bool force = false}) async {
    final scanner = _scanner;
    if (scanner == null || _loading || _exhausted || _selecting) {
      return;
    }
    if (!force && !_needsMoreItems()) {
      return;
    }
    _loading = true;
    _loadQueued = false;
    _stayAtScrollEnd = _atScrollEnd();
    rbClearErrorNotice();
    if (_items.isNotEmpty) {
      rbSetLoadingNotice('加载更多…');
    } else {
      rbClearLoadingNotice();
    }
    if (_items.isEmpty) {
      setState(() {});
    }
    try {
      final batch = await scanner.fetchItems(count: _batchSize());
      if (!mounted) {
        return;
      }
      _appendItems(batch, exhausted: scanner.isExhausted);
      _syncVisibleThumbs();
      _scheduleFillViewport();
      _scheduleContinueLoading(batch: batch, stayAtScrollEnd: _stayAtScrollEnd);
    } catch (e, st) {
      rbLog.warning('media grid load failed', e, st);
      if (!mounted) {
        return;
      }
      final anchor = _scroll.hasClients ? _scroll.offset : 0.0;
      setState(() {
        _loading = false;
        _initialLoaded = true;
      });
      rbClearLoadingNotice();
      rbShowErrorNotice(rbUserMessage(e));
      _loadQueued = false;
      _pinScrollOffset(anchor);
    } finally {
      if (mounted && _loading && _items.isEmpty) {
        setState(() => _loading = false);
        rbClearLoadingNotice();
      }
    }
  }

  void _syncVisibleThumbs() {
    final loader = widget.thumbs;
    if (loader == null || _items.isEmpty) {
      return;
    }
    final slice = _collectThumbSlice(fallbackFirstScreen: true);
    if (slice.isEmpty) {
      return;
    }
    loader.prefetchVisibleRange(slice, maxEdge: _thumbMaxEdge(), bypassPause: true);
  }

  List<({String path, RbFileEntry entry})> _collectThumbSlice({bool fallbackFirstScreen = false}) {
    final slice = <({String path, RbFileEntry entry})>[];
    if (_items.isEmpty) {
      return slice;
    }
    final rowH = _rowStride();
    if (!_scroll.hasClients) {
      return _firstScreenThumbSlice(rowH, fallbackFirstScreen);
    }
    final pos = _scroll.position;
    if (!pos.hasPixels || !pos.hasContentDimensions) {
      return _firstScreenThumbSlice(rowH, fallbackFirstScreen);
    }
    final offset = pos.pixels;
    final viewBottom = offset + _gridViewportHeight();
    final firstRow = max(0, (offset / rowH).floor() - RbThumbLoader.prefetchBehind);
    final lastRow = (viewBottom / rowH).ceil() + RbThumbLoader.prefetchAhead;
    final first = firstRow * _crossCount;
    final last = min(_items.length - 1, (lastRow + 1) * _crossCount - 1);
    for (var i = first; i <= last; i++) {
      slice.add((path: _items[i].relPath, entry: _items[i].entry));
    }
    if (slice.isNotEmpty || !fallbackFirstScreen) {
      return slice;
    }
    return _firstScreenThumbSlice(rowH, true);
  }

  List<({String path, RbFileEntry entry})> _firstScreenThumbSlice(double rowH, bool enabled) {
    if (!enabled) {
      return const [];
    }
    final maxCells = RbThumbLoader.initialPrefetchLimit(
      viewportHeight: _gridViewportHeight(),
      mainAxisExtent: rowH,
      crossAxisCount: _crossCount,
      thumbMaxEdge: _thumbMaxEdge(),
    );
    return [
      for (var i = 0; i < min(_items.length, maxCells); i++) (path: _items[i].relPath, entry: _items[i].entry),
    ];
  }

  int? _indexForTileKey(Key key) {
    if (key is! ValueKey<String>) {
      return null;
    }
    final path = key.value;
    final i = _items.indexWhere((e) => e.relPath == path);
    return i >= 0 ? i : null;
  }

  String _mediaStableId(RbMediaItem it) => rbNormRemotePath(it.relPath);

  String _relPathForEntry(RbFileEntry e) {
    for (final it in _items) {
      if (_mediaStableId(it) == e.id) {
        return it.relPath;
      }
    }
    return '';
  }

  String _parentDirForItem(String relPath) =>
      rbRemotePathParent(rbNormRemotePath(relPath)) ?? rbNormRemotePath(widget.rootPath);

  void _onMediaFileTap(RbFileEntry entry) {
    if (_selecting) {
      _toggleSelection(entry);
      return;
    }
    final i = _items.indexWhere((x) => _mediaStableId(x) == entry.id);
    if (i < 0) {
      return;
    }
    _openEntryAt(i);
  }

  void _openEntryAt(int index) {
    final it = _items[index];
    final entry = it.entry;
    if (!rbIsPreviewableEntry(entry)) {
      unawaited(
        RbFileInfoBottomPanel.showFromList(
          context,
          entry: entry,
          remoteCurrentDir: _parentDirForItem(it.relPath),
          remoteBrowseRoot: widget.browseRootPath,
          relPath: it.relPath,
          onNavigateToDir: (dir) {
            if (mounted) {
              Navigator.pop(context, RbPreviewNavigateToDir(dir));
            }
          },
          onDownload: () => rbRunFileInfoDownload(
            context,
            wire: widget.wire,
            relPath: it.relPath,
            entry: entry,
          ),
        ),
      );
      return;
    }
    final previewables = [
      for (final x in _items)
        if (rbIsPreviewableEntry(x.entry)) RbPreviewFileSlot(entry: x.entry, relPath: x.relPath),
    ];
    if (previewables.isEmpty) {
      return;
    }
    var idx = previewables.indexWhere((s) => s.entry.id == entry.id && s.relPath == it.relPath);
    if (idx < 0) {
      idx = 0;
    }
    unawaited(_openPreview(idx, previewables));
  }

  Future<void> _deleteSelected() async {
    if (_selectedIds.isEmpty || _deleting) {
      return;
    }
    final targets = [for (final it in _items) if (_selectedIds.contains(_mediaStableId(it))) it];
    if (targets.isEmpty) {
      _exitSelection();
      return;
    }
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('删除 ${targets.length} 个文件'),
        content: Text(targets.length == 1 ? targets.first.entry.name : '确定删除所选 ${targets.length} 个文件？'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Theme.of(ctx).colorScheme.error),
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('删除'),
          ),
        ],
      ),
    );
    if (ok != true || !mounted) {
      return;
    }
    setState(() => _deleting = true);
    var failed = 0;
    final deletedIds = <String>{};
    try {
      const batchSize = 4;
      for (var i = 0; i < targets.length; i += batchSize) {
        final batch = targets.sublist(i, min(i + batchSize, targets.length));
        await Future.wait(batch.map((it) async {
          try {
            await widget.wire.deleteFile(it.relPath);
            widget.thumbs?.evict(it.relPath, entry: it.entry);
            deletedIds.add(_mediaStableId(it));
          } catch (_) {
            failed++;
          }
        }));
      }
      if (!mounted) {
        return;
      }
      setState(() {
        _items.removeWhere((e) => deletedIds.contains(_mediaStableId(e)));
        _selectedIds.removeWhere(deletedIds.contains);
        if (_selectedIds.isEmpty) {
          _selecting = false;
        }
        _deleting = false;
        _rebuildCells();
      });
      if (failed > 0) {
        rbShowErrorNotice('有 $failed 项删除失败');
      } else if (deletedIds.isNotEmpty) {
        rbShowInfoNotice('已删除 ${deletedIds.length} 项');
      }
    } catch (e) {
      if (mounted) {
        rbShowErrorNotice(rbUserMessage(e));
      }
    } finally {
      if (mounted) {
        setState(() => _deleting = false);
      }
    }
  }

  Future<void> _openPreview(int idx, List<RbPreviewFileSlot> previewables) async {
    widget.thumbs?.setPaused(true);
    Object? pop;
    try {
      pop = await AppNavigator.push<Object?>(
        RemoteBrowseFilePreviewPage(
          nav: RbPreviewNav.fromSlots(
            wire: widget.wire,
            thumbs: widget.thumbs,
            files: previewables,
            fileIndex: idx,
            remoteBrowseRoot: widget.browseRootPath,
          ),
          previewStore: widget.previewStore,
          wireResolver: () => widget.wire,
          thumbsResolver: () => widget.thumbs,
          wireEpoch: widget.wireEpoch,
          onWireLost: widget.onWireLost,
        ),
      );
    } catch (e, st) {
      rbLog.warning('media grid preview failed', e, st);
      if (mounted) {
        rbShowErrorNotice('预览失败：${rbUserMessage(e)}');
      }
      return;
    } finally {
      widget.thumbs?.setPaused(false);
    }
    if (!mounted) {
      return;
    }
    if (pop is RbPreviewNavigateToDir) {
      Navigator.pop(context, pop);
    } else if (pop == rbPreviewResultWireLost) {
      widget.onWireLost?.call();
    }
  }

  String _panelTitle() {
    final segs = rbRemotePathSegments(widget.rootPath);
    if (segs.isEmpty) {
      return '媒体';
    }
    return segs.last.label;
  }

  String? _errorBannerMessage() => rbTopNoticeMessage;

  RbBannerTone _errorBannerTone() => rbTopNoticeTone;

  bool _errorBannerShowProgress() => rbTopNoticeShowProgress;

  List<RbGridCell> _mediaCells() => _cells;

  @override
  Widget build(BuildContext context) {
    return PopScope(
      canPop: !_selecting,
      onPopInvokedWithResult: (didPop, _) {
        if (!didPop && _selecting) {
          _exitSelection();
        }
      },
      child: Scaffold(
        appBar: AppBar(
          title: Text(
            _selecting ? '已选 ${_selectedIds.length}' : _panelTitle(),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          leading: IconButton(
            icon: Icon(_selecting ? Icons.close : Icons.arrow_back),
            onPressed: () {
              if (_selecting) {
                _exitSelection();
              } else {
                Navigator.maybePop(context);
              }
            },
          ),
          actions: [
            if (_selecting && _selectedIds.length < _items.length)
              TextButton(
                onPressed: _deleting
                    ? null
                    : () => setState(() {
                          _selectedIds
                            ..clear()
                            ..addAll(_items.map(_mediaStableId));
                        }),
                child: const Text('全选'),
              ),
          ],
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
                        Text('${_selectedIds.length} 项', style: Theme.of(context).textTheme.titleSmall),
                        const Spacer(),
                        FilledButton.tonalIcon(
                          onPressed: _deleting ? null : () => unawaited(_deleteSelected()),
                          icon: _deleting
                              ? const SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2))
                              : Icon(Icons.delete_outline, color: Theme.of(context).colorScheme.error),
                          label: Text(
                            _deleting ? '删除中…' : '删除',
                            style: TextStyle(color: Theme.of(context).colorScheme.error),
                          ),
                        ),
                        const SizedBox(width: 12),
                      ],
                    ),
                  ),
                ),
              )
            : null,
        body: RbStatusOverlayHost(
          message: _errorBannerMessage(),
          tone: _errorBannerTone(),
          showProgress: _errorBannerShowProgress(),
          child: _buildBody(context),
        ),
      ),
    );
  }

  Widget _buildBody(BuildContext context) {
    if (_initialLoaded && _items.isEmpty && !_loading) {
      if (rbTopNoticeMessage != null && rbTopNoticeTone == RbBannerTone.warning) {
        return const SizedBox.shrink();
      }
      return Center(
        child: Text(
          '未找到媒体文件',
          style: Theme.of(context).textTheme.bodyLarge,
          textAlign: TextAlign.center,
        ),
      );
    }
    if (!_initialLoaded && _loading && _items.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }
    return RbSelectableGridView(
      key: _gridKey,
      scrollController: _scroll,
      pageStorageKey: const PageStorageKey<String>('rb-media-grid'),
      findChildIndexCallback: _indexForTileKey,
      padding: const EdgeInsets.all(_gridSpacing),
      syncThumbs: false,
      onScroll: _onScroll,
      cells: _mediaCells(),
      pathForEntry: _relPathForEntry,
      thumbs: widget.thumbs,
      selecting: _selecting,
      selectedIds: _selectedIds,
      onTapFile: _onMediaFileTap,
      onLongPressFile: _enterSelection,
      onToggleFile: _toggleSelection,
      onSlideSelectApply: _applySlideSelection,
    );
  }
}
