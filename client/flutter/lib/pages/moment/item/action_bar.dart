
import 'package:app/generated/protobuf/content/action.model.pb.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';

class ActionBar extends StatelessWidget{
  ActionBar({this.action,this.ext}):super();

  final UserAction? action;
  final ContentExt? ext;
  static const size = 25.0;

  @override
  Widget build(BuildContext context) {
    return  Padding(
        //上下各添加8像素补白
        padding: const EdgeInsets.symmetric(vertical: 10.0),
        child:Row(children: [
          Expanded(
            flex: 1,
            child:  Row(children: [Expanded(flex:1,child:Icon(Icons.share, color: Colors.green,size:size))]),
          ),
          Expanded(
            flex: 1,
            child:   Row(children: [Expanded(flex:1,child:FaIcon(FontAwesomeIcons.commentAlt,size:size)),Expanded(flex:1,child: Text(ext != null?ext!.commentCount.toStringUnsigned():'0'),)],),
          ),
          Expanded(
            flex: 1,
            child:  Row(children: [Expanded(flex:1,child:Icon(Icons.star, color: Colors.blueAccent[200],size:size)),Expanded(flex:1,child:Text(ext != null?ext!.collectCount.toStringUnsigned():'0'))],),
          ),
          Expanded(
            flex: 1,
            child:   Row(children: [Expanded(flex:1,child:Icon(Icons.favorite, color: Colors.red,size:size)),Expanded(flex:1,child: Text(ext != null?ext!.likeCount.toStringUnsigned():'0'),)],),
          ),
          Expanded(
            flex: 1,
            child:  Icon(Icons.more_horiz_outlined, size:size),
          ),
        ]));
  }

}