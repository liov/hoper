import 'package:app/util/async.dart';

import 'package:app/pages/comment/comment_add_view.dart';
import 'package:app/pages/comment/comment_list_view_v2.dart';
import 'package:app/pages/moment/item/moment_item_view.dart';
import 'package:app/global/service.dart';
import 'package:flutter/material.dart';
import 'package:app/gen/pb/content/moment.model.pb.dart' as $pb;
import 'package:fixnum/fixnum.dart';
import 'package:app/gen/pb/hopeio/request/param.pb.dart' as $1;

class MomentDetailView extends StatelessWidget {
  MomentDetailView.detail($pb.Moment moment, {super.key})
      : id = moment.id,
        future = Future.value(moment);

  MomentDetailView.byId(this.id, {super.key}) : future = _loadMoment(id);

  final Int64 id;
  final Future<$pb.Moment> future;

  static Future<$pb.Moment> _loadMoment(Int64 id) async {
    return globalService.momentClient.stub.info($1.Id(id: id));
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<$pb.Moment>(
      future: future,
      builder: (BuildContext context, AsyncSnapshot<$pb.Moment> snapshot) {
        final noReady = snapshot.handle();
        if (noReady != null) return Scaffold(body: noReady);
        final moment = snapshot.data!;
        return SafeArea(
          child: Scaffold(
            appBar: AppBar(centerTitle: true, title: const Text('瞬间')),
            body: Center(
              child: Column(
                children: [
                  MomentItem(moment: moment),
                  Expanded(
                    flex: 10,
                    child: CommentListViewV2(moment.statistics),
                  ),
                  const Expanded(flex: 1, child: Text('')),
                ],
              ),
            ),
            bottomSheet: const CommentAdd(),
          ),
        );
      },
    );
  }
}
