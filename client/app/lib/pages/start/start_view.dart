
import 'package:app/pages/start/splash_conroller.dart';
import 'package:app/pages/start/splash_view.dart';
import 'package:app/pages/weibo/image_view.dart';
import 'package:flutter/cupertino.dart';
import 'package:get/get.dart';

import '../../global/state.dart';
import '../home/home_view.dart';
import '../route.dart';

class StartView extends StatelessWidget {
  final SplashController controller = Get.find();

  StartView({super.key});

  @override
  Widget build(BuildContext context) {
    controller.init();
    controller.adCompleter.future.then((value) {
      Get.offNamed(Routes.HOME);
    });
    globalState.rebuildTimes++;
    globalService.logger.d("StartView重绘${globalState.rebuildTimes}次");

    return splash;
  }
}