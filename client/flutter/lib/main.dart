
import 'dart:io';
import 'dart:isolate';

import 'package:app/routes/route.dart';
import 'package:app/global/theme.dart';
import 'package:app/translations/zh_CN/local.dart';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'package:app/global/global_state.dart';
import 'package:flutter_localizations/flutter_localizations.dart';


import 'ffi/ffi.dart';
import 'global/app_info.dart';


Future<void> main() async {
  ErrorWidget.builder = (FlutterErrorDetails flutterErrorDetails){
    print(flutterErrorDetails.toString());
    return Center(
      child: Text("找不到页面"),
    );
  };
  assert(AppInfo.isDebug = true);
  print(AppInfo.isDebug);
  runApp(GetMaterialApp(
    title: 'hoper',
    themeMode: globalState.isDarkMode.value?ThemeMode.dark:ThemeMode.system,
    theme: AppTheme.light,
    darkTheme:AppTheme.dark,
      builder: (context, child) => Scaffold(
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
      ),
    //home: HomeView(),
    initialRoute: Routes.HOME,
    initialBinding: BindingsBuilder.put(() =>globalState),
    getPages: AppPages.routes,
    localeListResolutionCallback:
        (List<Locale>? locales, Iterable<Locale> supportedLocales) {
      return const Locale('zh');
    },
    localeResolutionCallback:
        (Locale? locale, Iterable<Locale> supportedLocales) {
      return const Locale('zh');
    },
    localizationsDelegates: [
      GlobalMaterialLocalizations.delegate,
      GlobalWidgetsLocalizations.delegate,
      GlobalCupertinoLocalizations.delegate,
    ],
    supportedLocales: [
      const Locale('zh', 'CN'),
      const Locale('en', 'US'),
    ],
  ));
}

