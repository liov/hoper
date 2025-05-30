import 'package:get/get.dart';

import '../../rpc/weibo.dart';
import 'controller.dart';
class WeiboBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => WeiboController());
    Get.lazyPut(() => WeiboClient());

  }
}
