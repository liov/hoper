

import 'package:app/generated/protobuf/user/user.model.pb.dart';

class AuthState {
  UserBaseInfo? user = null;
  String? key = null;
  var isActive = false;
}

const AuthKey = "Authorization";