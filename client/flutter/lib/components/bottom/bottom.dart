import 'package:flutter/cupertino.dart';


class Bottom {
 Bottom(this.widget, this.label);

  final Widget widget;

  final String? label;

  factory Bottom.icon(IconData icon, String? label){
    return Bottom(Icon(icon),label);
  }

  BottomNavigationBarItem bottomNavigationBarItem(){
    return BottomNavigationBarItem(
      icon: widget,
      label: label,
    );
  }
}