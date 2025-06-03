
import 'package:app/pages/home/splash_conroller.dart';
import 'package:app/pages/home/splash_view.dart';
import 'package:app/pages/weibo/image_view.dart';
import 'package:flutter/cupertino.dart';
import 'package:get/get.dart';

import '../../global/state.dart';
import 'home_view.dart';

class StartView extends StatelessWidget {
  final SplashController controller = Get.find();

  StartView({super.key});

  @override
  Widget build(BuildContext context) {

    controller.startAd();
    globalState.rebuildTimes++;
    globalService.logger.d("StartView重绘${globalState.rebuildTimes}次");

    return FutureBuilder(
      // Replace the 3 second delay with your initialization code:
      future: controller.adCompleter.future,
      builder: (context, AsyncSnapshot snapshot) {
        // Show splash screen while waiting for app resources to load:
        if (snapshot.connectionState == ConnectionState.waiting) {
          return splash;
        } else {
/*          final app = App();
          WidgetsBinding.instance.addObserver(app);*/
          // Loading is done, return the app:
          return ImageView();
        }
      },
    );
  }
}