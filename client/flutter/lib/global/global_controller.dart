import 'dart:io';

import 'package:app/generated/protobuf/utils/empty/empty.pb.dart';
import 'package:app/model/const/const.dart';
import 'package:app/service/upload.dart';
import 'package:app/service/user.dart';
import 'package:app/utils/dio.dart' as $dio;
import 'package:dio/dio.dart';
import 'package:flutter/cupertino.dart';
import 'package:grpc/grpc.dart';
import 'package:hive/hive.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart' as path;
import '/pages/login_view.dart';
import '/service/dao.dart';
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
  final UploadClient uploadClient = Get.put(UploadClient());
  final Dio httpClient = $dio.httpClient;
  late CallOptions options;


  late final Box box;
  late final Database db;

  var _initialized = false;

  init() async {
    if (_initialized) return;
    _initialized = true;
    Directory appDocDir = await getApplicationDocumentsDirectory();
    appDocDir.list().forEach((element) {
      print(element);
    });

    final boxfuture = () async {
      box = await Hive.openBox('box', path: path.join(appDocDir.path, "hive"));
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

      Get.put(db, tag: 'db');
    };
    await Future.wait([boxfuture(), dbfuture()]);

  }

  getAuth() async {
    if (authState.user != null) return;
    final authKey = box.get(AuthState.StringAuthKey);
    setAuth(authKey);
    if (authKey != null) {
      try {
        final user = await userClient.stub.authInfo(Empty(),options:options);
        authState.user = user;
        return null;
      } catch (err) {
        print(err);
      }
    }
  }

  setAuth(String authKey){
    options = CallOptions(metadata: {Authorization: authKey});
    httpClient.options.headers[Authorization] = authKey;
  }

  Widget? authCheck() {
    if (authState.user == null) return LoginView();
    return null;
  }
}
