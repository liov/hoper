import 'package:get/get.dart';

import 'global.dart';

class GlobalBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => GlobalState());
  }
}
