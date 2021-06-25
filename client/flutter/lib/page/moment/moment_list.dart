import 'package:app/model/moment.dart';
import 'package:app/model/state/moment.dart';
import 'package:app/model/state/user.dart';
import 'package:app/model/user.dart';
import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:get/get.dart';

import 'moment_item.dart';

class MomentListView extends StatelessWidget {
  final MomentState momentState = Get.put(MomentState());
  final UserState userState = Get.put(UserState());

  _getList() async {
    var response = await getMomentList(
        momentState.pageNo.value, momentState.pageSize.value);
    if (response == null) return;
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    response.users.forEach((e) => userState.users[e.id] = e);
    momentState.list.addAll(response.list);
    momentState.timesIncrement();
    momentState.pageNoIncrement();
  }

  late final ScrollController _controller = ScrollController()
    ..addListener(() {
      if (this._controller.position.atEdge) {
        _getList();
      }
      }
    );

  @override
  Widget build(BuildContext context) {
    var list = momentState.list;
    var users = userState.users;
    _getList();
    return Obx(() => RefreshIndicator(
        onRefresh: () {
          momentState.reset();
          return _getList();
          },
        child:ListView.separated(
        controller: this._controller,
        itemCount: list.length,
        separatorBuilder: (BuildContext context, int index) {
          return Divider();
        },
        itemBuilder: (context, index) {
          return MomentItem(moment:list[index]);
        })));
  }
}
