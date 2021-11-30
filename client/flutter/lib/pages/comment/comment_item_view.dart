import 'package:app/generated/protobuf/content/action.model.pb.dart';
import 'package:app/global/controller.dart';
import 'package:app/model/const/const.dart';
import 'package:app/pages/action_bar/action_bar.dart';
import 'package:app/pages/photo/slide_photo.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';

class CommentItem extends StatelessWidget {
  CommentItem({Key? key, required this.comment}) : super(key: key) {
    if (this.comment.image != "")
      this.images = this
          .comment
          .image
          .split(",")
          .map((url) => BASE_STATIC_URL + url)
          .toList();
    else
      this.images = null;
  }

  final Comment comment;

  late final List<String>? images;

  @override
  Widget build(BuildContext context) {
    print("${this.toStringShort()}");
    final user = globalState.userState.getUser(comment.userId);
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
                        BASE_STATIC_URL + user!.avatarUrl,
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
                        slidePhotoRoute(BASE_STATIC_URL + user.avatarUrl),
                  )),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [Text('${user.name}'), Text('${comment.createdAt}')],
              ),
            ])),
        Expanded(
            flex: 1,
            child: Align(
                alignment: Alignment.centerRight,
                child: Padding(
                    padding: const EdgeInsets.all(10.0),
                    child: IconButton(icon: Icon(Icons.favorite), onPressed: () {  },

                    )))),
      ]),
      Padding(
          padding: const EdgeInsets.symmetric(vertical:5.0,horizontal:10),
          child: MarkdownBody(data: '${comment.content}')),
      if (images != null)
        GridView.builder(
            gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
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
    ]);
  }
}
