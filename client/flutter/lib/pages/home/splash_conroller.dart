
import 'dart:async';
import 'package:flutter/material.dart';

import 'package:app/global/controller.dart';
import 'package:get/get.dart';


class SplashController extends GetxController {
  late Completer<Null> adCompleter;
  DateTime? pausedTime;
  int countdown = 0;
  set duration(Duration duration){
    countdown = (duration.inMilliseconds/1000).round();
    print(countdown);
    update();
    Timer.periodic(Duration(seconds: 1), (timer) {
      countdown--;
      if(countdown == 0){
        timer.cancel();
        if (!adCompleter.isCompleted) adCompleter.complete();
      }
      update();
    });
  }
  void advertising(Widget splash){
    if(pausedTime ==null) return;
    final current = DateTime.now();
    if (current.difference(pausedTime!)  < Duration(seconds:1)) return;
    print('advertising');
    Get.showOverlay(loadingWidget:splash, asyncFunction: () {
      adCompleter = Completer();
      duration = Duration(seconds:3);
      return adCompleter.future;
    });
  }

  void startAd(){
    if(globalState.initialized) return;
    adCompleter = Completer();
    final startTime = DateTime.now();
    globalState.init().then((v){
      final current = DateTime.now();
      final duration = current.difference(startTime);
      if (duration < Duration(seconds:3)){
        this.duration = Duration(seconds:3) - duration;
      }else{
        if (!adCompleter.isCompleted) adCompleter.complete();
      };
    });
  }


}