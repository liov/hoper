import 'package:applib/util/async.dart';
import 'package:app/generated/protobuf/content/action.model.pb.dart';
import 'package:app/generated/protobuf/content/action.service.pb.dart';
import 'package:app/pages/comment/comment_controller.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get_core/src/get_main.dart';
import 'package:get/get_instance/src/extension_instance.dart';
import 'package:get/get_state_manager/src/simple/get_controllers.dart';
import 'package:get/get_state_manager/src/simple/get_state.dart';

import 'comment_item_view.dart';

class CommentListViewV2 extends StatelessWidget {
  CommentListViewV2(this.ext) : super();
  final Statistics ext;
  late final CommentController controller = Get.find()..init(ext.type, ext.refId);
  late final ScrollController scrollController = ScrollController()
    ..addListener(() {
      if (scrollController.position.atEdge) {
        controller.grpcGetList();
      }
    }
    );
  @override
  Widget build(BuildContext context) {
    return GetBuilder(
      id: "list",
      builder: (CommentController _) {
        return FutureBuilder<void>(
            future: controller.future,
            builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
              return snapshot.handle() ??  RefreshIndicator(
                  onRefresh: () {
                    return controller.resetList();
                  },
                  child: controller.list.isEmpty? Center(child:TextButton.icon(onPressed: () {
                    controller.resetList();
                    },
                    label: const Text('暂无评论，点击刷新'),
                    icon: const Icon(Icons.refresh),)) : ListView.separated(
                      //shrinkWrap: true,
                      physics: const BouncingScrollPhysics(),
                      controller: scrollController,
                      itemCount: controller.list.length,
                      separatorBuilder: (BuildContext context, int index) {
                        return Divider();
                      },
                      itemBuilder: (context, index) {
                        return InkWell(
                          onTap: (){

                            final tmp = controller.list[index];
                            controller.rootId = tmp.rootId == 0
                                ? tmp.id : tmp.rootId;
                            controller.recvId = tmp.userId;
                            controller.replyId = tmp.id;
                            if(controller.focusNode.hasFocus){
                              SystemChannels.textInput.invokeListMethod('TextInput.show');
                            }
                            controller.focusNode.requestFocus();

                          },
                          child: CommentItem(
                              comment: controller.list[index]),
                        );
                      }));
            });
      },
    );
  }
}
