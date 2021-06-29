import 'package:app/generated/protobuf/content/content.model.pb.dart';
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/model/state/user.dart';
import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';

class MomentItem extends StatelessWidget {
  MomentItem({Key? key, required this.moment}) : super(key: key) {
    if (this.moment.images != "")
      this.images = this.moment.images.split(",");
    else this.images = null;
  }

  final Moment moment;
  final RxMap<Int64, UserBaseInfo> users = Get.find<UserState>().users;
  late final List<String>? images;

  @override
  Widget build(BuildContext context) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Row(children: [
        Expanded(
          flex: 1,
          child: Text('${users[moment.userId]!.name}'),
        ),
        Expanded(
          flex: 1,
          child: Text('${moment.createdAt}'),
        ),
      ]),
      MarkdownBody(data: '${moment.content}'),
      if (images != null)
        GridView(
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 3, //横轴三个子widget
              childAspectRatio: 1.0 //宽高比为1时，子widget
              ),
          shrinkWrap: true,
          children: images!.map((e) => Image(image: NetworkImage(e))).toList(),
        ),
      Row(children: [
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
      ])
    ]);
  }
}

