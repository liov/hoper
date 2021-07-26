import 'package:app/generated/protobuf/user/user.model.pb.dart';

import 'app_info.dart';
import 'user.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';
import 'auth.dart';

class GlobalController extends GetxController {
  var appInfo = AppInfo();
  var authState = AuthState();
  var userState = UserState();

}
