import 'package:get/get.dart';

import 'moment_controller.dart';

class MomentBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => MomentController());
  }
}
