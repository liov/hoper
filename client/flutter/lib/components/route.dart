
import 'package:flutter/cupertino.dart';

class SimpleRoute extends PageRoute {
  SimpleRoute({
    this.name,
    required this.builder,
  }) : super(
    settings: RouteSettings(name: name),
  );


  final String? name;
  final WidgetBuilder builder;

  @override
  String? get barrierLabel => null;

  @override
  bool get opaque => false;

  @override
  bool get maintainState => true;

  @override
  Duration get transitionDuration => Duration(milliseconds: 0);

  @override
  Widget buildPage(
      BuildContext context,
      Animation<double> animation,
      Animation<double> secondaryAnimation,
      ) {
    return builder(context);
  }

  /// 页面切换动画
  @override
  Widget buildTransitions(
      BuildContext context,
      Animation<double> animation,
      Animation<double> secondaryAnimation,
      Widget child,
      ) {
    return FadeTransition(
      opacity: animation,
      child: child,
    );
  }

  @override
  Color? get barrierColor => null;
}