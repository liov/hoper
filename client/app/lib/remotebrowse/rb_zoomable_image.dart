import 'dart:math' as math;
import 'dart:typed_data';
import 'dart:ui' as ui;

import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';

enum _TwoFingerMode { undecided, pinch, rotate }

/// 全屏单指滑动切预览文件（横/竖），与 [RbZoomableImage] 共用阈值。
abstract final class RbPreviewFileSwipe {
  static const swipeMin = 56.0;
  static const axisRatio = 1.25;

  static void tryAxisSwipe({
    required Offset start,
    required Offset end,
    void Function(bool next)? onVertical,
    void Function(bool next)? onHorizontal,
    bool Function(bool next)? canVertical,
    bool Function(bool next)? canHorizontal,
  }) {
    final d = end - start;
    if (d.distance < swipeMin) {
      return;
    }
    final vert = d.dy.abs() > d.dx.abs() * axisRatio;
    final horiz = d.dx.abs() > d.dy.abs() * axisRatio;
    if (vert && onVertical != null) {
      final next = d.dy < 0;
      if (canVertical?.call(next) ?? true) {
        onVertical(next);
      }
      return;
    }
    if (horiz && onHorizontal != null) {
      final next = d.dx < 0;
      if (canHorizontal?.call(next) ?? true) {
        onHorizontal(next);
      }
    }
  }
}

/// 远程浏览图片：contain、缩放拖动、双指旋转（绘制层旋转，避免缩放变形）；双击还原。
///
/// [onVerticalFileSwipe]：`next==true` 为上滑切下一文件；[onHorizontalFileSwipe]：左滑下一、右滑上一（需未缩放或竖向边界拉出）。
class RbZoomableImage extends StatefulWidget {
  const RbZoomableImage({
    super.key,
    required this.bytes,
    this.fileName = '',
    this.previewMime = '',
    this.onTap,
    this.onVerticalFileSwipe,
    this.canVerticalFileSwipe,
    this.onHorizontalFileSwipe,
    this.canHorizontalFileSwipe,
  });

  final Uint8List bytes;
  final String fileName;
  final String previewMime;
  final VoidCallback? onTap;
  final void Function(bool next)? onVerticalFileSwipe;
  final bool Function(bool next)? canVerticalFileSwipe;
  final void Function(bool next)? onHorizontalFileSwipe;
  final bool Function(bool next)? canHorizontalFileSwipe;

  @override
  State<RbZoomableImage> createState() => _RbZoomableImageState();
}

class _RbZoomableImageState extends State<RbZoomableImage> {
  final _gestureKey = GlobalKey<ExtendedImageGestureState>();
  final _activePointers = <int, Offset>{};
  Offset? _tapDown;
  var _pointers = 0;
  var _rotation = 0.0;
  var _rotationBase = 0.0;
  double? _rotateStartAngle;
  double? _pinchStartSpan;
  double? _pinchStartAngle;
  var _twoFingerAccumSpan = 0.0;
  var _twoFingerAccumAngle = 0.0;
  _TwoFingerMode _twoFingerMode = _TwoFingerMode.undecided;
  Offset? _fileSwipeStart;
  var _edgeOverscroll = 0.0;
  bool? _edgeSwipeNext;

  static const _initialScale = 1.0;
  static const _tapSlop = 40.0;
  static const _fileSwipeMin = 56.0;
  static const _zoomedScaleEpsilon = 0.05;
  /// 旋转灵敏度（<1 更钝）；需累计角度大于捏合才进入旋转模式。
  static const _rotateDamping = 0.32;
  static const _rotateDecideAngle = 0.14;
  static const _rotateDecideSpan = 0.07;
  static const _rotateDominance = 1.65;
  static const _rotateApplyEpsilon = 0.002;

  GestureConfig _gestureConfig(ExtendedImageState state) {
    return GestureConfig(
      minScale: _initialScale,
      animationMinScale: 0.6,
      maxScale: 8.0,
      animationMaxScale: 9.0,
      initialScale: _initialScale,
      speed: 1.0,
      inertialSpeed: 100,
      inPageView: false,
      initialAlignment: InitialAlignment.center,
      cacheGesture: false,
    );
  }

