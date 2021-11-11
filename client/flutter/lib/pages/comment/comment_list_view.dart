import 'package:app/components/async/async.dart';
import 'package:app/generated/protobuf/content/action.model.pb.dart';
import 'package:app/generated/protobuf/content/action.service.pb.dart';
import 'package:app/generated/protobuf/content/content.model.pb.dart' as $moment;
import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/global/controller.dart';

import 'package:app/pages/moment/item/moment_item_view.dart';
import 'package:app/pages/moment/detail/moment_detail_view.dart';
import 'package:app/service/action.dart';
import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'comment_item_view.dart';



class CommentListView extends StatefulWidget {
  CommentListView(this.ext) : super();

  final ContentExt ext;
  @override
  _CommentListState createState() => _CommentListState();
}

class _CommentListState extends State<CommentListView> with AutomaticKeepAliveClientMixin {

  final ActionClient actionClient = Get.find();

  late final req = CommentListReq(type:widget.ext.type,refId: widget.ext.refId,pageNo:1,pageSize:10);
  var times = 0;
  var list = List<Comment>.empty(growable: true);


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
    Get.log(req.toString());
    var response = await actionClient.stub.commentList(req);
    if (response.list.isEmpty) return;
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    globalState.userState.appendUsers(response.users);
    list.addAll(response.list);
    times++;
    req.pageNo++;
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    Get.log("CommentList重绘");
    final _future = grpcGetList();
    return FutureBuilder<void>(
        future: _future,
        builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
          return snapshot.handle() ??  RefreshIndicator(
              onRefresh: () {
                return resetList();
              },
              child: list.isEmpty? Center(child:TextButton.icon(onPressed: () { setState(() {}); },
              label: const Text('暂无评论，点击刷新'),
              icon: const Icon(Icons.refresh),)) : ListView.separated(
                  physics: BouncingScrollPhysics(),
                  controller: _controller,
                  itemCount: list.length,
                  separatorBuilder: (BuildContext context, int index) {
                    return Divider();
                  },
                  itemBuilder: (context, index) {
                    return InkWell(
                      onTap: (){
                        Get.to(()=>MomentDetailView(),arguments: list[index]);
                      },
                      child: CommentItem(
                          comment: list[index]),
                    );
                  }));
        });
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  void didUpdateWidget(covariant CommentListView oldWidget) {
    super.didUpdateWidget(oldWidget);
    times = 0;
    req.pageNo = 1;
    list.clear();
  }

  @override
  void reassemble() {
    print('调用了');

    super.reassemble();
  }

  @override
  bool get wantKeepAlive => true;
}