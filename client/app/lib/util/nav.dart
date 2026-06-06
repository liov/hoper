import 'package:flutter/material.dart';

final rootNavigatorKey = GlobalKey<NavigatorState>();

/// 全局导航与弹层（原 Get 导航 API）。
abstract final class AppNavigator {
  static NavigatorState? get nav => rootNavigatorKey.currentState;
  static BuildContext? get context => rootNavigatorKey.currentContext;

  static Future<T?> push<T>(Widget page) {
    final n = nav;
    if (n == null) {
      return Future.value();
    }
    return n.push<T>(MaterialPageRoute(builder: (_) => page));
  }

  static Future<T?> pushNamed<T>(String name, {Object? arguments}) {
    final n = nav;
    if (n == null) {
      return Future.value();
    }
    return n.pushNamed<T>(name, arguments: arguments);
  }

  static void pop<T>([T? result]) {
    final n = nav;
    if (n != null && n.canPop()) {
      n.pop(result);
    }
  }

  static void offAllNamed(String name, {Object? arguments}) {
    nav?.pushNamedAndRemoveUntil(name, (_) => false, arguments: arguments);
  }

  static Future<T?> dialog<T>(Widget child) {
    final ctx = context;
    if (ctx == null) {
      return Future.value();
    }
    return showDialog<T>(context: ctx, builder: (_) => child);
  }

  static void snackbar(String title, [String? message, Duration duration = const Duration(seconds: 3)]) {
    final ctx = context;
    if (ctx == null) {
      return;
    }
    final text = message == null || message.isEmpty ? title : '$title: $message';
    ScaffoldMessenger.of(ctx).showSnackBar(
      SnackBar(content: Text(text), duration: duration, behavior: SnackBarBehavior.floating),
    );
  }

  static void rawSnackbar(String message) => snackbar(message);

  static OverlayEntry? _overlay;

  static Future<void> showOverlay({
    required Widget loadingWidget,
    required Future<void> Function() asyncFunction,
  }) async {
    final overlay = nav?.overlay;
    if (overlay == null) {
      await asyncFunction();
      return;
    }
    _overlay?.remove();
    _overlay = OverlayEntry(builder: (_) => Positioned.fill(child: IgnorePointer(child: loadingWidget)));
    overlay.insert(_overlay!);
    try {
      await asyncFunction();
    } finally {
      _overlay?.remove();
      _overlay = null;
    }
  }

  static double get width {
    final ctx = context;
    if (ctx == null) {
      return 0;
    }
    return MediaQuery.sizeOf(ctx).width;
  }

  static double get height {
    final ctx = context;
    if (ctx == null) {
      return 0;
    }
    return MediaQuery.sizeOf(ctx).height;
  }

  static ThemeData get theme => Theme.of(context!);
}
