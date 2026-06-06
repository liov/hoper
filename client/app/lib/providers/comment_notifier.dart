import 'package:app/gen/pb/content/action.model.pb.dart';
import 'package:app/gen/pb/content/action.service.pb.dart';
import 'package:app/gen/pb/content/content.model.pbenum.dart';
import 'package:app/global/state.dart';
import 'package:app/util/dialog.dart';
import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:flutter/material.dart';
import 'package:grpc/grpc.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'comment_notifier.g.dart';

class CommentState {
  const CommentState({
    this.list = const [],
    this.draftEmpty = true,
    this.revision = 0,
    this.loading = false,
  });

  final List<Comment> list;
  final bool draftEmpty;
  final int revision;
  final bool loading;

  CommentState copyWith({List<Comment>? list, bool? draftEmpty, int? revision, bool? loading}) {
    return CommentState(
      list: list ?? this.list,
      draftEmpty: draftEmpty ?? this.draftEmpty,
      revision: revision ?? this.revision,
      loading: loading ?? this.loading,
    );
  }
}

@riverpod
class Comments extends _$Comments {
  final req = CommentListReq(pageNo: 1, pageSize: 10);
  var times = 0;
  late Future<void> future;

  final textEditingController = TextEditingController();
  final focusNode = FocusNode();

  $fixnum.Int64? refId;
  $fixnum.Int64? replyId;
  $fixnum.Int64? rootId;
  String? image;
  $fixnum.Int64? recvId;
  ContentType? type;

  @override
  CommentState build() {
    future = Future.value();
    ref.onDispose(() {
      textEditingController.dispose();
      focusNode.dispose();
    });
    return const CommentState();
  }

  void setup(ContentType contentType, $fixnum.Int64 contentRefId) {
    type = contentType;
    refId = contentRefId;
    req.type = contentType;
    req.refId = contentRefId;
    future = grpcGetList();
  }

  Future<void> grpcGetList() async {
    globalService.logger.fine(req.toString());
    try {
      final response = await globalService.actionClient.stub.commentList(req);
      if (response.list.isEmpty) return;
      globalState.userState.appendUsers(response.users);
      final next = [...state.list, ...response.list];
      times++;
      req.pageNo++;
      state = state.copyWith(list: next, revision: state.revision + 1);
    } catch (e) {
      final grpce = e as GrpcError;
      toast(grpce.message!);
    }
  }

  Future<void> resetList() {
    req.pageNo = 1;
    times = 0;
    state = state.copyWith(list: [], revision: state.revision + 1);
    future = grpcGetList();
    return future;
  }

  void onDraftChanged(String value) {
    final empty = value.isEmpty;
    if (state.draftEmpty == empty) return;
    state = state.copyWith(draftEmpty: empty, revision: state.revision + 1);
  }

  void onDraftCleared() {
    if (state.draftEmpty) return;
    state = state.copyWith(draftEmpty: true, revision: state.revision + 1);
  }

  Future<void> save(String content) async {
    try {
      final object = await globalService.actionClient.stub.comment(
        CommentReq(
          type: type,
          content: content,
          refId: refId,
          replyId: replyId,
          rootId: rootId,
          recvId: recvId,
          image: image,
        ),
      );
      final next = [
        ...state.list,
        Comment(
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
        ),
      ];
      state = state.copyWith(list: next, revision: state.revision + 1);
    } catch (e) {
      globalService.logger.warning('$e');
    }
  }
}
