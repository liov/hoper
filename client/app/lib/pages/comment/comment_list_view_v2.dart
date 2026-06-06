import 'package:app/util/async.dart';
import 'package:app/gen/pb/content/action.model.pb.dart';
import 'package:app/providers/providers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'comment_item_view.dart';

class CommentListViewV2 extends ConsumerStatefulWidget {
  const CommentListViewV2(this.ext, {super.key});

  final Statistics ext;

  @override
  ConsumerState<CommentListViewV2> createState() => _CommentListViewV2State();
}

class _CommentListViewV2State extends ConsumerState<CommentListViewV2> {
  ScrollController? _scrollController;
  var _initialized = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (_initialized) {
      return;
    }
    _initialized = true;
    final comments = ref.read(commentsProvider.notifier);
    comments.setup(widget.ext.type, widget.ext.refId);
    _scrollController = ScrollController()
      ..addListener(() {
        if (_scrollController!.position.atEdge) {
          comments.grpcGetList();
        }
      });
  }

  @override
  void dispose() {
    _scrollController?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final commentState = ref.watch(commentsProvider);
    final comments = ref.read(commentsProvider.notifier);
    return FutureBuilder<void>(
        future: comments.future,
        builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
          return snapshot.handle() ??
              RefreshIndicator(
                  onRefresh: () {
                    return comments.resetList();
                  },
                  child: ListView.separated(
                      physics: const BouncingScrollPhysics(),
                      controller: _scrollController,
                      itemCount: commentState.list.length,
                      separatorBuilder: (BuildContext context, int index) {
                        return const Divider();
                      },
                      itemBuilder: (context, index) {
                        return CommentItem(comment: commentState.list[index]);
                      }));
        });
  }
}
