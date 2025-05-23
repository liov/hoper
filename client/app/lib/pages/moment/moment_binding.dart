import 'package:app/global/service.dart';
import 'package:app/rpc/action.dart';
import 'package:app/rpc/content.dart';
import 'package:app/rpc/moment.dart';
import 'package:get/get.dart';

import 'list/moment_list_controller.dart';
import 'moment_controller.dart';

class MomentBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => MomentController());

    Get.lazyPut(() => MomentListController());
    Get.lazyPut(() => MomentClient(globalService.subject));
    Get.lazyPut(() => ContentClient(globalService.subject));
    Get.lazyPut(() => ActionClient(globalService.subject));
  }
}
