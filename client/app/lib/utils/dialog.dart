import 'package:app/global/state.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:get/get.dart';

typedef Callback = void Function();

void dialog(String message, Callback success, Callback cancel) {
  Get.dialog(CupertinoAlertDialog(
    content: Text(message),
    actions: <Widget>[
      CupertinoDialogAction(
        child: Text('取消'),
        onPressed: () {
          navigator!.pop('ok');
        },
      ),
      CupertinoDialogAction(
        child: Text('确认'),
        onPressed: () {
          navigator!.pop('ok');
        },
      ),
    ],
  ));
}

void toast(String message) {
  Widget toast = Center(child:ConstrainedBox(
    constraints: BoxConstraints(maxHeight: 60, maxWidth: 300),
    child: Container(padding: const EdgeInsets.symmetric(horizontal: 24.0, vertical: 12.0),
      decoration: BoxDecoration(
        color: globalState.isDarkMode.value?Colors.black:Colors.blueAccent,
        borderRadius: BorderRadius.circular(25.0),
        //color: Colors.greenAccent,
      ),child:Text(message,textAlign: TextAlign.center,softWrap: true,
            style: TextStyle(color:globalState.isDarkMode.value?Colors.white:Colors.black)
        )),
  ));
  Get.showOverlay(
    opacityColor: Colors.transparent,
      asyncFunction: () async {
        return Future.delayed(Duration(seconds: 1));
      },
      loadingWidget: toast);
}