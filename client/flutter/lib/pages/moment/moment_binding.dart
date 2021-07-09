import 'package:app/service/moment.dart';
import 'package:get/get.dart';

import 'list/moment_list_controller.dart';
import 'moment_controller.dart';

class MomentBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => MomentController());

    Get.lazyPut(() => MomentListController());
    Get.lazyPut(() => MomentClient());
  }
}
