
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../item/moment_item_view.dart';
import 'moment_list_controller.dart';
// TabBarView 无法保持状态
class MomentListView extends StatelessWidget {
  MomentListView({this.tag = "default"}) : super();

  final MomentListController controller = Get.find();

  final String tag;

  late final ScrollController _controller = ScrollController()
    ..addListener(() {
      if (_controller.position.atEdge) {
        controller.pullList(tag);
      }
    }
    );
  @override
  Widget build(BuildContext context) {
    print(controller.entityMap[tag]?.list.length);
    Get.log('重建');
    final _future = controller.newList(tag);
    return FutureBuilder<void>(
        future: _future,
        builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
          switch (snapshot.connectionState) {
            case ConnectionState.none:
              return Text('还没有开始网络请求');
            case ConnectionState.active:
              return Text('ConnectionState.active');
            case ConnectionState.waiting:
              return Center(
                child: CircularProgressIndicator(),
              );
            case ConnectionState.done:
              return GetBuilder<MomentListController>(
                  id: tag,
                  builder: (_) => RefreshIndicator(
                      onRefresh: () {
                        return controller.resetList(tag);
                      },
                      child: ListView.separated(
                          physics: BouncingScrollPhysics(),
                          controller: _controller,
                          itemCount: controller.entityMap[tag]!.list.length,
                          separatorBuilder: (BuildContext context, int index) {
                            return Divider();
                          },
                          itemBuilder: (context, index) {
                            return MomentItem(
                                moment: controller.entityMap[tag]!.list[index]);
                          })));
          }
        });
  }
}
