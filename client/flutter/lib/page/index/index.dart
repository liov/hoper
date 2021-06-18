import 'package:app/model/state/user.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';

import '../loginView.dart';

class IndexPage extends StatefulWidget {
  IndexPage({Key? key, required this.title}) : super(key: key);

  final String title;

  @override
  _IndexPageState createState() => _IndexPageState();
}

class _IndexPageState extends State<IndexPage> {
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
          Text(
            '点击下方加号:',
          ),
          Text(
            '$_counter',
            style: Theme.of(context).textTheme.headline4,
          ),
          Row(
              mainAxisAlignment: MainAxisAlignment.center,
              verticalDirection: VerticalDirection.up,
              children: [
                FloatingActionButton(
                  heroTag: 'Increment',
                  onPressed: _incrementCounter,
                  tooltip: 'Increment',
                  child: Icon(Icons.add),
                ),
                SizedBox(width: 20),
                FloatingActionButton(
                  heroTag: 'Reduce',
                  onPressed: _reduceCounter,
                  tooltip: 'Reduce',
                  child: Icon(Icons.remove),
                )
              ]),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        heroTag: 'login',
        onPressed: (){
          _methodChannel.invokeMethod("toNative",{"route":"/"}).then((value) => null);
          final user = Get.find<AuthState>().user;
          if ( user!= null) {
              showDialog(
                  context: context,
                  builder: (context) {
                    return AlertDialog(
                      title: Text('提示',textAlign:TextAlign.center),
                      content: Text('确认退出吗？',textAlign:TextAlign.center),
                      actions: <Widget>[
                        TextButton(
                          child: Text('取消'),
                          onPressed: () {
                            Navigator.of(context).pop('cancel');
                          },
                        ),
                        TextButton(
                          child: Text('确认'),
                          onPressed: () {
                            Navigator.of(context).pop('ok');
                          },
                        ),
                      ],
                    );
                  });
            }
          else {Get.to(LoginView());}
        },

        tooltip: 'ToBrowser',
        child: Icon(Icons.send),
      ),
      // This trailing comma makes auto-formatting nicer for build methods.
      floatingActionButtonAnimator: FloatingActionButtonAnimator.scaling,
    );
  }
}
