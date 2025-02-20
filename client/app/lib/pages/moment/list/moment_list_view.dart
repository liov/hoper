
import 'package:app/components/async/async.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../global/service.dart';
import '../item/moment_item_view.dart';
import 'moment_list_controller.dart';
// TabBarView StatelessWidget无法保持状态
class MomentListView extends StatelessWidget  {
  MomentListView({this.tag = "default"}) : super(key:Key(tag));

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
    globalService.logger.d('MomentListView重绘');
    print(_controller);
    final _future = controller.newList(tag);
    return FutureBuilder<void>(
        future: _future,
        builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
    return snapshot.handle() ?? GetBuilder<MomentListController>(
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
          });
  }
}
