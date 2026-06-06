import 'dart:async';

import 'package:app/remotebrowse/rb_preview_edge_swipe.dart';
import 'package:app/remotebrowse/rb_video_timeline.dart';
import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

/// 视频预览手势：横向 seek；左/右条带竖滑亮度/音量；中间宽区域竖滑切文件；底部横滑切目录。
class RbVideoGestureLayer extends StatefulWidget {
  const RbVideoGestureLayer({
    super.key,
    required this.controller,
    required this.child,
    this.onEdgeVerticalSwipe,
    this.onBottomHorizontalSwipe,
    this.onHlsSeek,
    this.timelineOffsetMs = 0,
    this.timelineTotalMs = 0,
    this.playbackAnchorMs = 0,
    this.playStartedAt,
    this.pausedAbsPosMs,
    this.controlsVisible = false,
    this.holdSpeedActive = false,
    this.seekSpanPerWidth = 3.0,
    this.seekMinMsPerWidth = 150000,
    this.bottomStripHeight = 72,
    /// 左/右亮度、音量条带各占屏宽比例（中间为切文件区，宜偏小以放大中间）。
    this.sideBandFraction = 0.14,
  });

  final VideoPlayerController controller;
  final Widget child;
  final void Function(bool next)? onEdgeVerticalSwipe;
  final void Function(bool next)? onBottomHorizontalSwipe;
  final Future<void> Function(Duration target, {bool playAfter})? onHlsSeek;
  final int timelineOffsetMs;
  final int timelineTotalMs;
  final int playbackAnchorMs;
  final DateTime? playStartedAt;
  final int? pausedAbsPosMs;
  final double bottomStripHeight;
  final bool controlsVisible;
  final bool holdSpeedActive;
  /// 横向快进灵敏度：屏宽 1 倍滑动 ≈ `seekSpanPerWidth` 倍视频时长（越小越不灵敏）。
  final double seekSpanPerWidth;
  /// 长视频横向滑动兜底：整屏宽度至少跳多少毫秒。
  final int seekMinMsPerWidth;
  final double sideBandFraction;

  @override
  State<RbVideoGestureLayer> createState() => _RbVideoGestureLayerState();
}

enum _RbVideoGestureKind { none, seek, brightness, volume, file, sibling }

class _RbVideoGestureLayerState extends State<RbVideoGestureLayer> {
  static const _axisRatio = 1.15;
  static const _maxSideBandFraction = 0.22;
  static const _seekMinInterval = Duration(milliseconds: 45);
  static const _volPxStep = 28.0;
  static const _brightPxStep = 28.0;
  static const _fileSwipeMin = 56.0;

  _RbVideoGestureKind _kind = _RbVideoGestureKind.none;
  Offset? _start;
  int? _seekAnchorMs;
  int _lastSeekTargetMs = -1;
  DateTime? _lastSeekAt;
  double _gestureWidth = 0;
  double _gestureHeight = 0;
  double _volAccum = 0;
  double _brightAccum = 0;
  var _brightnessDim = 0.0;
  var _volume = 1.0;
  String? _hud;
  Timer? _hudTimer;

  @override
  void initState() {
    super.initState();
    _volume = widget.controller.value.volume.clamp(0.0, 1.0);
  }

