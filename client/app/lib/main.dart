import 'dart:async';
import 'dart:isolate';

import 'package:app/routes/route.dart';
import 'package:app/global/theme.dart';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'package:app/global/state.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

import 'package:applib/widget/hot_load_warp.dart';
import 'package:app/ffi/ffi.dart';
import 'package:app/global/state/app.dart';


Future<void> main() async {
  runZonedGuarded(() {
    WidgetsFlutterBinding.ensureInitialized();
    ErrorWidget.builder = (FlutterErrorDetails flutterErrorDetails) {
      globalService.logger.d(flutterErrorDetails.toString());
      return const Center(
        child: Text("找不到页面"),
      );
    };
    globalService.logger.d('runZonedGuarded');
    assert(AppState.isDebug = true);
    globalService.logger.d("${AppState.isDebug}");
    runApp( GetMaterialApp(
      title: 'hoper',
      themeMode: globalState.isDarkMode.value ? ThemeMode.dark : ThemeMode
          .system,
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      builder: (context, child) {
      globalService.logger.d('GetMaterialApp');
        return  Scaffold(
            body: GestureDetector(
              onTap: () {
                FocusScopeNode focusScopeNode = FocusScope.of(context);
                if (!focusScopeNode.hasPrimaryFocus &&
                    focusScopeNode.focusedChild != null) {
                  FocusManager.instance.primaryFocus?.unfocus();
                }
              },
              child: child,
            ),
          );
        },
      //home: HomeView(),
      initialRoute: Routes.HOME,
      initialBinding: BindingsBuilder.put(() => globalState),
      getPages: AppPages.routes,
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
    globalService.logger.e("Something went wrong!", error: error,  stackTrace: stack);
  });
}
