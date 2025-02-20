import 'dart:async';

import 'package:app/components/bottom/bottom.dart';
import 'package:app/global/state.dart';
import 'package:app/pages/home/splash_view.dart';
import 'package:app/utils/dialog.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:get/get.dart';

class HomeController extends GetxController {
  var selectedIndex = 0.obs;
  var now = DateTime.now().obs;
  var pageIndex = 0;

  final pageController = PageController();
  var bottomNavigationBarList = [
    Bottom.icon(Icons.home, label: "瞬间", pageIndex: 0),
    Bottom.icon(Icons.gamepad, label: "Profile", pageIndex: 1),
    Bottom.icon(Icons.add,  onTap: () => toast("测试")),
    Bottom.icon(Icons.account_box_rounded, label: "Moments", pageIndex: 2),
    Bottom.icon(FontAwesomeIcons.user, label: "我的", pageIndex: 3),
  ];
  // 底部索引对应的page索引
  var pageBottomIdx = [0,1,3,4];

  static HomeController get to => Get.find<HomeController>();

  void onPageChanged(int index) {
    selectedIndex.value = pageBottomIdx[index];
  }

  void onItemTapped(int index) {
    final bottom = bottomNavigationBarList[index];
    if (bottom.pageIndex != null) {
      pageController.jumpToPage(bottomNavigationBarList[index].pageIndex!);
      pageIndex = bottom.pageIndex!;
    }
    if (bottom.onTap != null) bottom.onTap!();
    selectedIndex.value = index;
  }

  // PageView + TabView 连续滑动
  var scrollNum = 0;

  void continueScroll() {
    final index = pageIndex + (scrollNum < 0 ? -1 : 1);
    if (index < 0) return;
    scrollNum = 0;
    pageController.animateToPage(index,
        duration: kTabScrollDuration, curve: Curves.ease);
  }

  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
    Timer.periodic(
      const Duration(seconds: 1),
      (timer) {
        now.value = DateTime.now();
      },
    );
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }
}
