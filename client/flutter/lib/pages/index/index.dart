import 'dart:io';

import 'package:animate_do/animate_do.dart';
import 'package:app/components/async/async.dart';
import 'package:app/global/controller.dart';

import 'package:app/pages/home/splash_view.dart';
import 'package:app/routes/route.dart';
import 'package:app/service/baoyu.dart';
import 'package:app/utils/httpserver.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_inappwebview/flutter_inappwebview.dart';
import 'package:get/get.dart';
import 'package:shimmer/shimmer.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:shelf/shelf_io.dart' as shelf_io;
import 'package:shelf_proxy/shelf_proxy.dart';
import '../user/login_view.dart';

class IndexPage extends StatefulWidget {
  @override
  _IndexPageState createState() => _IndexPageState();
}

class _IndexPageState extends State<IndexPage>
    with AutomaticKeepAliveClientMixin {
  final MethodChannel _methodChannel = MethodChannel('xyz.hoper.native/view');

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
    Get.log("IndexPage重绘");
    super.build(context);
    final theme = Theme.of(context);
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
              loop: 1,
              baseColor: theme.textTheme.bodyText2!.color!,
              highlightColor: Colors.blue,
              child: Text(
                '点击下方加号:',
              )),
          Text(
            '$_counter',
            style: theme.textTheme.headline4,
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
                  child: Icon(Icons.add),
                )),
                SizedBox(width: 20),
                FadeInRight(
                    child: FloatingActionButton(
                  heroTag: 'Reduce',
                  onPressed: _reduceCounter,
                  tooltip: 'Reduce',
                  child: Icon(Icons.remove),
                ))
              ]),
          CupertinoSwitch(
              value: Get.isDarkMode,
              onChanged: (value) {
                print(Get.isDarkMode);
                Get.changeThemeMode(
                    Get.isDarkMode ? ThemeMode.light : ThemeMode.dark);
              }),
          Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                      height: 60,
                      width: 150,
                      child: TextField(
                          controller: _controller,
                          focusNode: _focusNode,
                          maxLines: 1,
                          maxLength: 3,
                          decoration: InputDecoration(
                            counterText: '',
                            fillColor: Color(0x30cccccc),
                            filled: true,
                            border: OutlineInputBorder(
                                borderSide:
                                    const BorderSide(color: Colors.blue),
                                borderRadius: const BorderRadius.all(
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
                      icon: Icon(Icons.check))
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
              title: Text('提示', textAlign: TextAlign.center),
              content: Text('确认退出吗？', textAlign: TextAlign.center),
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
        child: Icon(Icons.send),
      ),
      // This trailing comma makes auto-formatting nicer for build methods.
      floatingActionButtonAnimator: FloatingActionButtonAnimator.scaling,
    );
  }

  @override
  bool get wantKeepAlive => true;
}

class MyApp extends StatefulWidget {
  @override
  _MyAppState createState() => new _MyAppState();
}

class _MyAppState extends State<MyApp> {
  final GlobalKey webViewKey = GlobalKey();

  InAppWebViewController? webViewController;
  InAppWebViewGroupOptions options = InAppWebViewGroupOptions(
      crossPlatform: InAppWebViewOptions(
        useShouldOverrideUrlLoading: true,
        mediaPlaybackRequiresUserGesture: false,
      ),
      android: AndroidInAppWebViewOptions(
        useHybridComposition: true,
      ),
      ios: IOSInAppWebViewOptions(
        allowsInlineMediaPlayback: true,
      ));
  late Future<void> _localServer;
  late PullToRefreshController pullToRefreshController;
  String url = "";
  double progress = 0;
  final urlController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _localServer = localServer();
    pullToRefreshController = PullToRefreshController(
      options: PullToRefreshOptions(
        color: Colors.blue,
      ),
      onRefresh: () async {
        if (Platform.isAndroid) {
          webViewController?.reload();
        } else if (Platform.isIOS) {
          webViewController?.loadUrl(
              urlRequest: URLRequest(url: await webViewController?.getUrl()));
        }
      },
    );
  }

  Future<void> localServer() async {
    //Directory.fromRawPath(path).create(recursive:true);
    await shelf_io.serve(
      apiHandler("https://hoper.xyz"),
      'localhost',
      8080,
    );
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        //appBar: AppBar(title: Text("Official InAppWebView website")),
        body: SafeArea(
            child: FutureBuilder<void>(
                future: _localServer,
                builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
                  return snapshot.handle() ??
                      Column(children: <Widget>[
                        Expanded(
                          child: Stack(
                            children: [
                              InAppWebView(
                                key: webViewKey,
                                initialUrlRequest: URLRequest(
                                    url: Uri.parse("http://localhost:8080/index.html")),
                                initialOptions: options,
                                pullToRefreshController:
                                    pullToRefreshController,
                                onWebViewCreated: (controller) {
                                  webViewController = controller;
                                },
                                onLoadStart: (controller, url) {
                                  setState(() {
                                    this.url = url.toString();
                                    urlController.text = this.url;
                                  });
                                },
                                androidOnPermissionRequest:
                                    (controller, origin, resources) async {
                                  return PermissionRequestResponse(
                                      resources: resources,
                                      action: PermissionRequestResponseAction
                                          .GRANT);
                                },
                                shouldOverrideUrlLoading:
                                    (controller, navigationAction) async {
                                  var uri = navigationAction.request.url!;

                                  if (![
                                    "http",
                                    "https",
                                    "file",
                                    "chrome",
                                    "data",
                                    "javascript",
                                    "about"
                                  ].contains(uri.scheme)) {
                                    if (await canLaunch(url)) {
                                      // Launch the App
                                      await launch(
                                        url,
                                      );
                                      // and cancel the request
                                      return NavigationActionPolicy.CANCEL;
                                    }
                                  }

                                  return NavigationActionPolicy.ALLOW;
                                },
                                onLoadStop: (controller, url) async {
                                  pullToRefreshController.endRefreshing();
                                  setState(() {
                                    this.url = url.toString();
                                    urlController.text = this.url;
                                  });
                                },
                                onLoadError: (controller, url, code, message) {
                                  pullToRefreshController.endRefreshing();
                                },
                                onProgressChanged: (controller, progress) {
                                  if (progress == 100) {
                                    pullToRefreshController.endRefreshing();
                                  }
                                  setState(() {
                                    this.progress = progress / 100;
                                    urlController.text = this.url;
                                  });
                                },
                                onUpdateVisitedHistory:
                                    (controller, url, androidIsReload) {
                                  setState(() {
                                    this.url = url.toString();
                                    urlController.text = this.url;
                                  });
                                },
                                onConsoleMessage: (controller, consoleMessage) {
                                  print(consoleMessage);
                                },
                              ),
                              progress < 1.0
                                  ? LinearProgressIndicator(value: progress)
                                  : Container(),
                            ],
                          ),
                        )
                      ]);
                })));
  }

  @override
  void didUpdateWidget(covariant MyApp oldWidget) {
    super.didUpdateWidget(oldWidget);
    webViewController?.clearCache();
    webViewController?.reload();
  }
}
