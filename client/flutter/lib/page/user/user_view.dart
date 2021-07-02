import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'user_logic.dart';
import 'user_state.dart';

class userPage extends StatefulWidget {
  @override
  _userPageState createState() => _userPageState();
}

class _userPageState extends State<userPage> {
  final logic = Get.find<userLogic>();
  final userState state = Get.find<userLogic>().state;

  @override
    Widget build(BuildContext context) {
      return Container();
    }

  @override
  void dispose() {
    Get.delete<userLogic>();
    super.dispose();
  }
}