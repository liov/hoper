import 'package:app/pages/home/splash_conroller.dart';
import 'package:get/get.dart';

import 'package:app/pages/home/home_controller.dart';

class HomeBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => HomeController());
  }
}
