
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get_core/src/get_main.dart';
import 'package:get/get_navigation/src/extension_navigation.dart';

void dialog(String message){
  Get.dialog(CupertinoAlertDialog(
    content: Text(message),
    actions: <Widget>[
      TextButton(
        child: Text('确认'),
        onPressed: () {
          navigator!.pop('ok');
        },
      ),
    ],
  ));
}