
import 'package:flutter/material.dart';
import 'package:flutter_spinkit/flutter_spinkit.dart';

final splash = Splash();

class Splash extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SpinKitFadingCircle(
        itemBuilder: (BuildContext context, int index) {
          return DecoratedBox(
            decoration: BoxDecoration(
              color: index.isEven ? Colors.red : Colors.green,
            ),
          );
        },
      ),
      floatingActionButton: Text("跳过"),
    );
  }
}