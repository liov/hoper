
import 'dart:io';


import 'package:app/utils/sqlite.dart';
import 'package:get/get.dart';
import 'package:hive/hive.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';

import 'dao.dart';
import 'global_state/app_info.dart';
import 'global_state/global_controller.dart';


final init = _initialize();

var _initialized = false;

_initialize() async{
  if (_initialized) return;
  final GlobalController globalController =  Get.find();
  Directory appDocDir = await getApplicationDocumentsDirectory();
  final box = await Hive.openBox('box',path:appDocDir.path);
  dao.box = box;
  globalController.appInfo.version = box.get(AppVersionKey,defaultValue:"none");
  Get.put(box,tag:'box');
  var db = await openDatabase('hoper.db',
      onCreate: (Database db, int version) async {
      final hoperExists = await tableExists(db,"hoper");
      if(!hoperExists){
        // When creating the db, create the table
        await db.execute(
            'CREATE TABLE hoper (id INTEGER PRIMARY KEY, name TEXT, value INTEGER, num REAL)');
      }

      });
  dao.db = db;
  Get.put(db,tag:'db');
  _initialized = true;
}