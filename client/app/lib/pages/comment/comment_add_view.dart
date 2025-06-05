import 'dart:io';


import 'package:app/generated/protobuf/content/action.service.pb.dart';
import 'package:app/pages/moment/add/moment_add_controller.dart';
import 'package:app/pages/moment/add/moment_add_view.dart';
import 'package:app/pages/image/slide_image.dart';
import 'package:app/rpc/action.dart';
import 'package:app/rpc/baoyu.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';

import 'comment_controller.dart';

class CommentAdd extends StatelessWidget {
  final CommentController controller = Get.find();
  late final TextEditingController _controller = controller.textEditingController;
  late final _focusNode = controller.focusNode;


  @override
  Widget build(BuildContext context) {
    print('@'.codeUnits);
    var mode = true;
    return
        Row(mainAxisSize: MainAxisSize.min, children: [
          Padding(
            padding: EdgeInsets.symmetric(horizontal: 10),
          ),
          Expanded(
              flex: 7,
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
                    borderSide:
                        const BorderSide(width: 0, color: Colors.transparent),
                    borderRadius: const BorderRadius.all(Radius.circular(10)),
                  ),
                  contentPadding:
                      EdgeInsets.symmetric(vertical: 0, horizontal: 10),
                  suffixIcon: IconButton(
                      onPressed: () async {
                        _focusNode.unfocus();
                        showBottomSheet( context: context,
                          builder: (BuildContext context) {
                          return Text('测试');
                        },
                        );
                      },
                      color: Colors.blue,
                      icon: Icon(Icons.mood),
                      tooltip: '发送'),
                ),
                onChanged: (value){
                  if(!mode && value.isEmpty){
                    mode = true;
                    controller.update(['add']);
                  }
                  if(mode && !value.isEmpty){
                    mode = false;
                    controller.update(['add']);
                  }
                },
              )),
          Expanded(
            flex: 1,
            child: button1(),
          )


      ],
    );
  }

  Widget button1(){
    return GetBuilder<CommentController>(
      id:'add',
      builder:(CommentController _){
        return  IconButton(
            onPressed: () async {
              if(_controller.text.isEmpty){
                _controller.text = "为了遇见你我珍惜自己我穿越风和雨是"
                    "为交出我的心，直到遇见你";
                return;
              }

              await controller.save(_controller.text);
              _controller.text = '';
              _focusNode.unfocus();
              controller.update(['add']);
            },
            color: Colors.blue,
            icon: _controller.text.isEmpty ?Icon(Icons.add) : Text("发送"),
            tooltip: '更多');
      },);
  }


}
