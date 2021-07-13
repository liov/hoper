import 'package:app/generated/protobuf/content/content.model.pb.dart';
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/model/global_state/global_controller.dart';
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
  final GlobalController globalController = Get.find<GlobalController>();
  late final List<String>? images;

  @override
  Widget build(BuildContext context) {
    final user = globalController.userState.getUser(moment.userId);
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Row(children: [
        Expanded(
          flex: 1,
          child: Text('${user!.name}'),
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

    ]);
  }
}

