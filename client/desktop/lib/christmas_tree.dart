import 'dart:async';

import 'package:desktop/snow.dart';
import 'package:flutter/material.dart';

class ChristmasTree extends StatefulWidget {
  const ChristmasTree({Key? key}) : super(key: key);

  @override
  _ChristmasTreeState createState() => _ChristmasTreeState();
}

class _ChristmasTreeState extends State<ChristmasTree>
    with TickerProviderStateMixin {
  int branches = 7;
  late AnimationController controller;
  late Animation colorAnimation;

  @override
  void initState() {
    super.initState();
    controller =
        AnimationController(vsync: this, duration: const Duration(milliseconds: 350));

    colorAnimation =
        ColorTween(begin: Colors.blue, end: Colors.red).animate(controller);
    controller.repeat(reverse: true);
    controller.addListener(()=>setState((){}));
  }

  @override
  void dispose() {
    controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final width=MediaQuery.of(context).size.width;
    final height=MediaQuery.of(context).size.height;
    return Scaffold(
      backgroundColor: Colors.black,
      body: SingleChildScrollView(
        child: Stack(
          children: [
            SizedBox(
                width: width,
                height: height,
                child:buildTree()),
              SizedBox(
                  width: width,
                  height: height,
                  child:const SnowWidget(100,1,true)
              ),
            ],),
          ),
    );
  }

  buildTree() {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        buildStar(),
        for (var i = 1; i < branches; i++) buildRow(i),
        buildBase(),
        buildBark(),
      ],
    );
  }

  buildRow(int i) {
    return Wrap(
      children: [
        for (var j = 0; j <= i; j++)
          Text(
            " * ",
            style: TextStyle(
              fontSize: 50,
              color: colorAnimation.value,
            ),
          ),
      ],
    );
  }

  Widget buildBase() {
    return Container(
      width: (branches-1).toDouble() * 50,
      height: 3.5,
      color: Colors.lightGreen,
    );
  }

  Widget buildBark() {
    return Container(
      width: 30,
      height: 100,
      color: Colors.green[700],
    );
  }
  Widget buildStar(){
    return const Icon(Icons.star,size:100,color: Colors.yellow,);
  }
}