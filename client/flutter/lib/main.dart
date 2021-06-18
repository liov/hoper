import 'package:app/page/ffi/ffi.dart';
import 'package:app/page/index/fourth.dart';
import 'package:app/page/index/index.dart';
import 'package:app/page/loginView.dart';
import 'package:app/page/moment/momentListView.dart';

import 'package:flutter/material.dart';
import 'package:get/get.dart';


import 'model/state/global.dart';
import 'model/state/user.dart';

void main() async {
  runApp(
      GetMaterialApp(
          title: 'hoper',
          theme: ThemeData(
            // This is the theme of your application.
            //
            // Try running your application with "flutter run". You'll see the
            // application has a blue toolbar. Then, without quitting the app, try
            // changing the primarySwatch below to Colors.green and then invoke
            // "hot reload" (press "r" in the console where you ran "flutter run",
            // or simply save your changes to "hot reload" in a Flutter IDE).
            // Notice that the counter didn't reset back to zero; the application
            // is not restarted.
            primarySwatch: Colors.blue,
            primaryColor: Colors.blueAccent[700],
            // This makes the visual density adapt to the platform that you run
            // the app on. For desktop platforms, the controls will be smaller and
            // closer together (more dense) than on mobile platforms.
            visualDensity: VisualDensity.adaptivePlatformDensity,
          ),
          home: MyHomePage(),
        routes: {
          '/content': (context) => MomentListView(),
          '/login': (context) => LoginView(),
        },)
  );
}

class MyHomePage extends StatelessWidget {

  // 使用Get.put()实例化你的类，使其对当下的所有子路由可用。
  final AuthState userInfo = Get.put(AuthState());
  final GlobalState globalState = Get.put(GlobalState());

  static const TextStyle optionStyle =
  TextStyle(fontSize: 30, fontWeight: FontWeight.bold);
  final List<Widget> _widgetOptions = <Widget>[
    IndexPage(title: '🍀'),
    Container(
        alignment: Alignment.center,
        child: Text(
          greeting(),
          style: optionStyle,
        )
    ),
    FourthPage(),
  ];

  @override
  Widget build(BuildContext context) {


    return Scaffold(
      body: Obx(()=>_widgetOptions.elementAt(globalState.selectedIndex.value)),
      bottomNavigationBar: Obx(()=>BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        backgroundColor: Theme
            .of(context)
            .primaryColor
            .withAlpha(127),
        items: <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: Icon(Icons.movie),
            label: 'flutter',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.search),
            label: 'rustffi',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.twenty_one_mp_rounded),
            label: 'hoper',
          ),
        ],
        currentIndex: globalState.selectedIndex.value,
        selectedItemColor: Theme
            .of(context)
            .canvasColor,
        onTap: globalState.onItemTapped,
      ),
    ));
  }
}