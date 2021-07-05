

import 'package:get/get.dart';

import '../user.dart';


class AuthState {
  User? user = null;
  late String cookie;
  var isActive = false;
}
