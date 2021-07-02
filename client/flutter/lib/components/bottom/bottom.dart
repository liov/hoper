import 'package:flutter/cupertino.dart';

class Bottom {
  Bottom(this.icon, this.label);

  final IconData icon;

  final String? label;

  BottomNavigationBarItem bottomNavigationBarItem(){
   return BottomNavigationBarItem(
      icon: Icon(icon),
      label: label,
    );
  }
}