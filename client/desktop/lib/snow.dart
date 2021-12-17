import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/semantics.dart';


class SnowWidget extends StatefulWidget {
  final int totalSnow;
  final double speed;
  final bool isRunning;

  const SnowWidget(this.totalSnow, this.speed, this.isRunning, {Key? key})
      : super(key: key);

  @override
  _SnowWidgetState createState() => _SnowWidgetState();
}

class _SnowWidgetState extends State<SnowWidget> with SingleTickerProviderStateMixin {
  final Random _rnd = Random();
  late AnimationController controller;
  late Animation animation;
  late List<Snow> _snows;
  double angle = 0;
  double W =  0;
  double H = 300;

  @override
  void initState() {
    super.initState();
    _createSnow();

    controller = AnimationController(
        lowerBound: 0,
        upperBound: 1,
        vsync: this,
        duration: const Duration(milliseconds: 20000));
    controller.addListener(()=>setState(() {update();}));

    if (!widget.isRunning) {
      controller.stop();
    } else {
      controller.repeat();
    }
  }

  @override
  dispose() {
    controller.dispose();
    super.dispose();
  }

  _createSnow() {
    _snows = List.generate(widget.totalSnow, (index) => Snow(_rnd.nextDouble() * W, _rnd.nextDouble() * H,
        _rnd.nextDouble() * 4 + 1, _rnd.nextDouble() * widget.speed));
  }

  update() {
    if (kDebugMode) {
      print(" update" + widget.isRunning.toString());
    }
    angle += 0.01;
    if (widget.totalSnow != _snows.length) {
      _createSnow();
    }
    for (var i = 0; i < widget.totalSnow; i++) {
      var snow = _snows[i];
      //We will add 1 to the cos function to prevent negative values which will lead flakes to move upwards
      //Every particle has its own density which can be used to make the downward movement different for each flake
      //Lets make it more random by adding in the radius
      snow.y += (cos(angle + snow.d) + 1 + snow.r / 2) * widget.speed;
      snow.x += sin(angle) * 2 * widget.speed;
      if (snow.x > W + 5 || snow.x < -5 || snow.y > H) {
        if (i % 3 > 0) {
          //66.67% of the flakes
          _snows[i] = Snow(_rnd.nextDouble() * W, -10, snow.r, snow.d);
        } else {
          //If the flake is exitting from the right
          if (sin(angle) > 0) {
            //Enter from the left
            _snows[i] = Snow(-5, _rnd.nextDouble() * H, snow.r, snow.d);
          } else {
            //Enter from the right
            _snows[i] = Snow(W + 5, _rnd.nextDouble() * H, snow.r, snow.d);
          }
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    W =  MediaQuery.of(context).size.width;
    if (widget.isRunning && !controller.isAnimating) {
      controller.repeat();
    } else if (!widget.isRunning && controller.isAnimating) {
      controller.stop();
    }

    return LayoutBuilder(
      builder: (context, constraints) {
        if (_snows.isEmpty) {
          W = constraints.maxWidth;
          H = constraints.maxHeight;
        }
        return CustomPaint(
          willChange: widget.isRunning,
          painter: SnowPainter(_snows),
          size: Size.infinite,
        );
      },
    );
  }
}

class Snow {
  double x;
  double y;
  double r; //radius
  double d; //density
  Snow(this.x, this.y, this.r, this.d);
}

class SnowPainter extends CustomPainter {
  List<Snow> snows;


  SnowPainter(this.snows);

  @override
  void paint(Canvas canvas, Size size) {
    if (snows.isEmpty) return;
    //draw circle
    var paint = Paint()
      ..color = Colors.white
      ..style = PaintingStyle.fill
      ..strokeWidth = 5;
    for (var snow in snows) {canvas.drawCircle(Offset(snow.x,snow.y), snow.r, paint); }

  }

  @override
  bool shouldRepaint(SnowPainter oldDelegate) => false;
  @override
  bool shouldRebuildSemantics(SnowPainter oldDelegate) => false;
}
