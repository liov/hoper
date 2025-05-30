import 'package:app/generated/protobuf/content/content.model.pb.dart';
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/global/const.dart';
import 'package:app/global/state.dart';

import 'package:app/pages/photo/photo.dart';
import 'package:app/pages/photo/slide_photo.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:get/get.dart';

import 'package:app/pages/action_bar/action_bar.dart';

import 'package:applib/util/time.dart';

import 'package:app/generated/protobuf/content/moment.model.pb.dart';

class MomentItem extends StatelessWidget {
  MomentItem({Key? key, required this.moment}) : super(key: key) {
    if (moment.images != "") {
      images = moment
          .images
          .map((url) => BASE_STATIC_URL + url)
          .toList();
    } else {
      images = null;
    }
  }

  final Moment moment;

  late final List<String>? images;

  @override
  Widget build(BuildContext context) {
    final user = globalState.userState.getUser(moment.userId);
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Row(children: [
        Expanded(
            flex: 1,
            child: Row(children: [
              Padding(
                  padding: const EdgeInsets.all(10.0),
                  child: GestureDetector(
                    child: CircleAvatar(
                      child: ExtendedImage.network(
                        BASE_STATIC_URL + user!.avatar,
                        alignment: Alignment.centerLeft,
                        fit: BoxFit.fill,
                        shape: BoxShape.circle,
                        cache: true,
                      ),
                      /* backgroundImage: ExtendedNetworkImageProvider(
                BASE_STATIC_URL+user!.avatarUrl,
                cache: true,
              ),*/
                    ),
                    onTap: () =>
                        slidePhotoRoute(BASE_STATIC_URL + user.avatar),
                  )),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [Text(user.name), Text(getDateTime(moment.createdAt.seconds.toInt(),moment.createdAt.nanos.toInt()).toString())],
              ),
            ])),
        Expanded(
            flex: 1,
            child: Align(
                alignment: Alignment.centerRight,
                child: Padding(
                    padding: const EdgeInsets.all(10.0),
                    child: ElevatedButton(
                      style: ButtonStyle(
                          fixedSize:WidgetStateProperty.all(Size.fromHeight(5.0)),
                          backgroundColor:WidgetStateProperty.resolveWith((states)=>Colors.transparent),
                          shape: WidgetStateProperty.resolveWith((states)=>const StadiumBorder(
                              side: BorderSide(color: Colors.blue)))),
                      child: const Text('+关注'),
                      onPressed: () {},
                    )))),
      ]),
      Padding(
          padding: const EdgeInsets.symmetric(vertical:5.0,horizontal:10),
          child: MarkdownBody(data: moment.content)),
      if (images != null)
        GridView.builder(
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 3, //横轴三个子widget
              childAspectRatio: 1.0, //宽高比为1时，子widget
            ),
            shrinkWrap: true,
            itemCount: images!.length,
            itemBuilder: (BuildContext context, int index) {
              return Padding(
                  padding: const EdgeInsets.all(10.0),
                  child: GestureDetector(
                    child: ExtendedImage.network(
                      images![index],
                      alignment: Alignment.centerLeft,
                      fit: BoxFit.fill,
                      cache: true,
                      //cancelToken: cancellationToken,
                    ),
                    onTap: () => slidePhotoRoute(images![index]),
                  ));
            }),
      ActionBar(moment),
    ]);
  }
}
