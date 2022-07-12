
import 'dart:async';

import 'package:app/global/global_service.dart';

import 'state/app.dart';
import 'state/user.dart';
import 'package:get/get.dart';

import 'state/auth.dart';

export 'global_service.dart';

final globalState = GlobalState.instance;

class GlobalState extends GetxController{

  GlobalState._();

  static GlobalState? _instance;

  static GlobalState get instance => _instance ??= GlobalState._();

  var appState = AppState();
  var authState = AuthState();
  var userState = UserState();

  var initialized = false;

  var rebuildTimes = 0;

  Future<void> init() async {
    if (initialized) return;
    initialized = true;
    await globalService.init();
    await authState.getAuth();
  }

  var isDarkMode = (AppState.isDebug?true:false).obs;
}
