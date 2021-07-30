
import 'package:app/routes/route.dart';
import 'package:app/theme.dart';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'package:app/pages/home/global/global_controller.dart';



Future<void> main() async {
  runApp(GetMaterialApp(
    title: 'hoper',
    theme: appThemeData,
    darkTheme:darkThemeData,
    //home: HomeView(),
    initialRoute: Routes.HOME,
    initialBinding: BindingsBuilder.put(() =>globalController),
    getPages: AppPages.routes,
  ));
}

