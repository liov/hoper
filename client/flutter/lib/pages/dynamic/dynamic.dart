import 'dart:async';
import 'dart:convert';

import 'package:app/components/async/async.dart';
import 'package:app/components/json_viewer.dart';
import 'package:app/global/global_service.dart';
import 'package:app/utils/dialog.dart';
import 'package:flutter_js/flutter_js.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:get/get_core/src/get_main.dart';
import 'package:get/get_navigation/src/extension_navigation.dart';

class Dynamic extends StatelessWidget {
  final JavascriptRuntime flutterJs = getJavascriptRuntime(extraArgs:{"stackSize":1024*1024*5})..evaluate("""
  var alert = function(){
    sendMessage('Alert', JSON.stringify(['alert', ...arguments]));
  }
  """)..onMessage("Alert",(dynamic args) {
    toast("${args[1]}");
  });

  Future<String> exec({String path="assets/js/test.js"}) async {
    final body = await rootBundle.loadString(path);
    globalService.logger.d("js:" + body);
    JsEvalResult jsResult = await flutterJs.evaluateAsync(body);
    return jsResult.stringResult;
  }

  late final future = getMoment2();

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<dynamic>(
        future: future,
        builder: (BuildContext context, AsyncSnapshot<dynamic> snapshot) {
          final noReady = snapshot.handle();
          if (noReady != null) return Container();
          return Center(
            child: JsonViewerRoot(jsonObj:jsonDecode(snapshot.data! as String)),
          );
        });
  }

  Future<dynamic> getMoment() {
    final completer = Completer();
    final code = """
    console.log('开始请求');
    async function getMoment(){
      try {
        const response = await fetch('https://hoper.xyz/api/v1/moment?pageNo=1&pageSize=10');
        const body = await response.json();
        if (response.status === 200) {
          sendMessage('onRequestSuccess', JSON.stringify(body));
        } else {
          sendMessage('onRequestFailure', JSON.stringify(body));
        }
      } catch(e) {
        console.log(e.message);
        sendMessage('onError', e.message);
      }
    }
    getMoment();
  """;
    JsEvalResult jsResult = flutterJs.evaluate(code);
    flutterJs.executePendingJob();
    flutterJs.onMessage('onRequestSuccess', completer.complete);
    flutterJs.onMessage('onRequestFailure', (args) {
      completer.completeError(args);
    });
    flutterJs.onMessage('onError', (args) {
      completer.completeError(args);
    });
    return completer.future;
  }

  Future<String> getMoment2() async{

    final code = """
   fetch('https://hoper.xyz/api/v1/moment?pageNo=1&pageSize=1').then(response => response.text());
  """;
    var asyncResult = await flutterJs.evaluateAsync(code);
    //flutterJs.executePendingJob();
    final promiseResolved = await flutterJs.handlePromise(asyncResult);
    return promiseResolved.stringResult;
  }
}
