import 'package:get/get.dart';

import 'global_controller.dart';



class GlobalBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => GlobalController());
  }
}