  @override
  void didUpdateWidget(covariant RbVideoGestureLayer oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.controller != widget.controller) {
      _volume = widget.controller.value.volume.clamp(0.0, 1.0);
    }
  }

  @override
  void dispose() {
    _hudTimer?.cancel();
    super.dispose();
  }

  void _showHud(String text) {
    if (!mounted) {
      return;
    }
    _hudTimer?.cancel();
    setState(() => _hud = text);
    _hudTimer = Timer(const Duration(milliseconds: 900), () {
      if (mounted) {
        setState(() => _hud = null);
      }
    });
  }

  bool _inBottomStrip(double y, double h) => y > h - widget.bottomStripHeight;

  double get _sideFrac => widget.sideBandFraction.clamp(0.10, _maxSideBandFraction);

  double _leftBandEnd(BuildContext context, double w) {
    final inset = rbPreviewLeftEdgeInset(context);
    return (w * _sideFrac).clamp(inset + 36.0, w * _maxSideBandFraction);
  }

  double _rightBandStart(double w) => w * (1 - _sideFrac);

  bool _inLeftBand(BuildContext context, double x, double w) => x < _leftBandEnd(context, w);

  bool _inRightBand(double x, double w) => x >= _rightBandStart(w);

  bool _inCenterBand(BuildContext context, double x, double w) => !_inLeftBand(context, x, w) && !_inRightBand(x, w);

  _RbVideoGestureKind _zoneKind(BuildContext context, double x, double y, double w, double h, double dx, double dy) {
    final horiz = dx.abs() > dy.abs() * _axisRatio;
    final vert = dy.abs() > dx.abs() * _axisRatio;
    if (widget.onBottomHorizontalSwipe != null && horiz && _inBottomStrip(y, h)) {
      return _RbVideoGestureKind.sibling;
    }
    if (horiz) {
      return _RbVideoGestureKind.seek;
    }
    if (!vert) {
      return _RbVideoGestureKind.none;
    }
    if (_inCenterBand(context, x, w) && widget.onEdgeVerticalSwipe != null) {
      return _RbVideoGestureKind.file;
    }
    if (_inLeftBand(context, x, w)) {
      return _RbVideoGestureKind.brightness;
    }
    if (_inRightBand(x, w)) {
      return _RbVideoGestureKind.volume;
    }
    return _RbVideoGestureKind.none;
  }

  double _effectiveWidth(BuildContext context, BoxConstraints c) {
    if (c.maxWidth.isFinite && c.maxWidth > 0) {
      return c.maxWidth;
    }
    return MediaQuery.sizeOf(context).width;
  }

  double _effectiveHeight(BuildContext context, BoxConstraints c) {
    if (c.maxHeight.isFinite && c.maxHeight > 0) {
      return c.maxHeight;
    }
    return MediaQuery.sizeOf(context).height;
  }

  int _seekDeltaMs(double totalDx, double width, int durMs) {
    final prop = (totalDx / width * durMs * widget.seekSpanPerWidth).round();
    final floor = (totalDx / width * widget.seekMinMsPerWidth).round();
    if (prop.abs() >= floor.abs()) {
      return prop;
    }
    return floor;
  }

  bool _inChromeZone(double y, double h, Offset start) =>
      widget.controlsVisible && !_inBottomStrip(start.dy, h) && y > h - 130;

  void _onDown(PointerDownEvent e) {
    if (widget.holdSpeedActive) {
      return;
    }
    _kind = _RbVideoGestureKind.none;
    _start = e.position;
    _seekAnchorMs = null;
    _lastSeekTargetMs = -1;
    _lastSeekAt = null;
    _volAccum = 0;
    _brightAccum = 0;
  }

  void _onMove(PointerMoveEvent e, BoxConstraints c, BuildContext context) {
    if (widget.holdSpeedActive || _start == null) {
      return;
    }
    final w = _effectiveWidth(context, c);
    final h = _effectiveHeight(context, c);
    _gestureWidth = w;
    _gestureHeight = h;
    final start = _start!;
    if (_inChromeZone(e.position.dy, h, start)) {
      return;
    }
    final d = e.position - start;
    if (_kind == _RbVideoGestureKind.none) {
      if (d.distance < 8) {
        return;
      }
      _kind = _zoneKind(context, start.dx, start.dy, w, h, d.dx, d.dy);
    }
    switch (_kind) {
      case _RbVideoGestureKind.seek:
        _seekByDrag(e.position.dx - start.dx, w);
      case _RbVideoGestureKind.brightness:
        _brightAccum += e.delta.dy;
        _applyBrightnessAccum();
      case _RbVideoGestureKind.volume:
        _volAccum += e.delta.dy;
        _applyVolumeAccum();
      case _RbVideoGestureKind.file:
      case _RbVideoGestureKind.sibling:
        break;
      case _RbVideoGestureKind.none:
        break;
    }
  }

  void _finishGesture(BuildContext context, Offset end) {
    final start = _start;
    final kind = _kind;
    final seekAnchor = _seekAnchorMs;
    final w = _gestureWidth;
    _start = null;
    _kind = _RbVideoGestureKind.none;
    _seekAnchorMs = null;
    _lastSeekTargetMs = -1;
    _lastSeekAt = null;
    _volAccum = 0;
    _brightAccum = 0;
    if (start != null && !widget.holdSpeedActive) {
      final d = end - start;
      final h = _gestureHeight;
      if (kind == _RbVideoGestureKind.sibling &&
          d.distance >= _fileSwipeMin &&
          d.dx.abs() > d.dy.abs() * _axisRatio &&
          h > 0 &&
          _inBottomStrip(start.dy, h)) {
        widget.onBottomHorizontalSwipe?.call(d.dx < 0);
      } else if (kind == _RbVideoGestureKind.file &&
          d.dy.abs() >= _fileSwipeMin &&
          d.dy.abs() > d.dx.abs() * _axisRatio &&
          h > 0 &&
          _inCenterBand(context, start.dx, w)) {
        widget.onEdgeVerticalSwipe?.call(d.dy < 0);
      }
    }
    if (kind == _RbVideoGestureKind.seek && start != null && w > 0 && seekAnchor != null) {
      _applySeekTarget(seekAnchor, end.dx - start.dx, w, force: true);
    }
  }

  void _seekByDrag(double totalDx, double width) {
    final c = widget.controller;
    final v = c.value;
    if (!v.isInitialized || width <= 0) {
      return;
    }
    final totalMs = _timelineTotalMs(v);
    if (totalMs <= 0) {
      return;
    }
    _seekAnchorMs ??= _timelinePosMs(v);
    _applySeekTarget(_seekAnchorMs!, totalDx, width, force: false);
  }

  int _timelineTotalMs(VideoPlayerValue v) {
    if (widget.timelineTotalMs > 0) {
      return widget.timelineTotalMs;
    }
    return widget.timelineOffsetMs + v.duration.inMilliseconds;
  }

  int _timelinePosMs(VideoPlayerValue v) {
    return rbVideoAbsPosMs(
      playerPosMs: v.position.inMilliseconds,
      timelineOffsetMs: widget.timelineOffsetMs,
      playbackAnchorMs: widget.playbackAnchorMs,
      isPlaying: v.isPlaying,
      playStartedAt: widget.playStartedAt,
      pausedAbsPosMs: widget.pausedAbsPosMs,
      playbackSpeed: v.playbackSpeed,
    );
  }

  void _applySeekTarget(int anchorMs, double totalDx, double width, {required bool force}) {
    final c = widget.controller;
    final v = c.value;
    final totalMs = _timelineTotalMs(v);
    if (totalMs <= 0) {
      return;
    }
    const lo = 0;
    final hi = totalMs;
    final spanMs = hi - lo;
    final target = (anchorMs + _seekDeltaMs(totalDx, width, spanMs > 0 ? spanMs : totalMs)).clamp(lo, hi);
    if (!force && target == _lastSeekTargetMs) {
      return;
    }
    final now = DateTime.now();
    if (!force &&
        _lastSeekAt != null &&
        now.difference(_lastSeekAt!) < _seekMinInterval &&
        (target - _lastSeekTargetMs).abs() < 150) {
      return;
    }
    _lastSeekTargetMs = target;
    _lastSeekAt = now;
    unawaited(_seekToMs(target, hi, commit: force || widget.onHlsSeek == null));
  }

  Future<void> _seekToMs(int targetAbsMs, int absDurMs, {required bool commit}) async {
    final off = widget.timelineOffsetMs;
    final target = Duration(milliseconds: targetAbsMs);
    final dur = Duration(milliseconds: absDurMs);
    _showHud('${_fmtClock(target)} / ${_fmtClock(dur)}');
    if (!commit) {
      return;
    }
    final hlsSeek = widget.onHlsSeek;
    if (hlsSeek != null) {
      final playing = widget.controller.value.isPlaying;
      await hlsSeek(target, playAfter: playing);
      return;
    }
    final c = widget.controller;
    if (!c.value.isInitialized) {
      return;
    }
    await c.seekTo(Duration(milliseconds: targetAbsMs - off));
  }

  String _fmtClock(Duration d) {
    final s = d.inSeconds;
    final h = s ~/ 3600;
    final m = (s % 3600) ~/ 60;
    final r = s % 60;
    if (h > 0) {
      return '${h.toString().padLeft(2, '0')}:${m.toString().padLeft(2, '0')}:${r.toString().padLeft(2, '0')}';
    }
    return '${m.toString().padLeft(2, '0')}:${r.toString().padLeft(2, '0')}';
  }

  void _applyBrightnessAccum() {
    while (_brightAccum <= -_brightPxStep) {
      _brightAccum += _brightPxStep;
      _brightnessDim = (_brightnessDim - 0.08).clamp(0.0, 0.85);
      if (mounted) {
        setState(() {});
      }
      _showHud('亮度 ${((1 - _brightnessDim) * 100).round()}%');
    }
    while (_brightAccum >= _brightPxStep) {
      _brightAccum -= _brightPxStep;
      _brightnessDim = (_brightnessDim + 0.08).clamp(0.0, 0.85);
      if (mounted) {
        setState(() {});
      }
      _showHud('亮度 ${((1 - _brightnessDim) * 100).round()}%');
    }
  }

  void _applyVolumeAccum() {
    while (_volAccum <= -_volPxStep) {
      _volAccum += _volPxStep;
      _setVolume(_volume + 0.06);
    }
    while (_volAccum >= _volPxStep) {
      _volAccum -= _volPxStep;
      _setVolume(_volume - 0.06);
    }
  }

  Future<void> _setVolume(double v) async {
    _volume = v.clamp(0.0, 1.0);
    final c = widget.controller;
    if (c.value.isInitialized) {
      await c.setVolume(_volume);
    }
    if (mounted) {
      setState(() {});
      _showHud('音量 ${(_volume * 100).round()}%');
    }
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, c) {
        return Listener(
          behavior: HitTestBehavior.translucent,
          onPointerDown: _onDown,
          onPointerMove: (e) => _onMove(e, c, context),
          onPointerUp: (e) => _finishGesture(context, e.position),
          onPointerCancel: (e) => _finishGesture(context, e.position),
          child: Stack(
            fit: StackFit.expand,
            children: [
              widget.child,
              if (_brightnessDim > 0)
                Positioned.fill(
                  child: IgnorePointer(
                    child: ColoredBox(color: Colors.black.withValues(alpha: _brightnessDim)),
                  ),
                ),
              if (_hud != null)
                Center(
                  child: DecoratedBox(
                    decoration: BoxDecoration(color: Colors.black54, borderRadius: BorderRadius.circular(8)),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 12),
                      child: Text(_hud!, style: const TextStyle(color: Colors.white, fontSize: 16)),
                    ),
                  ),
                ),
            ],
          ),
        );
      },
    );
  }
}
