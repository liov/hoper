

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'comment_controller.dart';

class CommentAdd extends StatelessWidget {

  const CommentAdd({super.key});


  @override
  Widget build(BuildContext context) {
    final CommentController controller = Get.find();
    final TextEditingController controller0 = controller.textEditingController;
    final focusNode = controller.focusNode;
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
                controller: controller0,
                focusNode: focusNode,
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
                        focusNode.unfocus();
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
                  if(mode && value.isNotEmpty){
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
    final CommentController controller = Get.find();
    final TextEditingController controller0 = controller.textEditingController;
    final focusNode = controller.focusNode;
    return GetBuilder<CommentController>(
      id:'add',
      builder:(CommentController _){
        return  IconButton(
            onPressed: () async {
              if(controller0.text.isEmpty){
                controller0.text = "为了遇见你我珍惜自己我穿越风和雨是"
                    "为交出我的心，直到遇见你";
                return;
              }

              await controller.save(controller0.text);
              controller0.text = '';
              focusNode.unfocus();
              controller.update(['add']);
            },
            color: Colors.blue,
            icon: controller0.text.isEmpty ?Icon(Icons.add) : Text("发送"),
            tooltip: '更多');
      },);
  }


}
