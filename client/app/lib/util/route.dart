
import 'package:flutter/material.dart';

class SimpleRoute<T> extends PageRoute<T> with MaterialRouteTransitionMixin<T> {
  /// Construct a MaterialPageRoute whose contents are defined by [builder].
  SimpleRoute({
    required this.builder,
    super.settings,
    super.requestFocus,
    this.maintainState = true,
    super.fullscreenDialog,
    super.allowSnapshotting = true,
    super.barrierDismissible = false,
    super.traversalEdgeBehavior,
    super.directionalTraversalEdgeBehavior,
  }) ;

  @override
  Widget buildContent(BuildContext context) => builder(context);

  @override
  final bool maintainState;

  @override
  String get debugLabel => '${super.debugLabel}(${settings.name})';

  final WidgetBuilder builder;

  @override
  String? get barrierLabel => null;

  @override
  bool get opaque => false;

  @override
  Color? get barrierColor => null;
}