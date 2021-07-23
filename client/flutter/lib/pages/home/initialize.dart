import 'dart:async';
import 'dart:io';

import 'package:app/model/global_state/global_controller.dart';
import 'package:get/get.dart';
import 'package:hive/hive.dart';
import 'package:path_provider/path_provider.dart';

import 'home_controller.dart';

initialize() async{
  Directory appDocDir = await getApplicationDocumentsDirectory();
  final box = await Hive.openBox('box',path:appDocDir.path);
  Get.put(box,tag:'box');
}