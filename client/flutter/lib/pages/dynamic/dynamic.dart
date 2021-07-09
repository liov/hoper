import 'package:flutter_js/flutter_js.dart';
import 'package:flutter/material.dart';


class Dynamic extends StatelessWidget {
  final JavascriptRuntime flutterJs = getJavascriptRuntime();

  @override
  Widget build(BuildContext context) {
    JsEvalResult jsResult = flutterJs.evaluate("Date();");
    return  Text(jsResult.stringResult);
  }
}
