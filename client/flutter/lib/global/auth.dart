

import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/model/const/const.dart';


class AuthState {
  UserAuthInfo? user = null;

  static const _PRE = "AuthState";
  static const StringAuthKey = _PRE+Authorization;
  static const StringAccountKey = _PRE+"AccountKey";

}


