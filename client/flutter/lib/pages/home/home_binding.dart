import 'package:app/pages/home/splash_conroller.dart';
import 'package:get/get.dart';

import '../../pages/home/home_controller.dart';

class HomeBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => HomeController());
    Get.lazyPut(() => SplashController());
  }
}
