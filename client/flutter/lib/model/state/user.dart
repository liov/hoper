
import 'package:get/get.dart';

import '../user.dart';

class AuthState extends GetxController {
  Rxn<User?> user = Rxn<User?>();
  late RxString cookie;
  var isActive =false.obs;
}

class UserState extends GetxController {
  var users = Map<int, User>().obs;
}