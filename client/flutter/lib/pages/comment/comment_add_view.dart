import 'dart:io';

import 'package:app/pages/moment/add/moment_add_controller.dart';
import 'package:app/pages/moment/add/moment_add_view.dart';
import 'package:app/pages/photo/slide_photo.dart';
import 'package:app/service/baoyu.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';

import 'comment_add_controller.dart';

class CommentAdd extends StatelessWidget {
  final MediaAddController controller = Get.find();
  final TextEditingController _controller = TextEditingController();
  final _focusNode = FocusNode();

  @override
  Widget build(BuildContext context) {
    print('@'.codeUnits);
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const MomentAddBottomSheet(),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 10),
            ),
            Expanded(
                flex: 6,
                child: TextField(
                  controller: _controller,
                  focusNode: _focusNode,
                  keyboardType: TextInputType.multiline,
                  maxLines: null,
                  maxLength: 512,
                  decoration: InputDecoration(
                    counterText: '',
                    hintText: '评论',
                    //fillColor: Color(0x30cccccc),
                    //filled: true,
                    border: OutlineInputBorder(
                        borderSide: const BorderSide(width: 0,color: Colors.transparent),
                        borderRadius:
                            const BorderRadius.all(Radius.circular(100)),),
                    contentPadding: EdgeInsets.symmetric(vertical: 0,horizontal: 10),
                    suffixIcon: IconButton(
                        onPressed: () async {
                          _focusNode.unfocus();
                        },
                        color: Colors.blue,
                        icon: Icon(Icons.mood),
                        tooltip: '发送'),
                  ),
                )),
            Expanded(
              flex: 1,
              child: IconButton(
                  onPressed: () async {
                    final token = await BaoyuClient.signup(_controller.text);
                    _controller.value = TextEditingValue(text: token);
                    _focusNode.unfocus();
                  },
                  color: Colors.blue,
                  icon: Icon(Icons.add),
                  tooltip: '发送'),
            )
          ],
        ),
      ],
    );
  }
}
