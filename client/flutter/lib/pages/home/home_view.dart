import 'package:app/ffi/ffi.dart';
import 'package:app/pages/index/index.dart';
import 'package:app/pages/moment/moment_view.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'dashboard_view.dart';
import 'home_controller.dart';


class HomeView extends GetView<HomeController> {

  static const TextStyle optionStyle =
  TextStyle(fontSize: 30, fontWeight: FontWeight.bold);
  final List<Widget> _widgetOptions = <Widget>[
    IndexPage(title: '🍀'),
    Container(
        alignment: Alignment.center,
        child: Text(
          greeting(),
          style: optionStyle,
        )),
    MomentView(),
    DashboardView(),
  ];



  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body:PageView(
          controller: controller.pageController,
          onPageChanged: controller.onPageChanged,
          children: _widgetOptions,
          physics: ClampingScrollPhysics(),
        ),
        bottomNavigationBar: Obx(() =>
            BottomNavigationBar(
              currentIndex: controller.selectedIndex.value,
              onTap: controller.onItemTapped,
              selectedItemColor: Theme
                  .of(context)
                  .canvasColor,
              type: BottomNavigationBarType.fixed,
              backgroundColor: Theme
                  .of(context)
                  .primaryColor
                  .withAlpha(127),
              items: controller.bottomNavigationBarList.map((item)=> item.bottomNavigationBarItem()).toList(),
            ))
    );
  }
}