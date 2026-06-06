
import 'package:app/providers/providers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class CommentAdd extends ConsumerWidget {
  const CommentAdd({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.watch(commentsProvider);
    final comments = ref.read(commentsProvider.notifier);
    final controller0 = comments.textEditingController;
    final focusNode = comments.focusNode;
    return Row(mainAxisSize: MainAxisSize.min, children: [
      const Padding(
        padding: EdgeInsets.symmetric(horizontal: 10),
      ),
      Expanded(
          flex: 7,
          child: TextField(
            controller: controller0,
            focusNode: focusNode,
            keyboardType: TextInputType.multiline,
            maxLines: null,
            maxLength: 512,
            decoration: InputDecoration(
              counterText: '',
              hintText: '评论',
              border: OutlineInputBorder(
                borderSide: const BorderSide(width: 0, color: Colors.transparent),
                borderRadius: const BorderRadius.all(Radius.circular(10)),
              ),
              contentPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 10),
              suffixIcon: IconButton(
                  onPressed: () async {
                    focusNode.unfocus();
                    showBottomSheet(
                      context: context,
                      builder: (BuildContext context) {
                        return const Text('测试');
                      },
                    );
                  },
                  color: Colors.blue,
                  icon: const Icon(Icons.mood),
                  tooltip: '发送'),
            ),
            onChanged: comments.onDraftChanged,
          )),
      Expanded(
        flex: 1,
        child: _SendButton(comments: comments),
      )
    ]);
  }
}

class _SendButton extends StatelessWidget {
  const _SendButton({required this.comments});
  final Comments comments;

  @override
  Widget build(BuildContext context) {
    final controller0 = comments.textEditingController;
    final focusNode = comments.focusNode;
    return IconButton(
        onPressed: () async {
          if (controller0.text.isEmpty) {
            controller0.text = "为了遇见你我珍惜自己我穿越风和雨是为交出我的心，直到遇见你";
            return;
          }
          await comments.save(controller0.text);
          controller0.text = '';
          focusNode.unfocus();
          comments.onDraftCleared();
        },
        color: Colors.blue,
        icon: controller0.text.isEmpty ? const Icon(Icons.add) : const Text("发送"),
        tooltip: '更多');
  }
}
