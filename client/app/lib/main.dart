import 'dart:async';

import 'package:flutter/foundation.dart';

import 'package:app/app.dart';
import 'package:app/remotebrowse/rb_app_guard.dart';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:app/global/logger.dart';
import 'package:app/global/state.dart';
import 'package:app/global/state/app.dart';

Future<void> main() async {
  runZonedGuarded(
    () async {
      WidgetsFlutterBinding.ensureInitialized();
      installRbAppCrashGuard();
      ErrorWidget.builder = (FlutterErrorDetails flutterErrorDetails) {
        globalService.logger.fine(flutterErrorDetails.toString());
        return const Center(child: Text("找不到页面"));
      };
      globalService.logger.fine('runZonedGuarded');
      AppState.isDebug = kDebugMode;
      globalService.logger.fine("${AppState.isDebug}");
      runApp(ProviderScope(child: AppRoot(key: AppRoot.restartKey)));
    },
    (Object error, StackTrace stack) {
      logZoneError(error, stack);
    },
  );
}
