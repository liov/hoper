import 'package:flutter/material.dart';

class Bottom extends StatefulWidget {
  BottomState createState() => new BottomState();
}

class BottomState extends State<Bottom> {
  Column buildButtonColumn(IconData icon, String label) {
    Color color = Theme.of(context).primaryColor;

    return new Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        new Icon(icon, color: color),
        new Container(
          margin: const EdgeInsets.only(top: 8.0),
          child: new Text(
            label,
            style: new TextStyle(
              fontSize: 12.0,
              fontWeight: FontWeight.w400,
              color: color,
            ),
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            buildButtonColumn(Icons.call, 'CALL'),
            buildButtonColumn(Icons.near_me, 'ROUTE'),
            buildButtonColumn(Icons.share, 'SHARE'),
            buildButtonColumn(Icons.shopping_cart, 'SHARE'),
          ],
        )
    );
  }
}
