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

/// Simple builder which extend [DelegateBuilder] to provide some necessary config.
abstract class InnerBuilder extends DelegateBuilder {
  /// List of [TabItem] stands for tabs.
  final List<TabItem> items;

  /// Color used when tab is active.
  final Color activeColor;

  /// Color used for tab.
  final Color color;

  /// Style hook to override the internal tab style
  StyleHook? _style;

  /// Create style builder.
  InnerBuilder(
      {required this.items, required this.activeColor, required this.color});

  /// Get style config
  StyleHook ofStyle(BuildContext context) {
    return StyleProvider.of(context)?.style ?? (_style ??= InternalStyle());
  }

  /// Return true if title text exists
  bool hasNoText(TabItem item) {
    return item.title == null || item.title!.isEmpty;
  }
}


/// Convex shape is fixed center with circle.
class FixedCircleTabStyle extends InnerBuilder  {
  /// Color used as background of appbar and circle icon.
  final Color backgroundColor;

  /// Index of the centered convex shape.
  final int convexIndex;

  /// Create style builder
  FixedCircleTabStyle(
      {required List<TabItem> items,
        required Color activeColor,
        required Color color,
        required this.backgroundColor,
        required this.convexIndex})
      : super();

  @override
  Widget build(BuildContext context, int index, bool active) {
    var c = active ? activeColor : color;
    var item = items[index];
    var style = ofStyle(context);
    var textStyle = style.textStyle(c);
    var margin = style.activeIconMargin;

    if (index == convexIndex) {
      final item = items[index];
      return Container(
        // necessary otherwise the badge will not large enough
        width: style.layoutSize,
        height: style.layoutSize,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          color: c,
        ),
        margin: EdgeInsets.all(margin),
        child: BlendImageIcon(
          active ? item.activeIcon ?? item.icon : item.icon,
          size: style.activeIconSize,
          color: item.blend ? backgroundColor : null,
        ),
      );
    }

    var noLabel = style.hideEmptyLabel && hasNoText(item);
    var icon = BlendImageIcon(
      active ? item.activeIcon ?? item.icon : item.icon,
      color: item.blend ? (c) : null,
      size: style.iconSize,
    );
    var children = noLabel
        ? <Widget>[icon]
        : <Widget>[icon, Text(item.title ?? '', style: textStyle)];
    return Container(
      padding: EdgeInsets.only(bottom: 2),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: children,
      ),
    );
  }

  @override
  bool fixed() {
    return true;
  }
}
