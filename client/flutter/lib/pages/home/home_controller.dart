import 'dart:async';

import 'package:app/components/bottom/bottom.dart';
import 'package:app/pages/home/global/splash_view.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:get/get.dart';

class HomeController extends GetxController {
  var selectedIndex = 0.obs;
  var now = DateTime.now().obs;
  var stackIndex = 1.obs;

  final pageController = PageController();
  var bottomNavigationBarList = [
    Bottom.icon(Icons.home, "home"),
    Bottom(FaIcon(FontAwesomeIcons.gamepad), "Profile"),
    Bottom.icon(Icons.account_box_rounded, "Moments"),
    Bottom.icon(Icons.account_balance_sharp, "Weather"),
  ];

  HomeController get to => Get.find<HomeController>();

  void advertising(){
    if(pausedTime ==null) return;
    final current = DateTime.now();
    if (current.difference(pausedTime!)  < Duration(minutes:1)) return;
    Get.dialog(splash);
    //Get.to(()=>splash);
    Future.delayed(Duration(seconds:1),(){navigator!.pop('ok');});
  }

  DateTime? pausedTime;

  void onPageChanged(int index) {
    selectedIndex.value = index;
  }

  void onItemTapped(int index) {
    pageController.jumpToPage(index);
    selectedIndex.value = index;
  }

  // PageView + TabView 连续滑动
  var  scrollNum = 0;
  void continueScroll() {
    final index = selectedIndex.value + (scrollNum<0 ? -1 : 1);
    scrollNum = 0;
    pageController.animateToPage(
        index, duration: kTabScrollDuration,
        curve: Curves.ease);

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
