import 'package:app/ffi/ffi.dart';
import 'package:app/pages/home/splash.dart';
import 'package:app/pages/index/index.dart';
import 'package:app/pages/moment/moment_view.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'dashboard_view.dart';
import 'home_controller.dart';
import 'initialize.dart';



/*
class HomeView extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => HomeState();
}
*/

class HomeView extends StatelessWidget{

  @override
  Widget build(BuildContext context) {
    Get.log("HomeView重绘");
    return FutureBuilder(
      // Replace the 3 second delay with your initialization code:
      future: initialize(),
      builder: (context, AsyncSnapshot snapshot) {
        // Show splash screen while waiting for app resources to load:
        if (snapshot.connectionState == ConnectionState.waiting) {
          return splash;
        } else {
          final app = App();
          WidgetsBinding.instance!.addObserver(app);
          // Loading is done, return the app:
          return app;
        }
      },
    );
  }

}

class App extends StatelessWidget with WidgetsBindingObserver{
  final HomeController controller = Get.find();
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
  void didChangeAppLifecycleState(AppLifecycleState state) {
    print("--" + state.toString());
    switch (state) {
      case AppLifecycleState.inactive: // 处于这种状态的应用程序应该假设它们可能在任何时候暂停。
        break;
      case AppLifecycleState.resumed: // 应用程序可见，前台
        controller.advertising();
        break;
      case AppLifecycleState.paused: // 应用程序不可见，后台
        break;
      case AppLifecycleState.detached: // 申请将暂时暂停
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: PageView(
          controller: controller.pageController,
          //onPageChanged: controller.onPageChanged,
          children: _widgetOptions,
          physics: NeverScrollableScrollPhysics(),
        ),
        bottomNavigationBar: Obx(() => BottomNavigationBar(
              currentIndex: controller.selectedIndex.value,
              onTap: controller.onItemTapped,
              selectedItemColor: Theme.of(context).canvasColor,
              type: BottomNavigationBarType.fixed,
              backgroundColor: Theme.of(context).primaryColor.withAlpha(127),
              items: controller.bottomNavigationBarList
                  .map((item) => item.bottomNavigationBarItem())
                  .toList(),
            )));
  }
}
