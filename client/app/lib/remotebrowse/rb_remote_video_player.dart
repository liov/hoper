import 'dart:async';

import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/remotebrowse/rb_transcode.dart';
import 'package:app/remotebrowse/rb_video_gesture_layer.dart';
import 'package:app/remotebrowse/rb_video_timeline.dart';
import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';

/// 远程预览视频播放器：全屏画面、单击显隐顶/底工具条、进度条、播放/暂停、倍速。
class RbRemoteVideoPlayer extends StatefulWidget {
  const RbRemoteVideoPlayer({
    super.key,
    required this.controller,
    this.topChrome,
    this.folderStrip,
    this.bottomBarLeading,
    this.bottomBarTrailing,
    this.hideChromeInitially = true,
    this.autoHideChrome = false,
    this.holdPlaybackSpeed = 2.0,
    this.holdSpeedStep = 0.1,
    this.holdSpeedMin = 0.5,
    this.holdSpeedMax = 4.0,
    this.holdSpeedDragWidthFraction = 0.20,
    this.seekSpanPerWidth = 3.0,
    this.seekMinMsPerWidth = 150000,
    this.onEdgeVerticalSwipe,
    this.onBottomHorizontalSwipe,
    this.transcodePresetId,
    this.transcodeAgentAvailable = false,
    this.onTranscodePreset,
    this.onHlsSeek,
    this.timelineOffsetMs = 0,
    this.timelineTotalMs = 0,
    this.playbackAnchorMs = 0,
    this.timelineUsePlayClock = false,
    this.timelinePausedAbsPosMs,
    this.timelinePlayStartedAt,
    this.hintLabel = '',
  });

  final VideoPlayerController controller;
  /// 画面中间区域上/下滑：切换预览文件（上=下一个，下=上一个）；左右条带为亮度/音量。
  final void Function(bool next)? onEdgeVerticalSwipe;
  /// 底部横滑：切换同级目录（左滑下一个，右滑上一个）。
  final void Function(bool next)? onBottomHorizontalSwipe;
  final String? transcodePresetId;
  final bool transcodeAgentAvailable;
  final ValueChanged<String>? onTranscodePreset;
  /// 转码 HLS 无法用 ExoPlayer 可靠 seek 时，由上层按 `begin_ms` 重载流。
  final Future<void> Function(Duration target, {bool playAfter})? onHlsSeek;
  /// 裁剪 HLS 列表时，UI/手势使用的全片时间轴偏移（毫秒）。
  final int timelineOffsetMs;
  /// 全片时长（毫秒）；转码切档后播放器 duration 不可靠时用目录元数据。
  final int timelineTotalMs;
  /// 连续流转码起播锚点（全片毫秒）；与 [timelineOffsetMs] 互斥。
  final int playbackAnchorMs;
  final bool timelineUsePlayClock;
  /// 切档/跳转缓冲时由上层提供的全片位置（毫秒），避免进度条归零。
  final int? timelinePausedAbsPosMs;
  final DateTime? timelinePlayStartedAt;
  /// 跳转/切档等状态的顶部浮层文案（不占布局）。
  final String hintLabel;
  /// 长按开始时的初始倍速（默认 2.0）。
  final double holdPlaybackSpeed;
  final double holdSpeedStep;
  final double holdSpeedMin;
  final double holdSpeedMax;
  /// 长按倍速拖动：屏宽该比例 ≈ 从 [holdSpeedMin] 滑到 [holdSpeedMax]（与横向快进无关）。
  final double holdSpeedDragWidthFraction;
  /// 横向快进灵敏度（传给 [RbVideoGestureLayer]，越小越不灵敏）。
  final double seekSpanPerWidth;
  final int seekMinMsPerWidth;
  /// 顶部栏（返回、标题、切文件等），与播放控制一同显隐。
  final Widget? topChrome;
  /// 底部控制条上方（如当前目录缩略图条）。
  final Widget? folderStrip;
  /// 底部控制栏左侧（如上一/下一文件）。
  final Widget? bottomBarLeading;
  /// 底部控制栏右侧附加控件（如删除）。
  final Widget? bottomBarTrailing;
  final bool hideChromeInitially;
  final bool autoHideChrome;

