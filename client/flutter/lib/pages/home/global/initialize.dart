
import 'dart:io';

import 'package:path/path.dart';
import 'package:app/service/user.dart';
import 'package:app/utils/sqlite.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:hive/hive.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'package:app/generated/protobuf/utils/request/param.pb.dart' as req;
import 'dao.dart';
import 'global_state/app_info.dart';
import 'global_state/auth.dart';
import 'global_state/global_controller.dart';
import 'package:fixnum/fixnum.dart' as $fixnum;

final init = _initialize();

var _initialized = false;

_initialize() async{
  if (_initialized) return;
  final GlobalController globalController =  Get.find();
  final UserClient userClient = Get.put(UserClient());
  final Dao dao = Get.put(Dao());
  Directory appDocDir = await getApplicationDocumentsDirectory();
  appDocDir.list().forEach((element) {print(element);});
  final box = await Hive.openBox('box',path:join(appDocDir.path,"hive"));
  dao.box = box;
  globalController.appInfo.version = box.get(AppVersionKey,defaultValue:"none");
  globalController.authState.key = box.get(AuthKey);
  if (globalController.authState.key!=null){
    userClient.stub.info(req.Object(id:$fixnum.Int64(0)),options:CallOptions(metadata:{AuthKey:globalController.authState.key!}));
  }
  Get.put(box,tag:'box');
  final db = await openDatabase(join(appDocDir.path,'databases','hoper.db'), version: 1,
      onCreate: (Database db, int version) async {
        // When creating the db, create the table
        await db.execute(
            'CREATE TABLE hoper (id INTEGER PRIMARY KEY, name TEXT, value INTEGER, num REAL)');
      });
  print(db);
  dao.db = db;
  Get.put(db,tag:'db');
  _initialized = true;
}