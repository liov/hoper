
import 'package:app/pages/home/global/initialize.dart';
import 'package:app/routes/route.dart';
import 'package:app/theme.dart';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'package:app/pages/home/global/global_state/global_controller.dart';



Future<void> main() async {
  runApp(GetMaterialApp(
    title: 'hoper',
    theme: appThemeData,
    darkTheme:darkThemeData,
    //home: HomeView(),
    initialRoute: Routes.HOME,
    initialBinding: BindingsBuilder.put(() =>GlobalController()),
    getPages: AppPages.routes,
  ));
}

