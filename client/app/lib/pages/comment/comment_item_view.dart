import 'package:app/components/feed/feed_images.dart';
import 'package:app/gen/pb/content/action.model.pb.dart';
import 'package:app/global/const.dart';
import 'package:app/global/state.dart';
import 'package:app/pages/image/slide_image.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';

class CommentItem extends StatelessWidget {
  const CommentItem({super.key, required this.comment});

  final Comment comment;

  List<String> get _images {
    if (comment.image.isEmpty) return const [];
    return comment.image.split(',').map((url) => BASE_STATIC_URL + url).toList(growable: false);
  }

  @override
  Widget build(BuildContext context) {
    final user = globalState.userState.getUser(comment.userId);
    if (user == null) {
      return const SizedBox.shrink();
    }
    return RepaintBoundary(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Padding(
                padding: const EdgeInsets.all(10),
                child: CircleAvatar(
                  radius: 20,
                  child: ClipOval(
                    child: ExtendedImage.network(
                      BASE_STATIC_URL + user.avatar,
                      width: 40,
                      height: 40,
                      fit: BoxFit.cover,
                      cache: true,
                      enableLoadState: false,
                    ),
                  ),
                ),
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(user.name),
                    Text('${comment.modelTime.createdAt}', style: Theme.of(context).textTheme.bodySmall),
                  ],
                ),
              ),
              IconButton(icon: const Icon(Icons.favorite), onPressed: () {}),
            ],
          ),
          if (comment.content.isNotEmpty)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              child: Text(comment.content, maxLines: 8, overflow: TextOverflow.ellipsis),
            ),
          if (_images.isNotEmpty) FeedImageGrid(urls: _images, onTap: slideImageRoute),
        ],
      ),
    );
  }
}
