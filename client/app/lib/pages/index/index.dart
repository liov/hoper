import 'dart:io';

import 'package:animate_do/animate_do.dart';
import 'package:applib/util/async.dart';
import 'package:app/global/state.dart';

import 'package:app/pages/home/splash_view.dart';
import 'package:app/pages/route.dart';
import 'package:app/rpc/baoyu.dart';
import 'package:app/utils/httpserver.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import 'package:shimmer/shimmer.dart';
import 'package:webview_flutter/webview_flutter.dart';

import '../webview/webview.dart';

class IndexPage extends StatefulWidget {
  const IndexPage({super.key});

  @override
  _IndexPageState createState() => _IndexPageState();
}

class _IndexPageState extends State<IndexPage> with AutomaticKeepAliveClientMixin {
  final MethodChannel _methodChannel = const MethodChannel('xyz.hoper.native/view');

  final TextEditingController _controller = TextEditingController();
  final _focusNode = FocusNode();
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
    globalService.logger.d("IndexPage重绘");
    super.build(context);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        // Here we take the value from the MyHomePage object that was created by
        // the App.build method, and use it to set our appbar title.
        title: const Text('🍀'),
      ),
      body: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          Shimmer.fromColors(
              loop: 1,
              baseColor: theme.textTheme.bodyMedium!.color!,
              highlightColor: Colors.blue,
              child: const Text(
                '点击下方加号:',
              )),
          Text(
            '$_counter',
            style: theme.textTheme.headlineMedium,
          ),
          Row(
              mainAxisAlignment: MainAxisAlignment.center,
              verticalDirection: VerticalDirection.up,
              children: [
                FadeInLeft(
                    child: FloatingActionButton(
                  heroTag: 'Increment',
                  onPressed: _incrementCounter,
                  tooltip: 'Increment',
                  child: const Icon(Icons.add),
                )),
                const SizedBox(width: 20),
                FadeIn(
                    child: FloatingActionButton(
                      heroTag: 'Webview',
                      onPressed: ()=>{Get.to(()=> const WebViewExample())},
                      tooltip: 'Webview',
                      child: const Icon(Icons.web),
                    )),
                const SizedBox(width: 20),
                FadeInRight(
                    child: FloatingActionButton(
                  heroTag: 'Reduce',
                  onPressed: _reduceCounter,
                  tooltip: 'Reduce',
                  child: const Icon(Icons.remove),
                ))
              ]),
          CupertinoSwitch(
              value: globalState.isDarkMode.value,
              onChanged: (value) {
                globalState.isDarkMode.toggle();
                Get.changeThemeMode(
                    globalState.isDarkMode.value ? ThemeMode.dark : ThemeMode.light);
              }),
          Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                mainAxisSize: MainAxisSize.min,
                children: [
                  SizedBox(
                      height: 60,
                      width: 150,
                      child: TextField(
                          controller: _controller,
                          focusNode: _focusNode,
                          maxLines: 1,
                          maxLength: 3,
                          decoration: const InputDecoration(
                            counterText: '',
                            fillColor: Color(0x30cccccc),
                            filled: true,
                            border: OutlineInputBorder(
                                borderSide:
                                    BorderSide(color: Colors.blue),
                                borderRadius: BorderRadius.all(
                                    Radius.circular(100))),
                          ))),
                  IconButton(
                      onPressed: () async {
                        final token =
                            await BaoyuClient.signup(_controller.text);
                        _controller.value = TextEditingValue(text: token);
                        _focusNode.unfocus();
                        setState(() {});
                      },
                      color: Colors.blue,
                      icon: const Icon(Icons.check))
                ],
              ))
        ],
      ),
      floatingActionButton: FloatingActionButton(
        heroTag: 'login',
        onPressed: () {
          //_methodChannel.invokeMethod("toNative",{"route":"/"}).then((value) => null);
          final user = globalState.authState.userAuth;
          if (user != null) {
            Get.dialog(AlertDialog(
              title: const Text('提示', textAlign: TextAlign.center),
              content: const Text('确认退出吗？', textAlign: TextAlign.center),
              actions: <Widget>[
                TextButton(
                  child: const Text('取消'),
                  onPressed: () {
                    navigator!.pop('cancel');
                  },
                ),
                TextButton(
                  child: const Text('确认'),
                  onPressed: () {
                    globalState.authState.logout();
                    navigator!.pop('ok');
                  },
                ),
              ],
            ));
          } else {
            Get.toNamed(Routes.LOGIN);
          }
        },
        tooltip: 'ToBrowser',
        child: const Icon(Icons.send),
      ),
      // This trailing comma makes auto-formatting nicer for build methods.
      floatingActionButtonAnimator: FloatingActionButtonAnimator.scaling,
    );
  }

  @override
  bool get wantKeepAlive => true;
}

