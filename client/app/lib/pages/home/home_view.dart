

import 'dart:io';

import 'package:app/components/christmas_tree.dart';

import 'package:app/ffi/ffi.dart';
import 'package:app/global/state.dart';
import 'package:app/global/theme.dart';
import 'package:app/pages/home/splash_conroller.dart';
import 'package:app/pages/home/splash_view.dart';
import 'package:app/pages/index/index.dart';
import 'package:app/pages/moment/moment_view.dart';
import 'package:app/pages/user/user_view.dart';
import 'package:app/routes/route.dart';

import 'package:convex_bottom_bar/convex_bottom_bar.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:webview_flutter/webview_flutter.dart';

import 'package:app/components/media/pick.dart';
import 'home_controller.dart';

/*
class HomeView extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => HomeState();
}
*/

class HomeView extends StatelessWidget {
  final SplashController controller = Get.find();

   HomeView({super.key});

  @override
  Widget build(BuildContext context) {

    controller.startAd();
    globalState.rebuildTimes++;
    globalService.logger.d("HomeView重绘${globalState.rebuildTimes}次");

    return FutureBuilder(
      // Replace the 3 second delay with your initialization code:
      future: controller.adCompleter.future,
      builder: (context, AsyncSnapshot snapshot) {
        // Show splash screen while waiting for app resources to load:
        if (snapshot.connectionState == ConnectionState.waiting) {
          return splash;
        } else {
          final app = App();
          WidgetsBinding.instance.addObserver(app);
          // Loading is done, return the app:
          return app;
        }
      },
    );
  }
}

class App extends StatelessWidget with WidgetsBindingObserver {
  final HomeController controller = Get.find();
  final SplashController splashController = Get.find();

  static const TextStyle optionStyle =
      TextStyle(fontSize: 30, fontWeight: FontWeight.bold);

  final List<Widget> _widgetOptions = <Widget>[
    const MomentView(),
    const IndexPage(),
    const MediaPick(),
    UserView(),
  ];

  App({super.key});

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    globalService.logger.d("--$state");
    switch (state) {
      case AppLifecycleState.inactive: // 处于这种状态的应用程序应该假设它们可能在任何时候暂停。
        break;
      case AppLifecycleState.resumed: // 应用程序可见，前台
        splashController.advertising(splash);
        break;
      case AppLifecycleState.paused: // 应用程序不可见，后台
        splashController.pausedTime = DateTime.now();
        break;
      case AppLifecycleState.detached: // 申请将暂时暂停
        break;
      case AppLifecycleState.hidden:
        // TODO: Handle this case.
    }
  }

  @override
  Widget build(BuildContext context) {
    globalService.logger.d('重建了');
    return Scaffold(
        body: PageView(
          controller: controller.pageController,
          physics: const NeverScrollableScrollPhysics(),
          //onPageChanged: controller.onPageChanged,
          children: _widgetOptions,
        ),
 /*       floatingActionButtonLocation: CustomFloatingActionButtonLocation(FloatingActionButtonLocation.miniCenterDocked, 0, 15),
        floatingActionButton: FloatingActionButton(
            backgroundColor:Get.theme.backgroundColor,
          onPressed: () {  },),*/
        bottomNavigationBar: Obx(()=>_bottom2()));
  }

  Widget _bottom1(){
    final ThemeData theme = globalState.isDarkMode.value ? AppTheme.dark : AppTheme.light;
    return BottomNavigationBar(
      currentIndex: controller.selectedIndex.value,
      onTap: controller.onItemTapped,
      selectedItemColor: theme.canvasColor,
      type: BottomNavigationBarType.fixed,
      backgroundColor: theme.primaryColor,
        landscapeLayout:BottomNavigationBarLandscapeLayout.linear,
      items: controller.bottomNavigationBarList
          .map((item) => item.navigationBarItem())
          .toList(),
    );
  }

  Widget _bottom2(){
    globalService.logger.d(Get.theme.colorScheme.background);
    final ThemeData theme = globalState.isDarkMode.value ? AppTheme.dark : AppTheme.light;
    return ConvexAppBar(
      disableDefaultTabController: true,
      initialActiveIndex: controller.selectedIndex.value,
      onTabNotify: (i) {
        var intercept = i == 2;
        if (intercept) {
          Get.toNamed(Routes.MOMENT_ADD);
        }
        return !intercept;
      },
      onTap: controller.onItemTapped,
      activeColor: theme.canvasColor,
      backgroundColor: theme.primaryColor,
      style: TabStyle.fixedCircle,
      items: controller.bottomNavigationBarList.map((e) => e.tabItem()).toList(),
    );
  }
}
