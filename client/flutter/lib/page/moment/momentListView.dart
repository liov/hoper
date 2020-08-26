import 'package:app/model/moment.dart';
import 'package:app/service/moment.dart';
import 'package:app/util/dio.dart';
import 'package:flutter/material.dart';
import 'dart:io';

import 'package:flutter_markdown/flutter_markdown.dart';

class MomentListView extends StatefulWidget {
  MomentListStage createState() => MomentListStage();
}

class MomentListStage extends State<MomentListView> {
  List<Moment> items;
  int pageNo = 0;
  int pageSize = 10;

  _getItems() async{

    var response = await getMomentList(pageNo, pageSize);
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    if (!mounted) return;

    setState(() {
      if(items == null) items = response.data;
        else items.addAll(response.data);
    });

  }

  initState(){
    super.initState();
    _getItems();
  }

  @override
  Widget build(BuildContext context) {
    return ListView.separated(
        itemCount: items!=null?items.length:0,
        separatorBuilder: (BuildContext context, int index){
          return Divider();
        },
        itemBuilder: (context, index) {
          return Column(
            children:[
              Row(
                children:[
                  Text('${items[index].user.name}  ${items[index].createdAt}'),
                ]
              ),
              MarkdownBody(data:'${items[index].content}'),
              Row(
              children:[
                Text('收藏：'),
                Icon(Icons.star,color: Colors.blue,),
                Text('喜欢'),
                Icon(Icons.favorite,color: Colors.red,),
              ]
              )
            ]
          );
        }
    );
  }

}