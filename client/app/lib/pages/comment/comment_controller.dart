import 'package:app/generated/protobuf/content/action.model.pb.dart';
import 'package:app/generated/protobuf/content/action.service.pb.dart';
import 'package:app/generated/protobuf/content/content.model.pbenum.dart';
import 'package:app/global/state.dart';
import 'package:app/rpc/action.dart';
import 'package:app/components/media/media.dart';
import 'package:get/get.dart';
import 'package:flutter/material.dart';
import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:grpc/grpc.dart';

import 'package:app/utils/dialog.dart';

class CommentController extends GetxController with MediaController {
  final req = CommentListReq(pageNo: 1, pageSize: 10);
  var times = 0;
  var list = List<Comment>.empty(growable: true);
  late Future<void> future;
  final ActionClient actionClient = Get.find();

  void init(ContentType type, $fixnum.Int64 refId) {
    req.type = type;
    req.refId = refId;
    future = grpcGetList();
  }

  Future<void> grpcGetList() async {
    globalService.logger.d(req.toString());
    try {
      var response = await actionClient.stub.commentList(req);
      if (response.list.isEmpty) return;
      // If the widget was removed from the tree while the message was in flight,
      // we want to discard the reply rather than calling setState to update our
      // non-existent appearance.
      globalState.userState.appendUsers(response.users);
      list.addAll(response.list);
      times++;
      req.pageNo++;
      update(["list"]);
    } catch (e) {
      final grpce = e as GrpcError;
      toast(grpce.message!);
    }
  }

  Future<void> resetList() {
    list.removeRange(0, list.length);
    req.pageNo = 1;
    future = grpcGetList();
    return future;
  }

  final TextEditingController textEditingController = TextEditingController();
  final focusNode = FocusNode();
  $fixnum.Int64? refId;
  $fixnum.Int64? replyId;
  $fixnum.Int64? rootId;
  String? image;
  $fixnum.Int64? recvId;
  ContentType? type;

  Future<void> save(String content) async {
    try {
      final object = await actionClient.stub.comment(CommentReq(
        type: type,
        content: content,
        refId: refId,
        replyId: replyId,
        rootId: rootId,
        recvId: recvId,
        image: image,
      ));
      list.add(Comment(
        id: object.id,
        content: content,
        type: type,
        refId: refId,
        replyId: replyId,
        rootId: rootId,
        recvId: recvId,
        image: image,
        userId: globalState.authState.userAuth!.id,
        user: globalState.authState.userBaseInfo,
      ));
      update(["list"]);
    } catch (e) {
      print(e);
    }
  }

  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }
}
