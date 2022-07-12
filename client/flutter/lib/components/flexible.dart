import 'package:flutter/material.dart';

typedef Builder = Widget Function(BuildContext context,State state);

class FlexibleWidget extends StatefulWidget {
  FlexibleWidget(this.builder);
  final Builder builder;
  @override
  State<StatefulWidget> createState() => FlexibleWidgetState(builder);

}

class FlexibleWidgetState extends State<FlexibleWidget>{

  FlexibleWidgetState(this.builder);
  final Builder builder;
  @override
  Widget build(BuildContext context) {
    return builder(context,this);
  }
}