import 'package:app/components/async/async.dart';
import 'package:app/generated/protobuf/content/content.model.pb.dart' as $moment;
import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/pages/home/global/global_state/global_controller.dart';

import 'package:app/pages/moment/item/moment_item_view.dart';
import 'package:app/pages/moment/detail/moment_detail_view.dart';
import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';



class MomentListV2View extends StatefulWidget {
  MomentListV2View({this.tag = "default"}) : super();

  final String tag;
  @override
  _MomentListV2ViewState createState() => _MomentListV2ViewState();
}

class _MomentListV2ViewState extends State<MomentListV2View> with AutomaticKeepAliveClientMixin {

  final MomentClient momentClient = Get.find();
  final GlobalController globalController = Get.find();

  late final req = MomentListReq(pageNo:1,pageSize:10);
  var times = 0;
  var list = List<$moment.Moment>.empty(growable: true);


  late final ScrollController _controller = ScrollController()
    ..addListener(() {
      if (_controller.position.atEdge) {
        grpcGetList();
      }
    }
    );

  Future<void> resetList() {
    list.removeRange(0, list.length);
    req.pageNo = 1;
    return grpcGetList();
  }

  Future<void> grpcGetList() async {
    var response = await momentClient.stub.list(req);
    if (response.list.isEmpty) return;
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    globalController.userState.appendUsers(response.users);
    list.addAll(response.list);
    times++;
    req.pageNo++;
    setState(() {});
  }

  @override
    Widget build(BuildContext context) {
      super.build(context);
      final _future = grpcGetList();
      return FutureBuilder<void>(
          future: _future,
          builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
            return snapshot.handle() ??  RefreshIndicator(
                        onRefresh: () {
                          return resetList();
                        },
                        child: ListView.separated(
                            physics: BouncingScrollPhysics(),
                            controller: _controller,
                            itemCount: list.length,
                            separatorBuilder: (BuildContext context, int index) {
                              return Divider();
                            },
                            itemBuilder: (context, index) {
                              return InkWell(
                                onTap: (){
                                  Get.to(()=>MomentDetailView(moment:list[index],id:list[index].id));
                                },
                                child: MomentItem(
                                    moment: list[index]),
                              );
                            }));
            });
    }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  bool get wantKeepAlive => true;
}