import 'dart:io';

import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/utils/empty/empty.pb.dart';
import 'package:app/model/const/const.dart';
import 'package:app/service/user.dart';
import 'package:flutter/cupertino.dart';
import 'package:grpc/grpc.dart';
import 'package:hive/hive.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart' as path;
import '../../login_view.dart';
import '../../../../service/dao.dart';
import 'app_info.dart';
import 'user.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';
import 'auth.dart';

final globalController = GlobalController();

class GlobalController extends GetxController {
  var appState = AppInfo();
  var authState = AuthState();
  var userState = UserState();
  final UserClient userClient = Get.put(UserClient());
  late final Box box;
  late final Database db;
  var _initialized = false;

  init() async {
    if (_initialized) return;
    _initialized = true;
    final Dao dao = Get.put(Dao());
    Directory appDocDir = await getApplicationDocumentsDirectory();
    appDocDir.list().forEach((element) {
      print(element);
    });

    final boxfuture = () async {
      box = await Hive.openBox('box', path: path.join(appDocDir.path, "hive"));
      dao.box = box;
      appState.init(box);
      getAuth();
      Get.put(box, tag: 'box');
    };

    final dbfuture = () async {
      db = await openDatabase(
          path.join(appDocDir.path, 'databases', 'hoper.db'),
          version: 1, onCreate: (Database db, int version) async {
        // When creating the db, create the table
        await db.execute(
            'CREATE TABLE hoper (id INTEGER PRIMARY KEY, name TEXT, value INTEGER, num REAL)');
      });

      dao.db = db;
      Get.put(db, tag: 'db');
    };
    await Future.wait([boxfuture(), dbfuture()]);
  }

  getAuth() async {
    if (authState.user != null) return;
    final authKey = box.get(AuthState.StringAuthKey);
    if (authKey != null) {
      try {
        final user = await userClient.stub.authInfo(Empty(),
            options: CallOptions(metadata: {Authorization: authKey}));
        authState.user = user;
        return null;
      } catch (err) {
        print(err);
      }
    }
  }

  Widget? authCheck() {
    if (authState.user == null) return LoginView();
    return null;
  }
}
