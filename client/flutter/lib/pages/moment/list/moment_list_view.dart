import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/model/global_state/global_controller.dart';

import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'moment_item_view.dart';
import 'moment_list_controller.dart';


class MomentListView extends StatelessWidget {
  final MomentListController controller = Get.put(MomentListController());
  final GlobalController globalController = Get.find();
  final MomentClient momentClient = Get.put(MomentClient());

  _getList() async {
    var response = await momentClient.getMomentList(
        controller.pageNo, controller.pageSize);
    if (response == null) return;
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    response.users.forEach((e) => globalController.userState.value.users$[e.id] = e);
    controller.list$.addAll(response.list);
    controller.timesIncrement();
    controller.pageNoIncrement();
  }

  _grpcGetList() async {
    print(controller.pageNo);
    var response = await momentClient.stub.list(MomentListReq(pageNo:controller.pageNo, pageSize:controller.pageSize));
    if (response.list.isEmpty) return;
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    response.users.forEach((e) => globalController.userState.value.users[e.id] = e);
    controller.list.addAll(response.list);
    controller.timesIncrement();
    controller.pageNoIncrement();
  }

  late final ScrollController _controller = ScrollController()
    ..addListener(() {
      if (this._controller.position.atEdge) {
        _grpcGetList();
      }
      }
    );

  @override
  Widget build(BuildContext context) {
    _grpcGetList();
    return Obx(() => RefreshIndicator(
        onRefresh: () {
          controller.resetList();
          return _grpcGetList();
          },
        child:ListView.separated(
        controller: this._controller,
        itemCount: controller.list.length,
        separatorBuilder: (BuildContext context, int index) {
          return Divider();
        },
        itemBuilder: (context, index) {
          return MomentItem(moment:controller.list[index]);
        })));
  }
}
