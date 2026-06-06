import 'dart:async';
import 'dart:math';

import 'package:app/pages/remotebrowse/rb_grid_tile.dart';
import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:app/remotebrowse/rb_thumb_loader.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart';

enum RbGridCellKind { up, dir, file }

class RbGridCell {
  const RbGridCell.up() : kind = RbGridCellKind.up, entry = null, relPath = null;
  const RbGridCell.dir(this.entry, this.relPath) : kind = RbGridCellKind.dir;
  const RbGridCell.file(this.entry, this.relPath) : kind = RbGridCellKind.file;

  final RbGridCellKind kind;
  final RbFileEntry? entry;
  final String? relPath;
}

/// 框选时禁止手指直接拖动列表，仅边缘自动滚动。
class RbSlideSelectScrollPhysics extends ScrollPhysics {
  const RbSlideSelectScrollPhysics({super.parent});

  @override
  RbSlideSelectScrollPhysics applyTo(ScrollPhysics? ancestor) =>
      RbSlideSelectScrollPhysics(parent: buildParent(ancestor));

  @override
  bool shouldAcceptUserOffset(ScrollMetrics position) => false;
}

enum RbGridPointerAxis { undecided, scroll, select }

class RbSelectableGridView extends StatefulWidget {
  const RbSelectableGridView({
    super.key,
    required this.cells,
    required this.pathForEntry,
    this.scrollController,
    this.thumbs,
    this.selecting = false,
    this.selectedIds = const {},
    required this.onTapFile,
    required this.onLongPressFile,
    required this.onToggleFile,
    required this.onSlideSelectApply,
    this.onGoUp,
    this.onEnterDir,
    this.pageStorageKey,
    this.findChildIndexCallback,
    this.padding = EdgeInsets.zero,
    this.onScroll,
    this.syncThumbs = true,
  });

  final List<RbGridCell> cells;
  final String Function(RbFileEntry) pathForEntry;
  final ScrollController? scrollController;
  final RbThumbLoader? thumbs;
  final bool selecting;
  final Set<String> selectedIds;
  final void Function(RbFileEntry) onTapFile;
  final void Function(RbFileEntry) onLongPressFile;
  final void Function(RbFileEntry) onToggleFile;
  final void Function(Set<String> ids) onSlideSelectApply;
  final Future<void> Function()? onGoUp;
  final Future<void> Function(RbFileEntry)? onEnterDir;
  final PageStorageKey<String>? pageStorageKey;
  final int? Function(Key)? findChildIndexCallback;
  final EdgeInsetsGeometry padding;
  final VoidCallback? onScroll;
  final bool syncThumbs;

  @override
  State<RbSelectableGridView> createState() => RbSelectableGridViewState();
}

class RbSelectableGridViewState extends State<RbSelectableGridView> {
  late final ScrollController _scroll;
  var _ownsScroll = false;
  final _gridKey = GlobalKey();

  int get crossAxisCount => _crossCount;
  double get layoutWidth => _gridWidth;

  int thumbMaxEdgeForLayout() => _thumbMaxEdgeForCell(_gridWidth > 1 ? _gridWidth : 320, _crossCount);