  @override
  State<RbRemoteVideoPlayer> createState() => _RbRemoteVideoPlayerState();
}

class _RbRemoteVideoPlayerState extends State<RbRemoteVideoPlayer> {
  static const _playbackSpeeds = [0.5, 1.0, 1.25, 1.5, 2.0, 3.0, 4.0];

  late var _showControls = !widget.hideChromeInitially;
  var _dragging = false;
  var _speed = 1.0;
  var _holdSpeedActive = false;
  var _holdSpeed = 2.0;
  var _holdAccumDx = 0.0;
  var _holdAccumDy = 0.0;
  var _holdSpeedAnchor = 2.0;
  var _holdDragSpanPx = 320.0;
  var _ignoreNextTap = false;
  var _tickRebuildScheduled = false;
  var _wasPlaying = false;
  DateTime? _anchorPlayStartedAt;
  int? _anchorPausedAbsPosMs;
  Timer? _hideTimer;
  Timer? _timelineTickTimer;

  VideoPlayerController get _c => widget.controller;

  @override
  void initState() {
    super.initState();
    _c.addListener(_onTick);
    if (_showControls && widget.autoHideChrome) {
      _armHideControls();
    }
    if (widget.timelineUsePlayClock) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) {
          _syncTimelinePlayClock(_c.value);
        }
      });
    }
  }

  @override
  void didUpdateWidget(covariant RbRemoteVideoPlayer oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.controller != widget.controller) {
      try {
        oldWidget.controller.removeListener(_onTick);
      } catch (_) {}
      try {
        _c.addListener(_onTick);
      } catch (_) {}
      _wasPlaying = false;
      _anchorPlayStartedAt = null;
      _anchorPausedAbsPosMs = null;
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) {
          _syncTimelinePlayClock(_c.value);
        }
      });
    }
    if (oldWidget.transcodePresetId != widget.transcodePresetId || oldWidget.timelineUsePlayClock != widget.timelineUsePlayClock) {
      setState(() {});
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) {
          _syncTimelinePlayClock(_c.value);
        }
      });
    }
  }

  @override
  void dispose() {
    _timelineTickTimer?.cancel();
    _hideTimer?.cancel();
    _holdSpeedActive = false;
    _holdAccumDx = 0;
    _holdAccumDy = 0;
    try {
      _c.removeListener(_onTick);
    } catch (_) {}
    super.dispose();
  }

  void _scheduleRebuild() {
    if (!mounted || _tickRebuildScheduled) {
      return;
    }
    _tickRebuildScheduled = true;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _tickRebuildScheduled = false;
      if (mounted) {
        setState(() {});
      }
    });
  }

  void _onTick() {
    VideoPlayerValue v;
    try {
      v = _c.value;
    } catch (_) {
      return;
    }
    _syncTimelinePlayClock(v);
    _scheduleRebuild();
  }

  void _syncTimelinePlayClock(VideoPlayerValue v) {
    if (!widget.timelineUsePlayClock || !v.isInitialized) {
      _timelineTickTimer?.cancel();
      _timelineTickTimer = null;
      return;
    }
    final playing = v.isPlaying;
    if (playing && !_wasPlaying) {
      _anchorPlayStartedAt = DateTime.now();
      _anchorPausedAbsPosMs = null;
    } else if (!playing && _wasPlaying) {
      _anchorPausedAbsPosMs = _timelinePosMsFrom(v);
      _anchorPlayStartedAt = null;
    }
    _wasPlaying = playing;
    if (playing) {
      if (_timelineTickTimer != null) {
        return;
      }
      _timelineTickTimer = Timer.periodic(const Duration(milliseconds: 200), (_) {
        if (!mounted) {
          return;
        }
        VideoPlayerValue vv;
        try {
          vv = _c.value;
        } catch (_) {
          return;
        }
        if (!widget.timelineUsePlayClock || !vv.isInitialized || !vv.isPlaying) {
          _syncTimelinePlayClock(vv);
          return;
        }
        setState(() {});
      });
    } else {
      _timelineTickTimer?.cancel();
      _timelineTickTimer = null;
    }
  }

  int _timelinePosMsFrom(VideoPlayerValue v) {
    return rbVideoAbsPosMs(
      playerPosMs: v.position.inMilliseconds,
      timelineOffsetMs: widget.timelineOffsetMs,
      playbackAnchorMs: widget.playbackAnchorMs,
      isPlaying: v.isPlaying,
      playStartedAt: _anchorPlayStartedAt ?? widget.timelinePlayStartedAt,
      pausedAbsPosMs: _anchorPausedAbsPosMs ?? widget.timelinePausedAbsPosMs,
      playbackSpeed: v.playbackSpeed,
    );
  }

  void _armHideControls() {
    _hideTimer?.cancel();
    if (!_showControls || !widget.autoHideChrome) {
      return;
    }
    _hideTimer = Timer(const Duration(seconds: 4), () {
      if (mounted && _c.value.isPlaying && !_dragging) {
        setState(() => _showControls = false);
      }
    });
  }

  void _toggleControls() {
    if (!mounted) {
      return;
    }
    setState(() => _showControls = !_showControls);
    if (_showControls) {
      _armHideControls();
    } else {
      _hideTimer?.cancel();
    }
  }

  Widget _chromeGradient({required bool top}) {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: top ? Alignment.topCenter : Alignment.bottomCenter,
          end: top ? Alignment.bottomCenter : Alignment.topCenter,
          colors: [Colors.black.withValues(alpha: 0.72), Colors.transparent],
        ),
      ),
    );
  }

  Future<void> _togglePlay() async {
    final v = _c.value;
    if (!v.isInitialized) {
      return;
    }
    try {
      if (v.position >= v.duration && v.duration > Duration.zero) {
        final hlsSeek = widget.onHlsSeek;
        if (hlsSeek != null) {
          await hlsSeek(Duration.zero, playAfter: true);
          return;
        }
        await _c.seekTo(Duration.zero);
      }
      if (v.isPlaying) {
        await _c.pause();
      } else {
        await _c.play();
      }
    } catch (_) {
      return;
    }
    if (mounted) {
      setState(() => _showControls = true);
      _armHideControls();
    }
  }

  int get _timelineTotalMs {
    if (widget.timelineTotalMs > 0) {
      return widget.timelineTotalMs;
    }
    final d = _c.value.duration.inMilliseconds;
    return widget.timelineOffsetMs + d;
  }

  int get _timelinePosMs {
    try {
      return _timelinePosMsFrom(_c.value);
    } catch (_) {
      return 0;
    }
  }

  Future<void> _seekToFraction(double f) async {
    final total = _timelineTotalMs;
    if (total <= 0) {
      return;
    }
    final absMs = (total * f.clamp(0.0, 1.0)).round();
    final hlsSeek = widget.onHlsSeek;
    if (hlsSeek != null) {
      await hlsSeek(Duration(milliseconds: absMs), playAfter: _c.value.isPlaying);
      return;
    }
    try {
      await _c.seekTo(Duration(milliseconds: absMs - widget.timelineOffsetMs));
    } catch (_) {}
  }

  Future<void> _setSpeed(double speed) async {
    if (_holdSpeedActive) {
      return;
    }
    _speed = speed;
    if (!_c.value.isInitialized) {
      if (mounted) {
        setState(() {});
      }
      return;
    }
    try {
      await _c.setPlaybackSpeed(_speed);
    } catch (_) {
      return;
    }
    if (mounted) {
      setState(() {});
      _armHideControls();
    }
  }

  String _fmtSpeed(double s) => s == s.roundToDouble() ? '${s.toStringAsFixed(0)}x' : '${s.toStringAsFixed(1)}x';

  Future<void> _applyHoldSpeed(double speed) async {
    _holdSpeed = speed.clamp(widget.holdSpeedMin, widget.holdSpeedMax);
    if (!_c.value.isInitialized) {
      return;
    }
    try {
      await _c.setPlaybackSpeed(_holdSpeed);
    } catch (_) {
      return;
    }
    if (mounted) {
      setState(() {});
    }
  }

  Future<void> _startHoldSpeed() async {
    final v = _c.value;
    if (!v.isInitialized || _holdSpeedActive) {
      return;
    }
    _hideTimer?.cancel();
    _holdSpeedActive = true;
    _holdAccumDx = 0;
    _holdAccumDy = 0;
    _holdSpeedAnchor = widget.holdPlaybackSpeed;
    if (!v.isPlaying) {
      await _c.play();
    }
    await _applyHoldSpeed(_holdSpeedAnchor);
  }

  Future<void> _endHoldSpeed() async {
    if (!_holdSpeedActive) {
      return;
    }
    _holdSpeedActive = false;
    _holdAccumDx = 0;
    _holdAccumDy = 0;
    _ignoreNextTap = true;
    if (_c.value.isInitialized) {
      try {
        await _c.setPlaybackSpeed(_speed);
      } catch (_) {}
    }
    if (mounted) {
      setState(() {});
      _armHideControls();
    }
  }

  void _onHoldPointerMove(PointerMoveEvent e) {
    if (!_holdSpeedActive) {
      return;
    }
    _holdAccumDx += e.delta.dx;
    _holdAccumDy += e.delta.dy;
    final span = _holdDragSpanPx.clamp(80.0, 480.0);
    final range = widget.holdSpeedMax - widget.holdSpeedMin;
    final delta = (-_holdAccumDy + _holdAccumDx) / span * range;
    final target = (_holdSpeedAnchor + delta).clamp(widget.holdSpeedMin, widget.holdSpeedMax);
    if ((target - _holdSpeed).abs() < 0.03) {
      return;
    }
    unawaited(_applyHoldSpeed(target));
  }

  String _fmt(Duration d) {
    final s = d.inSeconds;
    final m = s ~/ 60;
    final r = s % 60;
    return '${m.toString().padLeft(2, '0')}:${r.toString().padLeft(2, '0')}';
  }

  @override
  Widget build(BuildContext context) {
    VideoPlayerValue v;
    try {
      v = _c.value;
    } catch (_) {
      return const Center(child: CircularProgressIndicator(strokeWidth: 2.5, color: Colors.white70));
    }
    if (!v.isInitialized) {
      return Stack(
        fit: StackFit.expand,
        children: [
          const Center(
            child: SizedBox(
              width: 36,
              height: 36,
              child: CircularProgressIndicator(strokeWidth: 2.5, color: Colors.white70),
            ),
          ),
          if (widget.hintLabel.isNotEmpty)
            Positioned(
              top: 0,
              left: 0,
              right: 0,
              child: RbStatusOverlay(message: widget.hintLabel, onDark: true, showProgress: false),
            ),
        ],
      );
    }
    final absDurMs = _timelineTotalMs;
    final absPosMs = absDurMs > 0 ? _timelinePosMs.clamp(0, absDurMs) : _timelinePosMs;
    final progress = absDurMs > 0 ? absPosMs / absDurMs : 0.0;
    final ar = v.aspectRatio > 0 ? v.aspectRatio : 16 / 9;
    final sz = v.size;
    final videoW = sz.width > 0 ? sz.width : ar;
    final videoH = sz.height > 0 ? sz.height : 1.0;
    final holdLabel = _fmtSpeed(_holdSpeed);
    return LayoutBuilder(
      builder: (context, constraints) {
        final w = constraints.maxWidth.isFinite && constraints.maxWidth > 0
            ? constraints.maxWidth
            : MediaQuery.sizeOf(context).width;
        _holdDragSpanPx = w * widget.holdSpeedDragWidthFraction;
        return Listener(
          onPointerMove: _onHoldPointerMove,
          child: GestureDetector(
        onTap: () {
          if (_ignoreNextTap) {
            _ignoreNextTap = false;
            return;
          }
          if (!_holdSpeedActive) {
            _toggleControls();
          }
        },
        onLongPressStart: (_) => unawaited(_startHoldSpeed()),
        onLongPressEnd: (_) => unawaited(_endHoldSpeed()),
        onLongPressCancel: () => unawaited(_endHoldSpeed()),
        behavior: HitTestBehavior.opaque,
        child: RbVideoGestureLayer(
          controller: _c,
          onEdgeVerticalSwipe: widget.onEdgeVerticalSwipe,
          onBottomHorizontalSwipe: widget.onBottomHorizontalSwipe,
          onHlsSeek: widget.onHlsSeek,
          timelineOffsetMs: widget.timelineOffsetMs,
          timelineTotalMs: widget.timelineTotalMs,
          playbackAnchorMs: widget.playbackAnchorMs,
          playStartedAt: _anchorPlayStartedAt ?? widget.timelinePlayStartedAt,
          pausedAbsPosMs: _anchorPausedAbsPosMs ?? widget.timelinePausedAbsPosMs,
          controlsVisible: _showControls,
          holdSpeedActive: _holdSpeedActive,
          seekSpanPerWidth: widget.seekSpanPerWidth,
          seekMinMsPerWidth: widget.seekMinMsPerWidth,
          child: Stack(
            fit: StackFit.expand,
            children: [
              const ColoredBox(color: Colors.black),
              Center(
                child: FittedBox(
                  fit: BoxFit.contain,
                  child: SizedBox(width: videoW, height: videoH, child: VideoPlayer(_c, key: ObjectKey(_c))),
                ),
              ),
          if (widget.hintLabel.isNotEmpty)
            Positioned(
              top: 0,
              left: 0,
              right: 0,
              child: IgnorePointer(child: RbStatusOverlay(message: widget.hintLabel, onDark: true, showProgress: false)),
            ),
          if (_holdSpeedActive)
            Center(
              child: DecoratedBox(
                decoration: BoxDecoration(color: Colors.black54, borderRadius: BorderRadius.circular(8)),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                  child: Column(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text('$holdLabel 快进中', style: const TextStyle(color: Colors.white, fontSize: 17)),
                      const SizedBox(height: 4),
                      const Text('↑/→ 加速  ↓/← 减速', style: TextStyle(color: Colors.white70, fontSize: 12)),
                    ],
                  ),
                ),
              ),
            ),
          if (_showControls && widget.topChrome != null)
            Positioned(
              left: 0,
              right: 0,
              top: 0,
              child: Stack(
                children: [
                  Positioned(left: 0, right: 0, top: 0, height: 120, child: _chromeGradient(top: true)),
                  SafeArea(bottom: false, child: widget.topChrome!),
                ],
              ),
            ),
          if (_showControls)
            Positioned(
              left: 0,
              right: 0,
              bottom: 0,
              child: DecoratedBox(
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        begin: Alignment.bottomCenter,
                        end: Alignment.topCenter,
                        colors: [Colors.black.withValues(alpha: 0.75), Colors.transparent],
                      ),
                    ),
                    child: SafeArea(
                      top: false,
                      child: Padding(
                        padding: const EdgeInsets.fromLTRB(12, 8, 12, 8),
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            if (widget.folderStrip != null) widget.folderStrip!,
                            SliderTheme(
                              data: SliderTheme.of(context).copyWith(
                                trackHeight: 2,
                                thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 6),
                                overlayShape: const RoundSliderOverlayShape(overlayRadius: 12),
                              ),
                              child: Slider(
                                value: _dragging ? _dragValue ?? progress : progress.clamp(0.0, 1.0),
                                onChangeStart: (_) {
                                  if (!mounted) {
                                    return;
                                  }
                                  setState(() {
                                    _dragging = true;
                                    _dragValue = progress;
                                  });
                                  _hideTimer?.cancel();
                                },
                                onChanged: (x) {
                                  if (mounted) {
                                    setState(() => _dragValue = x);
                                  }
                                },
                                onChangeEnd: (x) async {
                                  if (mounted) {
                                    setState(() => _dragging = false);
                                  }
                                  await _seekToFraction(x);
                                  _dragValue = null;
                                  _armHideControls();
                                },
                              ),
                            ),
                            Row(
                              children: [
                                IconButton(
                                  visualDensity: VisualDensity.compact,
                                  padding: EdgeInsets.zero,
                                  constraints: const BoxConstraints(minWidth: 40, minHeight: 40),
                                  onPressed: () => unawaited(_togglePlay()),
                                  icon: Icon(v.isPlaying ? Icons.pause : Icons.play_arrow, color: Colors.white),
                                ),
                                if (widget.bottomBarLeading != null) widget.bottomBarLeading!,
                                Expanded(
                                  child: Text(
                                    '${_fmt(Duration(milliseconds: absPosMs))} / ${_fmt(Duration(milliseconds: absDurMs))}',
                                    maxLines: 1,
                                    overflow: TextOverflow.ellipsis,
                                    style: const TextStyle(color: Colors.white70, fontSize: 12),
                                  ),
                                ),
                                Flexible(
                                  fit: FlexFit.loose,
                                  child: FittedBox(
                                    fit: BoxFit.scaleDown,
                                    alignment: Alignment.centerRight,
                                    child: Row(
                                      mainAxisSize: MainAxisSize.min,
                                      children: [
                                        if (widget.bottomBarTrailing != null) widget.bottomBarTrailing!,
                                        if (widget.onTranscodePreset != null) _transcodeMenu(),
                                        _speedMenu(holdLabel),
                                      ],
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
            ),
            ],
          ),
        ),
        ),
      );
      },
    );
  }

  double? _dragValue;

  Widget _transcodeMenu() {
    final cur = widget.transcodePresetId ?? RbTranscodePreset.source.id;
    final label = RbTranscodePreset.all.firstWhere((p) => p.id == cur, orElse: () => RbTranscodePreset.source).label;
    return PopupMenuButton<String>(
      tooltip: '清晰度：$label',
      initialValue: cur,
      padding: EdgeInsets.zero,
      onSelected: widget.onTranscodePreset,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 6),
        child: Text(label, style: const TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.w500)),
      ),
      itemBuilder: (ctx) => [
        for (final p in RbTranscodePreset.all)
          PopupMenuItem<String>(
            value: p.id,
            enabled: !p.agentOnly || widget.transcodeAgentAvailable,
            child: Text(p.label),
          ),
      ],
    );
  }

  Widget _speedMenu(String holdLabel) {
    final label = _holdSpeedActive ? holdLabel : _fmtSpeed(_speed);
    return PopupMenuButton<double>(
      tooltip: '倍速：$label',
      initialValue: _holdSpeedActive ? null : _speed,
      enabled: !_holdSpeedActive,
      padding: EdgeInsets.zero,
      onSelected: (s) => unawaited(_setSpeed(s)),
      icon: const Icon(Icons.speed, color: Colors.white, size: 22),
      itemBuilder: (ctx) => [
        for (final s in _playbackSpeeds)
          PopupMenuItem<double>(value: s, child: Text(_fmtSpeed(s))),
      ],
    );
  }
}
