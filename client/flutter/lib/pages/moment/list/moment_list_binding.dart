import 'package:get/get.dart';

import 'moment_list_controller.dart';

class MomentListBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => MomentListController());
  }
}
