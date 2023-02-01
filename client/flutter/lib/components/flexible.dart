import 'package:flutter/material.dart';

typedef Builder = Widget Function(BuildContext context,State state);

class FlexibleWidget extends StatefulWidget {
  const FlexibleWidget(this.builder, {super.key});
  final Builder builder;
  @override
  State<StatefulWidget> createState() => FlexibleWidgetState();

}

class FlexibleWidgetState extends State<FlexibleWidget>{

  FlexibleWidgetState();

  @override
  Widget build(BuildContext context) {
    return widget.builder(context,this);
  }
}