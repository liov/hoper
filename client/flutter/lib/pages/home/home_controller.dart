import 'dart:async';

import 'package:flutter/animation.dart';
import 'package:get/get.dart';

import 'home_state.dart';

class HomeController extends GetxController {
  final state = HomeState().obs;

  onItemTapped(int index) {
    state.update((state)=>state!.selectedIndex = index);
  }

  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
    Timer.periodic(
      Duration(seconds: 1),
          (timer) { state.update((state)=>state!.now = DateTime.now());},
    );
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }

}