  /// 在 [destinationRect] 内按 contain 绘制后旋转，手势边界仍与 extended_image 一致。
  bool _beforePaintImage(Canvas canvas, Rect destinationRect, ui.Image image, Paint paint) {
    final angle = _rotation;
    if (angle == 0) {
      return false;
    }
    final iw = image.width.toDouble();
    final ih = image.height.toDouble();
    final sourceRect = Alignment.center.inscribe(Size(iw, ih), Offset.zero & Size(iw, ih));
    canvas.save();
    canvas.translate(destinationRect.center.dx, destinationRect.center.dy);
    canvas.rotate(angle);
    final local = Rect.fromCenter(center: Offset.zero, width: destinationRect.width, height: destinationRect.height);
    canvas.drawImageRect(image, sourceRect, local, paint);
    canvas.restore();
    return true;
  }

  double? _twoFingerSpan() {
    if (_activePointers.length < 2) {
      return null;
    }
    final ids = _activePointers.keys.toList()..sort();
    final a = _activePointers[ids[ids.length - 2]]!;
    final b = _activePointers[ids[ids.length - 1]]!;
    return (b - a).distance;
  }

  double? _twoFingerAngle() {
    if (_activePointers.length < 2) {
      return null;
    }
    final ids = _activePointers.keys.toList()..sort();
    final a = _activePointers[ids[ids.length - 2]]!;
    final b = _activePointers[ids[ids.length - 1]]!;
    return math.atan2(b.dy - a.dy, b.dx - a.dx);
  }

  void _resetTwoFingerGesture() {
    _twoFingerMode = _TwoFingerMode.undecided;
    _pinchStartSpan = _twoFingerSpan();
    _pinchStartAngle = _twoFingerAngle();
    _twoFingerAccumSpan = 0.0;
    _twoFingerAccumAngle = 0.0;
    _rotateStartAngle = null;
  }

  void _clearTwoFingerGesture() {
    _twoFingerMode = _TwoFingerMode.undecided;
    _pinchStartSpan = null;
    _pinchStartAngle = null;
    _twoFingerAccumSpan = 0.0;
    _twoFingerAccumAngle = 0.0;
    _rotateStartAngle = null;
    _rotationBase = _rotation;
  }

  void _decideTwoFingerMode() {
    if (_twoFingerMode != _TwoFingerMode.undecided) {
      return;
    }
    final span = _twoFingerAccumSpan;
    final angle = _twoFingerAccumAngle;
    if (span >= _rotateDecideSpan && span > angle * _rotateDominance) {
      _twoFingerMode = _TwoFingerMode.pinch;
      return;
    }
    if (angle >= _rotateDecideAngle && angle > span * _rotateDominance) {
      _twoFingerMode = _TwoFingerMode.rotate;
      _rotateStartAngle = _twoFingerAngle();
      _rotationBase = _rotation;
    }
  }

  void _accumTwoFingerMove() {
    final span = _twoFingerSpan();
    final angle = _twoFingerAngle();
    final baseSpan = _pinchStartSpan;
    final baseAngle = _pinchStartAngle;
    if (span == null || angle == null || baseSpan == null || baseAngle == null || baseSpan < 1) {
      return;
    }
    _twoFingerAccumSpan = (_twoFingerAccumSpan + (span - baseSpan).abs() / baseSpan).clamp(0.0, 2.0);
    var da = angle - baseAngle;
    while (da > math.pi) {
      da -= 2 * math.pi;
    }
    while (da < -math.pi) {
      da += 2 * math.pi;
    }
    _twoFingerAccumAngle = (_twoFingerAccumAngle + da.abs()).clamp(0.0, math.pi * 2);
    _pinchStartSpan = span;
    _pinchStartAngle = angle;
    _decideTwoFingerMode();
  }

  bool get _atInitialScale {
    final scale = _gestureKey.currentState?.gestureDetails?.totalScale ?? _initialScale;
    return (scale - _initialScale).abs() <= _zoomedScaleEpsilon;
  }

