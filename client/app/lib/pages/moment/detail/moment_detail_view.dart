import 'package:app/components/async/async.dart';
import 'package:app/generated/protobuf/content/action.enum.pb.dart';
import 'package:app/generated/protobuf/content/content.enum.pb.dart';
import 'package:app/pages/comment/comment_add_view.dart';
import 'package:app/pages/comment/comment_controller.dart';
import 'package:app/pages/comment/comment_list_view_v2.dart';
import 'package:app/pages/moment/item/moment_item_view.dart';
import 'package:app/routes/route.dart';
import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:app/generated/protobuf/content/content.model.pb.dart' as $pb;
import 'package:fixnum/fixnum.dart';
import 'package:app/generated/protobuf/tiga/protobuf/request/param.pb.dart' as $1;

class MomentDetailView extends StatelessWidget {
  final CommentController commentController = Get.find();

  MomentDetailView({super.key}){
    if (Get.arguments != null) {
      moment = Get.arguments;
      commentController.refId = moment.id;
      commentController.recvId = moment.userId;
      commentController.type = ContentType.ContentMoment;
      future = Future.value(moment);
      return;
    }

    final idStr = Get.parameters['id'];
    if (idStr != null) {
      id = Int64.parseInt(idStr);
      future = getMoment();
      return;
    }
    Get.toNamed(Routes.NOTFOUND);
  }

  MomentDetailView.detail(this.moment) : super() {
    id = moment.id;
    future = Future.value(moment);
  }

  MomentDetailView.byId(this.id) : super() {
    future = getMoment();
  }

  final MomentClient momentClient = Get.find();
  late final $pb.Moment moment;
  late final Int64 id;
  late final Future<$pb.Moment> future;

  Future<$pb.Moment> getMoment() async {
    final rpcMoment = await momentClient.stub.info($1.Id(id: id));
    moment = rpcMoment;
    return rpcMoment;
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
              appBar: AppBar(
                centerTitle: true,
                title: const Text('瞬间'),
              ),
              body: Center(
                child: Column(
                  children: [
                    MomentItem(moment: moment),
                    Expanded(flex: 10, child: CommentListViewV2(moment.ext)),
                    const Expanded(flex: 1, child: Text('')),
                  ],
                ),
              ),
              bottomSheet: CommentAdd(),
            ),
          );
        });
  }
}
