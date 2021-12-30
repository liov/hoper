import 'package:convex_bottom_bar/convex_bottom_bar.dart';
import 'package:flutter/cupertino.dart';


class Bottom {
 Bottom(this.widget, this.label,this.pageIndex,this.onTap);

  final Widget widget;

  final String? label;
  final int? pageIndex;
 final Function? onTap;

  factory Bottom.icon(IconData icon, {String? label = "",int? pageIndex,Function? onTap}){
    return Bottom(Icon(icon),label,pageIndex,onTap);
  }

  BottomNavigationBarItem navigationBarItem(){
    return BottomNavigationBarItem(
      icon: widget,
      label: label,
    );
  }

 TabItem tabItem(){
   return TabItem(
     icon: widget,
     title: label,
   );
 }
}