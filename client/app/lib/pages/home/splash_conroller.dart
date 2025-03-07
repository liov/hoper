
import 'dart:async';
import 'package:flutter/material.dart';

import 'package:app/global/state.dart';
import 'package:get/get.dart';


class SplashController extends GetxController {
  late Completer<void> adCompleter;
  DateTime? pausedTime;
  int countdown = 0;
  set duration(Duration duration){
    countdown = (duration.inMilliseconds/1000).round();
    globalService.logger.d(countdown);
    update();
    Timer.periodic(const Duration(seconds: 1), (timer) {
      countdown--;
      if(countdown <= 0){
        timer.cancel();
        if (!adCompleter.isCompleted) adCompleter.complete();
      }
      update();
    });
  }
  void advertising(Widget splash){
    if(pausedTime ==null) return;
    const time = Duration(seconds:3);
    final current = DateTime.now();
    if (current.difference(pausedTime!)  < const Duration(minutes:10)) return;
    globalService.logger.d('advertising');
    Get.showOverlay(loadingWidget:splash, asyncFunction: () {
      adCompleter = Completer();
      duration = time;
      return adCompleter.future;
    });
  }

  void startAd(){
    if(globalState.initialized) return;
    const time = Duration(seconds:1);
    adCompleter = Completer();
    final startTime = DateTime.now();
    globalState.init().then((v){
      final current = DateTime.now();
      final duration = current.difference(startTime);
      if (duration < time){
        this.duration = time - duration;
      }else{
        if (!adCompleter.isCompleted) adCompleter.complete();
      };
    });
  }


}