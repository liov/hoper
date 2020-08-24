import 'dart:convert';

import 'package:app/model/moment/moment.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'dart:io';

class MomentListView extends StatefulWidget {
  MomentListStage createState() => MomentListStage();
}

class MomentListStage extends State<MomentListView> {
  List<Moment> items;

  _getItems() async{
    var url = '/moment';

    List<Moment> result;

    try {
      var response = await Dio().get(url);
      if (response.statusCode == HttpStatus.ok) {
        var data = jsonDecode(response.data);
        result = data['origin'];
      } else {
        result = null;
      }
    } catch (exception) {
      result = null;
    }

    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    if (!mounted) return;

    setState(() {
      items = result;
    });

  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
        itemCount: items.length,
        itemBuilder: (context, index) {
          return ListTile(
            title: Text('${items[index].content}'),
          );
        }
    );
  }

}