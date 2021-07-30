

import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:flutter/cupertino.dart';
import 'package:hive/hive.dart';


class AuthState {
  UserAuthInfo? user = null;

  static const _PRE = "AuthState";
  static const StringAuthKey = _PRE+"Authorization";
  static const StringAccountKey = _PRE+"AccountKey";

}