  var _crossCount = 4;
  var _gridWidth = 320.0;
  var _slideSelecting = false;
  var _longPressSlideActive = false;
  var _longPressEnteringSelection = false;
  bool? _longPressSlideAddOverride;
  var _ignoreFileTapAfterLongPress = false;
  var _gridPointerDragged = false;
  Timer? _gridLongPressTimer;
  RbFileEntry? _gridLongPressEntry;
  int? _gridLongPressPointer;
  Offset? _gridLongPressDownGlobal;
  RbGridPointerAxis _dragAxis = RbGridPointerAxis.undecided;
  int? _slideAnchorRow;
  int? _slideAnchorCol;
  Set<String> _slideBaseIds = {};
  var _slideGestureAdd = true;
  Offset? _lastPointerDownGlobal;
  Offset? _pointerDownGlobal;
  Offset? _dragGlobal;
  Offset? _lastDragGlobal;
  Timer? _autoScrollTimer;
  var _activePointers = 0;
  final Map<int, Offset> _pinchPointers = {};
  double? _pinchStartSpan;
  double _pinchZoomAtStart = 1.0;
  double _gridZoom = 1.0;
  static const _gridBaseCellWidth = 76.0;
  static const _gridZoomMin = 0.55;
  static const _gridZoomMax = 2.4;
  static const _gridSpacing = 1.0;
  static const _gridAspect = 1.0;
  static const _axisSlop = 12.0;
  static const _gridLongPressDelay = Duration(milliseconds: 480);
  static const _gridLongPressMoveSlop = 18.0;
  static const _autoScrollEdge = 48.0;
  static const _autoScrollMaxPxPerFrame = 14.0;

  /// [_gridZoom] 越大格越大：双指张开放大，捏合缩小。列数带滞回，避免宽度微变时 Grid 重建把滚动打回顶部。
  int _crossForWidth(double width) => rbGridCrossForWidth(
        width: width,
        cellBaseWidth: _gridBaseCellWidth * _gridZoom,
        currentCross: _crossCount,
        currentWidth: _gridWidth,
      );

  int _thumbMaxEdgeForCell(double width, int cross) =>
      (_cellWidth(width, cross) * 1.05).round().clamp(72, 320);

  double _pinchSpan() {
    final pts = _pinchPointers.values.toList();
    if (pts.length < 2) {
      return 0;
    }
    return (pts[0] - pts[1]).distance;
  }

