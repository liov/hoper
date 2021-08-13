
import 'package:app/global/service.dart';

import 'package:flutter/cupertino.dart';

import '../pages/user/login_view.dart';
import 'app_info.dart';
import 'user.dart';
import 'package:get/get.dart';

import 'auth.dart';

export 'service.dart';

final globalState = GlobalState.instance;

class GlobalState {

  GlobalState._();

  static GlobalState? _instance;

  static GlobalState get instance => _instance ??= GlobalState._();

  var appState = AppInfo();
  var authState = AuthState();
  var userState = UserState();

  var _initialized = false;

  init() async {
    if (_initialized) return;
    _initialized = true;
    await globalService.init();
    authState.getAuth();
  }

  Widget? authCheck() => authState.userAuth == null ? LoginView():null;
}
