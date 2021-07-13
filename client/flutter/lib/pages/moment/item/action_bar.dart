
import 'package:flutter/material.dart';

class ActionBar extends StatelessWidget{
  @override
  Widget build(BuildContext context) {
    return Row(children: [
      Expanded(
        flex: 1,
        child:  Icon(Icons.more, color: Colors.yellowAccent[700],),
      ),
      Expanded(
        flex: 1,
        child:  Icon(Icons.favorite, color: Colors.red,),
      ),
      Expanded(
        flex: 1,
        child:  Icon(Icons.star, color: Colors.blueAccent[200],),
      ),
      Expanded(
        flex: 1,
        child:  Icon(Icons.share, color: Colors.green,),
      ),
    ]);
  }

}