
import 'package:app/routes/route.dart';
import 'package:app/theme.dart';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'package:app/global/global_controller.dart';



Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
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

