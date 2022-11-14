
import 'package:app/pages/home/splash_conroller.dart';
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
      floatingActionButton: GetBuilder<SplashController>(
        builder:(controller){
          if (controller.countdown == 0) {
            return const Text('初始化');
          }
          return GestureDetector(
            onTap: (){
              if (! controller.adCompleter.isCompleted)  controller.adCompleter.complete();
            },
            child: Text('${controller.countdown}秒'),
          );
        }
      ),
    );
  }
}