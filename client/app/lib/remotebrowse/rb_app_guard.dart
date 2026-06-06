import 'package:app/global/logger.dart';
import 'package:flutter/foundation.dart';

/// 全局兜底：记录未捕获异常，尽量避免进程直接退出。
void installRbAppCrashGuard() {
  final prevFlutter = FlutterError.onError;
  FlutterError.onError = (FlutterErrorDetails details) {
    logZoneError(details.exception, details.stack ?? StackTrace.empty);
    if (kDebugMode) {
      if (prevFlutter != null) {
        prevFlutter(details);
      } else {
        FlutterError.presentError(details);
      }
    }
  };

  PlatformDispatcher.instance.onError = (Object error, StackTrace stack) {
    logZoneError(error, stack);
    return true;
  };
}
