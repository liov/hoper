import 'package:flutter/animation.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'moment_state.dart';

class MomentController extends GetxController with SingleGetTickerProviderMixin{
  final state = MomentState();
  late TabController ac;

  @override
  void onInit() {
    super.onInit();
    ac = TabController(
      length: state.tabValues.length,
      initialIndex: 1,
      vsync: this,
    );
  }

  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }
}
