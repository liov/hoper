import 'package:flutter/material.dart';

class HIcon extends StatefulWidget {
  HIcon({key,required this.icon, required this.label, this.color = Colors.black}):super(key:key);
  final IconData icon;
  final String label;
  final Color color;
  HIconState createState() => HIconState();
}

class HIconState extends State<HIcon> {
  @override
  Widget build(BuildContext context) {
    Color color = widget.color != null? widget.color:Theme.of(context).primaryColor;
    return Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Icon(widget.icon, color: color),
        Container(
          margin: const EdgeInsets.only(top:3.0),
          child: Text(
            widget.label,
            style: TextStyle(
              fontSize: 12.0,
              fontWeight: FontWeight.w400,
              color: color,
            ),
          ),
        ),
      ],
    );
  }
}

class RowIcon extends StatefulWidget {
  RowIconState createState() => RowIconState();
}

class RowIconState extends State<RowIcon> {

  Column buildButtonColumn(IconData icon, String label) {
    Color color = Theme.of(context).primaryColor;

    return Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Icon(icon, color: color),
        Container(
          margin: const EdgeInsets.only(top: 8.0),
          child: Text(
            label,
            style: TextStyle(
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
        color:Theme.of(context).primaryColor.withAlpha(101),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            HIcon(icon:Icons.call, label:'CALL'),
            HIcon(icon:Icons.near_me, label:'ROUTE'),
            HIcon(icon:Icons.share, label:'SHARE'),
            HIcon(icon:Icons.shopping_cart, label:'SHARE'),
          ],
        )
    );
  }
}