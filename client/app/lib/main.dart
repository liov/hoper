import 'dart:async';

import 'package:app/pages/route.dart';
import 'package:app/global/theme.dart';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'package:app/global/state.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import 'package:app/global/state/app.dart';


Future<void> main() async {
  runZonedGuarded(() async {
    WidgetsFlutterBinding.ensureInitialized();
    ErrorWidget.builder = (FlutterErrorDetails flutterErrorDetails) {
      globalService.logger.fine(flutterErrorDetails.toString());
      return const Center(
        child: Text("找不到页面"),
      );
    };
    globalService.logger.fine('runZonedGuarded');
    assert(AppState.isDebug = true);
    globalService.logger.fine("${AppState.isDebug}");
    runApp( GetMaterialApp(
      title: 'hoper',
      themeMode: globalState.isDarkMode.value ? ThemeMode.dark : ThemeMode
          .system,
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      builder: (context, child) {
        globalService.logger.fine('GetMaterialApp');
        return  Listener(
          onPointerDown: (_) {
              FocusScopeNode focusScopeNode = FocusScope.of(context);
              if (!focusScopeNode.hasPrimaryFocus &&
                  focusScopeNode.focusedChild != null) {
                FocusManager.instance.primaryFocus?.unfocus();
              }
            },
            child: child,
        );
      },
      //home: HomeView(),
      initialRoute: Routes.START,
      initialBinding: BindingsBuilder.put(() => globalState),
      getPages: Routes.pages,
      localeListResolutionCallback:
          (List<Locale>? locales, Iterable<Locale> supportedLocales) {
        return const Locale('zh');
      },
      localeResolutionCallback:
          (Locale? locale, Iterable<Locale> supportedLocales) {
        return const Locale('zh');
      },
      localizationsDelegates: const [
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: const [
        Locale('zh', 'CN'),
        Locale('en', 'US'),
      ],
    ));
  }, (dynamic error, StackTrace stack) {
    globalService.logger.severe("Something went wrong!", error, stack);
  });
}
