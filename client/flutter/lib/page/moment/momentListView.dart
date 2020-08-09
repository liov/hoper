import 'dart:convert';

import 'package:app/generated/user/user.model.pb.dart';
import 'package:app/const/const.dart';
import 'package:flutter/material.dart';
import 'dart:io';

class MomentListView extends StatefulWidget {
  MomentListStage createState() => MomentListStage();
}

class MomentListStage extends State<MomentListView> {
  List<User> items;

  _getItems() async{
    var url = '$baseHost/api/';
    var httpClient = HttpClient();

    List<User> result;
    try {
      var request = await httpClient.getUrl(Uri.parse(url));
      var response = await request.close();
      if (response.statusCode == HttpStatus.ok) {
        var json = await response.transform(utf8.decoder).join();
        var data = jsonDecode(json);
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
            title: Text('${items[index].name}'),
          );
        }
    );
  }

}