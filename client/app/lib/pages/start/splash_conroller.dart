
import 'dart:async';
import 'package:flutter/material.dart';

import 'package:app/global/state.dart';
import 'package:get/get.dart';


class SplashController extends GetxController {
  late Completer<void> adCompleter;
  DateTime? pausedTime;
  int countdown = 999;
  Timer? timer;
  Duration time = Duration(seconds:3);
  set duration(Duration duration){
    countdown = (duration.inMilliseconds/1000).round();
    globalService.logger.d(countdown);
    update();
    timer = Timer.periodic(const Duration(seconds: 1), (timer) {
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
    final current = DateTime.now();
    if (current.difference(pausedTime!)  < const Duration(minutes:10)) return;
    globalService.logger.d('advertising');
    Get.showOverlay(loadingWidget:splash, asyncFunction: () {
      adCompleter = Completer();
      duration = time;
      return adCompleter.future;
    });
  }

  void init(){
    adCompleter = Completer();
    if(globalState.initialized) {
      duration = time;
      return;
    }

    globalState.init().then((v){
      duration = time;
    });
  }

  void skip(){
    timer?.cancel();
    if (!adCompleter.isCompleted) adCompleter.complete();
  }

}