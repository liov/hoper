import 'dart:async';

import 'package:app/components/bottom/bottom.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

class HomeController extends GetxController {
  var selectedIndex = 0.obs;
  var now = DateTime.now().obs;

  final pageController = PageController();
  var bottomNavigationBarList = [
    Bottom(Icons.home, "home"),
    Bottom(Icons.account_box_rounded, "Profile"),
    Bottom(Icons.account_box_rounded, "Moments"),
    Bottom(Icons.account_balance_sharp, "Weather"),
  ];

  HomeController get to => Get.find<HomeController>();

  void onPageChanged(int index) {
    selectedIndex.value = index;
  }

  void onItemTapped(int index) {
    pageController.jumpToPage(index);
    selectedIndex.value = index;
  }

  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
    Timer.periodic(
      Duration(seconds: 1),
          (timer) { now.value = DateTime.now();},
    );
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }

}