  void _applyPinchZoom() {
    final start = _pinchStartSpan;
    if (start == null || start < 8) {
      return;
    }
    final span = _pinchSpan();
    if (span < 8) {
      return;
    }
    final next = (_pinchZoomAtStart * (span / start)).clamp(
      _gridZoomMin,
      _gridZoomMax,
    );
    if ((next - _gridZoom).abs() < 0.008) {
      return;
    }
    final anchor = _scroll.hasClients ? _scroll.offset : 0.0;
    final width = _gridWidth > 1 ? _gridWidth : MediaQuery.sizeOf(context).width;
    setState(() {
      _gridZoom = next;
      _crossCount = _crossForWidth(width);
    });
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted || !_scroll.hasClients) {
        return;
      }
      final max = _scroll.position.maxScrollExtent;
      final target = anchor.clamp(0.0, max);
      if ((_scroll.offset - target).abs() > 0.5) {
        _scroll.jumpTo(target);
      }
      _syncVisibleThumbs();
    });
  }

  void _onGridPointerDown(PointerDownEvent e) {
    _lastPointerDownGlobal = e.position;
    _pinchPointers[e.pointer] = e.position;
    if (_pinchPointers.length >= 2) {
      _cancelGridLongPressTimer();
      _pinchStartSpan = _pinchSpan();
      _pinchZoomAtStart = _gridZoom;
      _resetPointerGesture();
    } else {
      _scheduleGridLongPress(e);
    }
    if (!widget.selecting || _pinchPointers.length >= 2) {
      return;
    }
    _activePointers++;
    if (_activePointers != 1) {
      _resetPointerGesture();
      return;
    }
    _pointerDownGlobal = e.position;
    _dragGlobal = e.position;
    _lastDragGlobal = e.position;
    _dragAxis = RbGridPointerAxis.undecided;
    _setSlideAnchorAt(e.position);
  }

  void _onGridPointerMove(PointerMoveEvent e) {
    _pinchPointers[e.pointer] = e.position;
    if (_pinchPointers.length >= 2) {
      _applyPinchZoom();
      return;
    }
    if (_gridLongPressPointer == e.pointer) {
      final down = _gridLongPressDownGlobal;
      if (down != null &&
          (e.position - down).distance > _gridLongPressMoveSlop) {
        _cancelGridLongPressTimer();
      }
    }
    if (_longPressSlideActive) {
      _dragGlobal = e.position;
      final last = _lastDragGlobal;
      if (last != null) {
        _applySlideAlong(last, e.position);
      } else {
        _applySlideAt(e.position, clamp: true);
      }
      _lastDragGlobal = e.position;
      _tickAutoScroll(e.position);
      _ensureAutoScrollTimer();
      return;
    }
    if (!widget.selecting ||
        _activePointers != 1 ||
        _pointerDownGlobal == null) {
      return;
    }
    if (_dragAxis == RbGridPointerAxis.undecided) {
      _resolveDragAxis(e.position);
    }
    if (_dragAxis != RbGridPointerAxis.select || !_slideSelecting) {
      return;
    }
    _dragGlobal = e.position;
    final last = _lastDragGlobal;
    if (last != null) {
      _applySlideAlong(last, e.position);
    } else {
      _applySlideAt(e.position, clamp: true);
    }
    _lastDragGlobal = e.position;
    _tickAutoScroll(e.position);
    _ensureAutoScrollTimer();
  }

  void _onGridPointerUp(PointerUpEvent e) {
    _finishGridPointer(e, cancelled: false);
  }

  void _onGridPointerCancel(PointerCancelEvent e) {
    _finishGridPointer(e, cancelled: true);
  }

  void _finishGridPointer(PointerEvent e, {required bool cancelled}) {
    if (_gridLongPressPointer == e.pointer) {
      _cancelGridLongPressTimer();
    }
    _pinchPointers.remove(e.pointer);
    if (_pinchPointers.length < 2) {
      _pinchStartSpan = null;
    }
    if (_longPressSlideActive) {
      _resetPointerGesture();
      _releaseIgnoreFileTapAfterLongPress();
    } else if (_ignoreFileTapAfterLongPress) {
      _releaseIgnoreFileTapAfterLongPress();
    } else if (widget.selecting) {
      _resetPointerGesture();
    }
    if (widget.selecting) {
      _activePointers = max(0, _activePointers - 1);
    }
    if (_gridPointerDragged) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) {
          _gridPointerDragged = false;
        }
      });
    } else {
      _gridPointerDragged = false;
    }
  }

  void _releaseIgnoreFileTapAfterLongPress() {
    if (!_ignoreFileTapAfterLongPress) {
      return;
    }
    // 抬手后 InkWell 仍可能补发 onTap，需略晚于手势识别再恢复。
    Future<void>.delayed(const Duration(milliseconds: 320), () {
      if (mounted) {
        _ignoreFileTapAfterLongPress = false;
      }
    });
  }

  @override
  void initState() {
    super.initState();
    _ownsScroll = widget.scrollController == null;
    _scroll = widget.scrollController ?? ScrollController();
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
      _prefetchGridThumbs();
      _syncVisibleThumbs();
    });
  }


  void _prefetchGridThumbs() {
    if (!widget.syncThumbs) {
      return;
    }
    final loader = widget.thumbs;
    if (loader == null) {
      return;
    }
    final slice = _collectGridThumbSlice(fallbackFirstScreen: true);
    if (slice.isNotEmpty) {
      final edge = _thumbMaxEdgeForCell(_gridWidth, _crossCount);
      unawaited(
        loader.prefetch(
          slice,
          maxEdge: edge,
          bypassPause: true,
          todoCap: RbThumbLoader.drainTodoCap(slice.length),
        ),
      );
    }
  }

  @override
  void dispose() {
    _cancelGridLongPressTimer();
    _autoScrollTimer?.cancel();
    _thumbScrollDebounce?.cancel();
    _scroll.removeListener(_onScroll);
    if (_ownsScroll) {
      _scroll.dispose();
    }
    super.dispose();
  }

  void _cancelGridLongPressTimer() {
    _gridLongPressTimer?.cancel();
    _gridLongPressTimer = null;
    _gridLongPressEntry = null;
    _gridLongPressPointer = null;
    _gridLongPressDownGlobal = null;
  }

  void _scheduleGridLongPress(PointerDownEvent e) {
    _cancelGridLongPressTimer();
    final entry = _entryAtGlobal(e.position);
    if (entry == null) {
      return;
    }
    _gridLongPressDownGlobal = e.position;
    _gridLongPressEntry = entry;
    _gridLongPressPointer = e.pointer;
    _gridLongPressTimer = Timer(_gridLongPressDelay, () {
      if (!mounted || _gridLongPressEntry == null) {
        return;
      }
      final target = _gridLongPressEntry!;
      final inSelection = widget.selecting;
      _cancelGridLongPressTimer();
      if (inSelection) {
        _onFileLongPressInSelection(target);
      } else {
        _onFileLongPressEntry(target);
      }
    });
  }

  RbFileEntry? _entryAtGlobal(Offset global, {bool clamp = false}) {
    final rc = _rowColAtGlobal(global, clamp: clamp);
    if (rc == null) {
      return null;
    }
    final idx = rc.$1 * _crossCount + rc.$2;
    if (idx < 0 || idx >= widget.cells.length) {
      return null;
    }
    final t = widget.cells[idx];
    if (t.entry == null) {
      return null;
    }
    return switch (t.kind) {
      RbGridCellKind.file || RbGridCellKind.dir => t.entry,
      _ => null,
    };
  }

  double _cellWidth(double width, int cross) =>
      (width - _gridSpacing * (cross - 1)) / cross;

  double _cellHeight(double width, int cross) =>
      _cellWidth(width, cross) / _gridAspect;

  double _rowStride(double width, int cross) =>
      _cellHeight(width, cross) + _gridSpacing;

  Timer? _thumbScrollDebounce;

  void _onScroll() {
    widget.onScroll?.call();
    _thumbScrollDebounce?.cancel();
    _thumbScrollDebounce = Timer(const Duration(milliseconds: 150), () {
      if (mounted && _scroll.hasClients) {
        _syncVisibleThumbs();
      }
    });
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    final w = MediaQuery.sizeOf(context).width;
    final cross = _crossForWidth(w);
    if (_gridWidth <= 0) {
      _gridWidth = w;
      _crossCount = cross;
      return;
    }
    if (cross != _crossCount) {
      final anchor = _scroll.hasClients ? _scroll.offset : 0.0;
      setState(() {
        _gridWidth = w;
        _crossCount = cross;
      });
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (!mounted || !_scroll.hasClients) {
          return;
        }
        final max = _scroll.position.maxScrollExtent;
        final target = anchor.clamp(0.0, max);
        if ((_scroll.offset - target).abs() > 0.5) {
          _scroll.jumpTo(target);
        }
      });
    } else if ((w - _gridWidth).abs() >= 8) {
      _gridWidth = w;
    }
  }

  double _gridViewportHeight() {
    if (!_scroll.hasClients) {
      return 480;
    }
    final h = _scroll.position.viewportDimension;
    return h > 1 ? h : 480;
  }

  List<({String path, RbFileEntry entry})> _collectGridThumbSlice({
    bool fallbackFirstScreen = false,
  }) {
    final slice = <({String path, RbFileEntry entry})>[];
    if (widget.cells.isEmpty) {
      return slice;
    }
    final rowH = _rowStride(_gridWidth, _crossCount);
    if (rowH <= 0) {
      return slice;
    }
    final offset = _scroll.hasClients ? _scroll.offset : 0.0;
    final viewBottom = offset + _gridViewportHeight();
    final firstRow = max(
      0,
      (offset / rowH).floor() - RbThumbLoader.prefetchBehind,
    );
    final lastRow = (viewBottom / rowH).ceil() + RbThumbLoader.prefetchAhead;
    final first = firstRow * _crossCount;
    final last = min(widget.cells.length - 1, (lastRow + 1) * _crossCount - 1);
    for (var i = first; i <= last; i++) {
      final t = widget.cells[i];
      if (t.kind == RbGridCellKind.file && t.entry != null) {
        slice.add((path: _cellPath(t), entry: t.entry!));
      }
    }
    if (slice.isNotEmpty || !fallbackFirstScreen) {
      return slice;
    }
    final maxCells = RbThumbLoader.initialPrefetchLimit(
      viewportHeight: _gridViewportHeight(),
      mainAxisExtent: rowH,
      crossAxisCount: _crossCount,
      thumbMaxEdge: _thumbMaxEdgeForCell(_gridWidth, _crossCount),
    );
    final seen = slice.map((e) => rbNormRemotePath(e.path)).toSet();
    for (var i = 0; i < widget.cells.length && slice.length < maxCells; i++) {
      final t = widget.cells[i];
      if (t.kind == RbGridCellKind.file && t.entry != null) {
        final p = rbNormRemotePath(_cellPath(t));
        if (seen.add(p)) {
          slice.add((path: _cellPath(t), entry: t.entry!));
        }
      }
    }
    return slice;
  }

  void _syncVisibleThumbs() {
    if (!widget.syncThumbs) {
      return;
    }
    final loader = widget.thumbs;
    if (loader == null || widget.cells.isEmpty) {
      return;
    }
    final slice = _collectGridThumbSlice(fallbackFirstScreen: true);
    if (slice.isEmpty) {
      return;
    }
    loader.prefetchVisibleRange(
      slice,
      maxEdge: _thumbMaxEdgeForCell(_gridWidth, _crossCount),
    );
  }

  void _setSlideAnchorAt(Offset global) {
    final rc = _rowColAtGlobal(global);
    if (rc == null) {
      _slideAnchorRow = null;
      _slideAnchorCol = null;
      return;
    }
    _slideAnchorRow = rc.$1;
    _slideAnchorCol = rc.$2;
  }

  int _lastOccupiedColInRow(int row) {
    final cross = _crossCount;
    final endIdx = min(widget.cells.length - 1, (row + 1) * cross - 1);
    return endIdx % cross;
  }

  String? _entryIdAtRowCol(int row, int col) {
    final idx = row * _crossCount + col;
    if (idx < 0 || idx >= widget.cells.length) {
      return null;
    }
    final t = widget.cells[idx];
    if ((t.kind == RbGridCellKind.file || t.kind == RbGridCellKind.dir) &&
        t.entry != null) {
      return t.entry!.id;
    }
    return null;
  }

  bool _slideRegionIsAddMode() {
    if (_longPressSlideActive && _longPressEnteringSelection) {
      return true;
    }
    if (_longPressSlideActive && _longPressSlideAddOverride != null) {
      return _longPressSlideAddOverride!;
    }
    return _slideGestureAdd;
  }

  void _beginSlideGestureSession() {
    _slideBaseIds = Set<String>.from(widget.selectedIds);
    final ar = _slideAnchorRow;
    final ac = _slideAnchorCol;
    final anchorId = ar != null && ac != null ? _entryIdAtRowCol(ar, ac) : null;
    _slideGestureAdd = anchorId == null || !_slideBaseIds.contains(anchorId);
  }

  /// 多行梯形区域；单行闭区间。返回区域内文件 id。
  Set<String> _entryIdsInSlideRegion(int ar, int ac, int cr, int cc) {
    final cross = _crossCount;
    final rTop = min(ar, cr);
    final rBottom = max(ar, cr);
    final dragDown = ar < cr || (ar == cr && ac <= cc);
    final topCol = dragDown ? ac : cc;
    final bottomCol = dragDown ? cc : ac;
    final ids = <String>{};

    for (var i = 0; i < widget.cells.length; i++) {
      final t = widget.cells[i];
      if (t.entry == null ||
          (t.kind != RbGridCellKind.file && t.kind != RbGridCellKind.dir)) {
        continue;
      }
      final r = i ~/ cross;
      final c = i % cross;
      if (r < rTop || r > rBottom) {
        continue;
      }
      final int c0;
      final int c1;
      if (rTop == rBottom) {
        c0 = min(ac, cc);
        c1 = max(ac, cc);
      } else if (r == rTop) {
        c0 = topCol;
        c1 = _lastOccupiedColInRow(r);
      } else if (r == rBottom) {
        c0 = 0;
        c1 = bottomCol;
      } else {
        c0 = 0;
        c1 = _lastOccupiedColInRow(r);
      }
      if (c >= c0 && c <= c1) {
        ids.add(t.entry!.id);
      }
    }
    return ids;
  }

  /// 以手势开始时的选中为底，当前梯形区域加选或减选；手指回退时区域缩小选中同步回退。
  void _commitSlideRegion(int ar, int ac, int cr, int cc) {
    final region = _entryIdsInSlideRegion(ar, ac, cr, cc);
    final next = Set<String>.from(_slideBaseIds);
    final add = _slideRegionIsAddMode();
    if (add) {
      next.addAll(region);
    } else {
      next.removeAll(region);
    }
    widget.onSlideSelectApply(next);
  }

  void _applySlideRectByRowCol(int ar, int ac, int cr, int cc) =>
      _commitSlideRegion(ar, ac, cr, cc);

  void _applySlideAt(Offset global, {bool clamp = false}) {
    final cur = _rowColAtGlobal(global, clamp: clamp);
    final ar = _slideAnchorRow;
    final ac = _slideAnchorCol;
    if (cur == null || ar == null || ac == null) {
      return;
    }
    _applySlideRectByRowCol(ar, ac, cur.$1, cur.$2);
  }

  void _applySlideAlong(Offset from, Offset to) {
    final dist = (to - from).distance;
    final steps = max(1, (dist / 4).ceil());
    for (var i = 1; i <= steps; i++) {
      final t = i / steps;
      _applySlideAt(
        Offset(
          from.dx + (to.dx - from.dx) * t,
          from.dy + (to.dy - from.dy) * t,
        ),
        clamp: true,
      );
    }
  }

  void _tickAutoScroll(Offset global) {
    if (!_scroll.hasClients) {
      return;
    }
    final box = _gridKey.currentContext?.findRenderObject() as RenderBox?;
    if (box == null || !box.hasSize) {
      return;
    }
    final local = box.globalToLocal(global);
    final h = box.size.height;
    var delta = 0.0;
    if (local.dy < _autoScrollEdge) {
      final t = 1 - local.dy / _autoScrollEdge;
      delta = -_autoScrollMaxPxPerFrame * t;
    } else if (local.dy > h - _autoScrollEdge) {
      final t = (local.dy - (h - _autoScrollEdge)) / _autoScrollEdge;
      delta = _autoScrollMaxPxPerFrame * t;
    }
    if (delta == 0) {
      return;
    }
    final pos = _scroll.position;
    _scroll.jumpTo(
      (pos.pixels + delta).clamp(pos.minScrollExtent, pos.maxScrollExtent),
    );
  }

  void _ensureAutoScrollTimer() {
    if (_autoScrollTimer != null) {
      return;
    }
    _autoScrollTimer = Timer.periodic(const Duration(milliseconds: 16), (_) {
      final g = _dragGlobal;
      if (!_slideSelecting || g == null) {
        return;
      }
      _tickAutoScroll(g);
      _applySlideAt(g, clamp: true);
    });
  }

  void _stopAutoScrollTimer() {
    _autoScrollTimer?.cancel();
    _autoScrollTimer = null;
  }

  /// 超过 slop 后按首段位移方向锁定：竖向=滚动，横向=框选（与系统相册一致）。
  void _resolveDragAxis(Offset global) {
    final down = _pointerDownGlobal;
    if (down == null || _dragAxis != RbGridPointerAxis.undecided) {
      return;
    }
    final delta = global - down;
    if (delta.distance < _axisSlop) {
      return;
    }
    if (delta.dy.abs() > delta.dx.abs()) {
      _dragAxis = RbGridPointerAxis.scroll;
      return;
    }
    _dragAxis = RbGridPointerAxis.select;
    _gridPointerDragged = true;
    _cancelGridLongPressTimer();
    _beginSlideSelect();
  }

  void _beginSlideSelect() {
    if (_slideSelecting) {
      return;
    }
    _cancelGridLongPressTimer();
    setState(() => _slideSelecting = true);
    _beginSlideGestureSession();
    final ar = _slideAnchorRow;
    final ac = _slideAnchorCol;
    if (ar != null && ac != null) {
      _commitSlideRegion(ar, ac, ar, ac);
    }
    _ensureAutoScrollTimer();
  }

  void _resetPointerGesture() {
    _stopAutoScrollTimer();
    _longPressSlideActive = false;
    _longPressEnteringSelection = false;
    _longPressSlideAddOverride = null;
    _slideBaseIds = {};
    _pointerDownGlobal = null;
    _dragGlobal = null;
    _lastDragGlobal = null;
    _dragAxis = RbGridPointerAxis.undecided;
    if (!_slideSelecting) {
      _slideAnchorRow = null;
      _slideAnchorCol = null;
      return;
    }
    setState(() {
      _slideSelecting = false;
      _slideAnchorRow = null;
      _slideAnchorCol = null;
    });
  }

  void _onGridFileTap(RbFileEntry entry) {
    if (_ignoreFileTapAfterLongPress) {
      return;
    }
    widget.onTapFile(entry);
  }

  void _onGridFileTapInSelection(RbFileEntry entry) {
    if (_ignoreFileTapAfterLongPress ||
        _longPressSlideActive ||
        _gridPointerDragged) {
      return;
    }
    widget.onToggleFile(entry);
  }

  /// 多选模式下长按：切换该项选中，并可继续拖动框选（锚点已选=减选，未选=加选）。
  void _onFileLongPressInSelection(RbFileEntry entry) {
    _ignoreFileTapAfterLongPress = true;
    _longPressSlideActive = true;
    _longPressEnteringSelection = false;
    final wasSelected = widget.selectedIds.contains(entry.id);
    _longPressSlideAddOverride = !wasSelected;
    widget.onToggleFile(entry);
    final pos = _gridLongPressDownGlobal ?? _lastPointerDownGlobal;
    if (pos == null) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) {
        return;
      }
      _beginLongPressSlideSession(pos);
    });
  }

  /// 宫格 Listener 定时长按进多选；拖动由同层 Listener 承接（避免与 GridView 滚动手势冲突）。
  void _onFileLongPressEntry(RbFileEntry entry) {
    _ignoreFileTapAfterLongPress = true;
    _longPressSlideActive = true;
    _longPressEnteringSelection = true;
    widget.onLongPressFile(entry);
    final pos = _gridLongPressDownGlobal ?? _lastPointerDownGlobal;
    if (pos == null) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) {
        return;
      }
      _beginLongPressSlideSession(pos);
    });
  }

  void _beginLongPressSlideSession(Offset global) {
    _longPressSlideActive = true;
    _pointerDownGlobal = global;
    _dragGlobal = global;
    _lastDragGlobal = global;
    _activePointers = 1;
    _dragAxis = RbGridPointerAxis.select;
    _setSlideAnchorAt(global);
    setState(() => _slideSelecting = true);
    _beginSlideGestureSession();
    final ar = _slideAnchorRow;
    final ac = _slideAnchorCol;
    if (ar != null && ac != null) {
      _commitSlideRegion(ar, ac, ar, ac);
    }
    _ensureAutoScrollTimer();
  }

  /// 手指在宫格中的行列（与 .. / 目录 / 文件混排时的网格坐标一致）。
  (int, int)? _rowColAtGlobal(Offset global, {bool clamp = false}) {
    final box = _gridKey.currentContext?.findRenderObject() as RenderBox?;
    if (box == null || !box.hasSize) {
      return null;
    }
    final local = box.globalToLocal(global);
    final cross = _crossCount;
    final w = box.size.width;
    final cellW = _cellWidth(w, cross);
    final stride = _rowStride(w, cross);
    final contentY = local.dy + (_scroll.hasClients ? _scroll.offset : 0);
    final stepX = cellW + _gridSpacing;
    final cellH = _cellHeight(w, cross);
    final maxRow = max(0, (widget.cells.length - 1) ~/ cross);
    var col = (local.dx / stepX).floor();
    var row = (contentY / stride).floor();
    final xInStep = local.dx - col * stepX;
    if (xInStep > cellW && col < cross - 1) {
      col++;
    }
    final yInStep = contentY - row * stride;
    if (yInStep > cellH && row < maxRow) {
      row++;
    }
    if (clamp) {
      return (row.clamp(0, maxRow), col.clamp(0, cross - 1));
    }
    if (col < 0 || col >= cross || row < 0) {
      return null;
    }
    if (row * cross + col >= widget.cells.length) {
      return null;
    }
    return (row, col);
  }


  String _cellPath(RbGridCell c) {
    final p = c.relPath?.trim() ?? '';
    if (p.isNotEmpty) {
      return p;
    }
    final e = c.entry;
    return e != null ? widget.pathForEntry(e) : '';
  }

  @override
  Widget build(BuildContext context) {
    final w = _gridWidth > 1 ? _gridWidth : MediaQuery.sizeOf(context).width;
    final cross = _crossCount;
    final thumbEdge = _thumbMaxEdgeForCell(w, cross);
    final gridView = GridView.builder(
      scrollCacheExtent: ScrollCacheExtent.pixels(640), key: widget.pageStorageKey,
      primary: false,
      findChildIndexCallback: widget.findChildIndexCallback,
      controller: _scroll,
      physics: _slideSelecting
          ? const RbSlideSelectScrollPhysics(parent: AlwaysScrollableScrollPhysics())
          : const AlwaysScrollableScrollPhysics(),
      padding: widget.padding,
      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: cross,
        mainAxisSpacing: _gridSpacing,
        crossAxisSpacing: _gridSpacing,
        childAspectRatio: _gridAspect,
      ),
      itemCount: widget.cells.length,
      itemBuilder: (ctx, i) {
        final t = widget.cells[i];
        return switch (t.kind) {
          RbGridCellKind.up => RbGridUpTile(onTap: () { final fn = widget.onGoUp; if (fn != null) unawaited(fn()); }),
          RbGridCellKind.dir => RbGridDirTile(
            entry: t.entry!,
            selecting: widget.selecting,
            selected: widget.selectedIds.contains(t.entry!.id),
            onTap: widget.selecting
                ? () => widget.onToggleFile(t.entry!)
                : () { final fn = widget.onEnterDir; if (fn != null) unawaited(fn(t.entry!)); },
          ),
          RbGridCellKind.file => RbGridMediaTile(
            key: ValueKey(_cellPath(t)),
            thumbs: widget.thumbs,
            entry: t.entry!,
            filePath: _cellPath(t),
            thumbMaxEdge: thumbEdge,
            selecting: widget.selecting,
            selected: widget.selectedIds.contains(t.entry!.id),
            onTap: widget.selecting
                ? () => _onGridFileTapInSelection(t.entry!)
                : () => _onGridFileTap(t.entry!),
          ),
        };
      },
    );
    return Listener(
      key: _gridKey,
      behavior: HitTestBehavior.translucent,
      onPointerDown: _onGridPointerDown,
      onPointerMove: _onGridPointerMove,
      onPointerUp: _onGridPointerUp,
      onPointerCancel: _onGridPointerCancel,
      child: gridView,
    );
  }
}

