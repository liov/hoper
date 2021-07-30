import 'package:app/generated/protobuf/content/content.model.pb.dart';
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/model/const/const.dart';
import 'package:app/pages/home/global/global_controller.dart';

import 'package:app/pages/photo/photo.dart';
import 'package:app/pages/photo/slide_photo.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';

import 'action_bar.dart';

class MomentItem extends StatelessWidget {
  MomentItem({Key? key, required this.moment}) : super(key: key) {
    if (this.moment.images != "")
      this.images = this.moment.images.split(",").map((url) => BASE_HOST+"/static/"+url).toList();
    else this.images = null;
  }

  final Moment moment;

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
        GridView.builder(
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 3, //横轴三个子widget
              childAspectRatio: 1.0 //宽高比为1时，子widget
              ),
          shrinkWrap: true,
          itemCount:images!.length,
          itemBuilder: (BuildContext context, int index) {
            return  GestureDetector(
                  child:ExtendedImage.network(
                          images![index],
                          alignment: Alignment.centerLeft,
                          fit: BoxFit.fill,
                          cache: true,
                          //cancelToken: cancellationToken,
                  ),
                onTap:()=>slidePhotoRoute(images![index]),
        );}),
      ActionBar(),
    ]);
  }
}