  bool get _displayable => rbPreviewBytesDisplayable(
        fileName: widget.fileName,
        mime: widget.previewMime,
        bytesHead: widget.bytes,
      );

  void _onPointerDown(PointerDownEvent e) {
    _pointers++;
    _activePointers[e.pointer] = e.position;
    if (_pointers == 1) {
      _tapDown = e.position;
      if (widget.onVerticalFileSwipe != null || widget.onHorizontalFileSwipe != null) {
        _fileSwipeStart = e.position;
        _edgeOverscroll = 0;
        _edgeSwipeNext = null;
      }
    } else if (_pointers == 2) {
      _fileSwipeStart = null;
      _edgeOverscroll = 0;
      _edgeSwipeNext = null;
      _resetTwoFingerGesture();
    }
  }

  void _trackZoomedEdgePull(PointerMoveEvent e) {
    if (widget.onVerticalFileSwipe == null || _pointers != 1 || _atInitialScale) {
      return;
    }
    final g = _gestureKey.currentState?.gestureDetails;
    if (g == null || !g.computeVerticalBoundary) {
      return;
    }
    final dy = e.delta.dy;
    if (g.boundary.top && dy > 0) {
      _edgeOverscroll += dy;
      _edgeSwipeNext = false;
      return;
    }
    if (g.boundary.bottom && dy < 0) {
      _edgeOverscroll += -dy;
      _edgeSwipeNext = true;
      return;
    }
    final next = _edgeSwipeNext;
    if (next == false && dy > 0) {
      _edgeOverscroll += dy;
    } else if (next == true && dy < 0) {
      _edgeOverscroll += -dy;
    }
  }

  void _tryFileSwipe(Offset end) {
    final start = _fileSwipeStart;
    _fileSwipeStart = null;
    if (start == null) {
      return;
    }
    final onVertical = widget.onVerticalFileSwipe;
    if (onVertical != null) {
      final edgeNext = _edgeSwipeNext;
      if (_edgeOverscroll >= _fileSwipeMin && edgeNext != null) {
        _edgeOverscroll = 0;
        _edgeSwipeNext = null;
        if (widget.canVerticalFileSwipe?.call(edgeNext) ?? true) {
          onVertical(edgeNext);
        }
        return;
      }
    }
    _edgeOverscroll = 0;
    _edgeSwipeNext = null;
    if (!_atInitialScale) {
      return;
    }
    RbPreviewFileSwipe.tryAxisSwipe(
      start: start,
      end: end,
      onVertical: onVertical,
      onHorizontal: widget.onHorizontalFileSwipe,
      canVertical: widget.canVerticalFileSwipe,
      canHorizontal: widget.canHorizontalFileSwipe,
    );
  }

  void _onPointerMove(PointerMoveEvent e) {
    _activePointers[e.pointer] = e.position;
    _trackZoomedEdgePull(e);
    if (_activePointers.length < 2) {
      return;
    }
    _accumTwoFingerMove();
    if (_twoFingerMode != _TwoFingerMode.rotate) {
      return;
    }
    final start = _rotateStartAngle;
    if (start == null) {
      return;
    }
    final cur = _twoFingerAngle();
    if (cur == null) {
      return;
    }
    var delta = cur - start;
    while (delta > math.pi) {
      delta -= 2 * math.pi;
    }
    while (delta < -math.pi) {
      delta += 2 * math.pi;
    }
    final next = _rotationBase + delta * _rotateDamping;
    if ((next - _rotation).abs() < _rotateApplyEpsilon) {
      return;
    }
    setState(() => _rotation = next);
  }

  void _onPointerUp(PointerUpEvent e) {
    _activePointers.remove(e.pointer);
    if (_pointers > 0) {
      _pointers--;
    }
    if (_pointers < 2) {
      _clearTwoFingerGesture();
    }
    if (_pointers != 0) {
      return;
    }
    _tryFileSwipe(e.position);
    final down = _tapDown;
    _tapDown = null;
    final onTap = widget.onTap;
    if (down == null || onTap == null) {
      return;
    }
    if ((e.position - down).distance > _tapSlop) {
      return;
    }
    onTap();
  }

