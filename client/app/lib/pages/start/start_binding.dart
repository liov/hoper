import 'package:app/pages/start/splash_conroller.dart';
import 'package:get/get.dart';


class StartBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut(() => SplashController());
  }
}
