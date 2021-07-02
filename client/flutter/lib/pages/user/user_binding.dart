import 'package:get/get.dart';

import 'user_logic.dart';

class userBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => userLogic());
  }
}