  void _onPointerCancel(PointerCancelEvent e) {
    _activePointers.remove(e.pointer);
    if (_pointers > 0) {
      _pointers--;
    }
    if (_pointers < 2) {
      _clearTwoFingerGesture();
    }
    if (_pointers == 0) {
      _fileSwipeStart = null;
      _edgeOverscroll = 0;
      _edgeSwipeNext = null;
      _tapDown = null;
    }
  }

  void _onDoubleTap(ExtendedImageGestureState state) {
    setState(() {
      _rotation = 0;
      _rotationBase = 0;
      _clearTwoFingerGesture();
    });
    state.reset();
  }

  @override
  Widget build(BuildContext context) {
    if (!_displayable || widget.bytes.isEmpty) {
      final hint = rbImageNeedsAgentPreviewWebp(widget.fileName) ? '需 Agent 转码预览' : '图片无法显示';
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.broken_image_outlined, color: Colors.white54, size: 48),
            const SizedBox(height: 8),
            Text(hint, textAlign: TextAlign.center, style: const TextStyle(color: Colors.white54, fontSize: 13)),
          ],
        ),
      );
    }
    final raster = widget.bytes;
    return Listener(
      behavior: HitTestBehavior.translucent,
      onPointerDown: _onPointerDown,
      onPointerMove: _onPointerMove,
      onPointerUp: _onPointerUp,
      onPointerCancel: _onPointerCancel,
      child: ExtendedImage.memory(
        raster,
        fit: BoxFit.contain,
        mode: ExtendedImageMode.gesture,
        extendedImageGestureKey: _gestureKey,
        gaplessPlayback: true,
        filterQuality: FilterQuality.medium,
        initGestureConfigHandler: _gestureConfig,
        beforePaintImage: _beforePaintImage,
        onDoubleTap: _onDoubleTap,
        loadStateChanged: (state) {
          if (state.extendedImageLoadState != LoadState.failed) {
            return null;
          }
          return const Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(Icons.broken_image_outlined, color: Colors.white54, size: 48),
                SizedBox(height: 8),
                Text('图片解码失败', textAlign: TextAlign.center, style: TextStyle(color: Colors.white54, fontSize: 13)),
              ],
            ),
          );
        },
      ),
    );
  }
}

/// 缩略图占位等无缩放层：横/竖滑切文件 + 轻点。
class RbPreviewFileSwipeListener extends StatefulWidget {
  const RbPreviewFileSwipeListener({
    super.key,
    required this.child,
    this.onTap,
    this.onVerticalFileSwipe,
    this.canVerticalFileSwipe,
    this.onHorizontalFileSwipe,
    this.canHorizontalFileSwipe,
  });

  final Widget child;
  final VoidCallback? onTap;
  final void Function(bool next)? onVerticalFileSwipe;
  final bool Function(bool next)? canVerticalFileSwipe;
  final void Function(bool next)? onHorizontalFileSwipe;
  final bool Function(bool next)? canHorizontalFileSwipe;

  @override
  State<RbPreviewFileSwipeListener> createState() => _RbPreviewFileSwipeListenerState();
}

class _RbPreviewFileSwipeListenerState extends State<RbPreviewFileSwipeListener> {
  static const _tapSlop = 40.0;
  Offset? _down;

  @override
  Widget build(BuildContext context) {
    return Listener(
      behavior: HitTestBehavior.translucent,
      onPointerDown: (e) => _down = e.position,
      onPointerUp: (e) {
        final down = _down;
        _down = null;
        if (down == null) {
          return;
        }
        final moved = (e.position - down).distance;
        if (moved > _tapSlop) {
          RbPreviewFileSwipe.tryAxisSwipe(
            start: down,
            end: e.position,
            onVertical: widget.onVerticalFileSwipe,
            onHorizontal: widget.onHorizontalFileSwipe,
            canVertical: widget.canVerticalFileSwipe,
            canHorizontal: widget.canHorizontalFileSwipe,
          );
          return;
        }
        widget.onTap?.call();
      },
      onPointerCancel: (_) => _down = null,
      child: widget.child,
    );
  }
}
