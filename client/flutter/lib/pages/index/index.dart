import 'package:animate_do/animate_do.dart';
import 'package:app/global/global_controller.dart';

import 'package:app/pages/home/splash_view.dart';
import 'package:app/routes/route.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import 'package:shimmer/shimmer.dart';

import '../login_view.dart';

class IndexPage extends StatefulWidget {
  IndexPage({Key? key, required this.title}) : super(key: key);

  final String title;

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
        title: Text(widget.title),
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
        ],
      ),
      floatingActionButton: FloatingActionButton(
        heroTag: 'login',
        onPressed: (){
          //_methodChannel.invokeMethod("toNative",{"route":"/"}).then((value) => null);
          final user = Get.find<GlobalController>().authState.user;
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
                            globalController.authState.user = null;
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
