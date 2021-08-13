import 'package:app/pages/home/home_controller.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';



class MomentController extends GetxController with SingleGetTickerProviderMixin{

  final HomeController homeController = Get.find();

  var title = "moment";
  var tabValues = [
    '关注',
    '推荐',
    '刚刚',
  ];
  late final TabController tabController = TabController(
    length: tabValues.length,
    initialIndex: 2,
    vsync: this,
  );


  @override
  void onInit() {
    super.onInit();
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
