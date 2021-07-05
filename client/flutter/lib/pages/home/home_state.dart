
import 'package:app/components/bottom/bottom.dart';
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:flutter/material.dart';
import 'package:fixnum/fixnum.dart';

import 'package:app/model/user.dart' as $self;
import 'package:get/get_rx/src/rx_types/rx_types.dart';

class HomeState {
  var selectedIndex = 0;
  var now = DateTime.now();

  var bottomNavigationBarList = [
    Bottom(Icons.home, "home"),
    Bottom(Icons.account_box_rounded, "Profile"),
    Bottom(Icons.account_box_rounded, "Moments"),
    Bottom(Icons.account_balance_sharp, "Weather"),
  ];


  HomeState() {
    ///Initialize variables
  }
}
