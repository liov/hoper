import 'dart:math' as math;

import 'package:flutter/material.dart';

/// 左缘上下滑切文件的起始 x，避开 Android/iOS 系统返回手势区。
double rbPreviewLeftEdgeInset(BuildContext context) {
  final mq = MediaQuery.of(context);
  final reserved = math.max(mq.systemGestureInsets.left, mq.padding.left);
  const floor = 28.0;
  const pad = 12.0;
  return math.max(floor, reserved + pad);
}

/// 预览页手势：左右边缘上下滑切文件；底部横滑切同级目录。
class RbPreviewEdgeSwipeLayer extends StatefulWidget {
  const RbPreviewEdgeSwipeLayer({
    super.key,
    required this.child,
    required this.onEdgeVerticalSwipe,
    this.onHorizontalFileSwipe,
    this.onBottomHorizontalSwipe,
    this.edgeWidth = 56,
    this.bottomStripHeight = 72,
    this.leftEdgeVerticalEnabled = true,
    this.edgeVerticalSwipeEnabled = true,
    this.canTriggerVerticalSwipe,
    this.canTriggerHorizontalFileSwipe,
    this.horizontalFileSwipeOnEdgeOnly = false,
  });

  final Widget child;
  final void Function(bool next) onEdgeVerticalSwipe;
  /// 全屏横滑切预览文件：左滑下一、右滑上一。
  final void Function(bool next)? onHorizontalFileSwipe;
  final void Function(bool next)? onBottomHorizontalSwipe;
  final double edgeWidth;
  final double bottomStripHeight;
  /// 为 false 时左缘不触发上下切文件（避免与系统返回/竖向滚动冲突）。
  final bool leftEdgeVerticalEnabled;
  /// 为 false 时不启用左右缘上下滑（由子组件全屏手势接管，如图片）。
  final bool edgeVerticalSwipeEnabled;
  /// `next==true` 表示上滑切下一文件；返回 false 则忽略本次边缘手势。
  final bool Function(bool next)? canTriggerVerticalSwipe;
  final bool Function(bool next)? canTriggerHorizontalFileSwipe;
  /// 为 true 时横滑切文件须在左右边缘发起，避免与 Markdown 代码块等横向滚动冲突。
  final bool horizontalFileSwipeOnEdgeOnly;

  @override
  State<RbPreviewEdgeSwipeLayer> createState() => _RbPreviewEdgeSwipeLayerState();
}

class _RbPreviewEdgeSwipeLayerState extends State<RbPreviewEdgeSwipeLayer> {
  static const _swipeMin = 56.0;
  static const _axisRatio = 1.25;

  Offset? _start;

  bool _onVerticalEdge(BuildContext context, double x, double width) {
    final inset = rbPreviewLeftEdgeInset(context);
    final left = widget.leftEdgeVerticalEnabled && x >= inset && x < inset + widget.edgeWidth;
    final right = x > width - widget.edgeWidth;
    return left || right;
  }

  bool _onHorizontalEdge(BuildContext context, double x, double width) {
    final inset = rbPreviewLeftEdgeInset(context);
    final left = x >= inset && x < inset + widget.edgeWidth;
    final right = x > width - widget.edgeWidth;
    return left || right;
  }

  bool _inBottomStrip(double y, double height) => y > height - widget.bottomStripHeight;

  void _onDown(PointerDownEvent e) => _start = e.position;

  void _onUp(PointerUpEvent e) {
    final start = _start;
    _start = null;
    if (start == null) {
      return;
    }
    final d = e.position - start;
    if (d.distance < _swipeMin) {
      return;
    }
    final size = MediaQuery.sizeOf(context);
    final w = size.width;
    final h = size.height;
    final horiz = d.dx.abs() > d.dy.abs() * _axisRatio;
    final vert = d.dy.abs() > d.dx.abs() * _axisRatio;
    final bottomSwipe = widget.onBottomHorizontalSwipe;
    if (bottomSwipe != null && horiz && _inBottomStrip(start.dy, h)) {
      bottomSwipe(d.dx < 0);
      return;
    }
    final horizFile = widget.onHorizontalFileSwipe;
    if (horizFile != null && horiz) {
      if (widget.horizontalFileSwipeOnEdgeOnly && !_onHorizontalEdge(context, start.dx, w)) {
        return;
      }
      final next = d.dx < 0;
      if (widget.canTriggerHorizontalFileSwipe?.call(next) ?? true) {
        horizFile(next);
        return;
      }
    }
    if (vert && widget.edgeVerticalSwipeEnabled && _onVerticalEdge(context, start.dx, w)) {
      final next = d.dy < 0;
      if (widget.canTriggerVerticalSwipe?.call(next) ?? true) {
        widget.onEdgeVerticalSwipe(next);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Listener(
      behavior: HitTestBehavior.translucent,
      onPointerDown: _onDown,
      onPointerUp: _onUp,
      onPointerCancel: (_) => _start = null,
      child: widget.child,
    );
  }
}
