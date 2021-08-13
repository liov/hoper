import 'package:animate_do/animate_do.dart';
import 'package:app/global/controller.dart';

import 'package:app/pages/home/splash_view.dart';
import 'package:app/routes/route.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import 'package:shimmer/shimmer.dart';

import '../user/login_view.dart';

class IndexPage extends StatefulWidget {


  @override
  _IndexPageState createState() => _IndexPageState();
}

class _IndexPageState extends State<IndexPage> with AutomaticKeepAliveClientMixin {
  final MethodChannel _methodChannel = MethodChannel('xyz.hoper.native/view');

  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      _counter++;
    });
  }

  void _reduceCounter() {
    setState(() {
      _counter--;
    });
  }

  @override
  Widget build(BuildContext context) {
    Get.log("IndexPage重绘");
    super.build(context);
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: Text('🍀'),
      ),
      body: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
        Shimmer.fromColors(
          loop:1,
          baseColor: Colors.white,
          highlightColor: Colors.blue,
          child:Text('点击下方加号:',)
      ),
          Text(
            '$_counter',
            style: Get.theme.textTheme.headline4,
          ),
          Row(
              mainAxisAlignment: MainAxisAlignment.center,
              verticalDirection: VerticalDirection.up,
              children: [
                FadeInLeft(child:FloatingActionButton(
                  heroTag: 'Increment',
                  onPressed: _incrementCounter,
                  tooltip: 'Increment',
                  child: Icon(Icons.add),
                )),
                SizedBox(width: 20),
                FadeInRight(child:FloatingActionButton(
                  heroTag: 'Reduce',
                  onPressed: _reduceCounter,
                  tooltip: 'Reduce',
                  child: Icon(Icons.remove),
                ))
              ]),
            CupertinoSwitch(
            value: Get.isDarkMode,
            onChanged: (value){
              Get.changeTheme(Get.isDarkMode? ThemeData.light(): ThemeData.dark());
            }),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        heroTag: 'login',
        onPressed: (){
          //_methodChannel.invokeMethod("toNative",{"route":"/"}).then((value) => null);
          final user = globalState.authState.userAuth;
          if ( user!= null) {
              Get.dialog(
                 AlertDialog(
                      title: Text('提示',textAlign:TextAlign.center),
                      content: Text('确认退出吗？',textAlign:TextAlign.center),
                      actions: <Widget>[
                        TextButton(
                          child: Text('取消'),
                          onPressed: () {
                            navigator!.pop('cancel');
                          },
                        ),
                        TextButton(
                          child: Text('确认'),
                          onPressed: () {
                            globalState.authState.userAuth = null;
                            navigator!.pop('ok');
                          },
                        ),
                      ],
                  ));
            }
          else {Get.toNamed(Routes.LOGIN);}
        },

        tooltip: 'ToBrowser',
        child: Icon(Icons.send),
      ),
      // This trailing comma makes auto-formatting nicer for build methods.
      floatingActionButtonAnimator: FloatingActionButtonAnimator.scaling,
    );
  }

  @override
  bool get wantKeepAlive => true;
}
