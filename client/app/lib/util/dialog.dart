import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:app/util/nav.dart';

typedef Callback = void Function();

void dialog(String message, Callback success, Callback cancel) {
  AppNavigator.dialog(CupertinoAlertDialog(
    content: Text(message),
    actions: <Widget>[
      CupertinoDialogAction(
        child: const Text('取消'),
        onPressed: () => AppNavigator.nav!.pop('ok'),
      ),
      CupertinoDialogAction(
        child: const Text('确认'),
        onPressed: () => AppNavigator.nav!.pop('ok'),
      ),
    ],
  ));
}

void toast(String message, {Color color = Colors.blueAccent}) {
  final w = Center(
    child: ConstrainedBox(
      constraints: const BoxConstraints(maxHeight: 60, maxWidth: 300),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
        decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(25)),
        child: Text(message, textAlign: TextAlign.center, softWrap: true, style: TextStyle(color: color)),
      ),
    ),
  );
  AppNavigator.showOverlay(
    loadingWidget: w,
    asyncFunction: () async => Future<void>.delayed(const Duration(seconds: 1)),
  );
}
