import 'package:app/pages/start/splash_conroller.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

const splash = Splash();

class Splash extends StatelessWidget {
  const Splash({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: ExtendedImage.asset(
        "assets/splash/splash.jpg",
        alignment: Alignment.center,
        width: Get.width,
        height: Get.height,
        fit: BoxFit.fill,
        //cancelToken: cancellationToken,
      ),
      floatingActionButton: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
        decoration: BoxDecoration(
          color: Colors.black.withValues(alpha: 0.5),
          borderRadius: BorderRadius.circular(16),
        ),
        child: GetBuilder<SplashController>(
          builder: (controller) {
            if (controller.countdown == 999) {
              return const Text('初始化');
            }
            return GestureDetector(
              onTap: () {
                if (!controller.adCompleter.isCompleted) {
                  controller.adCompleter.complete();
                }
              },
              child: Text('跳过 ${controller.countdown}秒'),
            );
          },
        ),
      ),
    );
  }
}
