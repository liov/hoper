import 'package:app/components/feed/feed_images.dart';
import 'package:app/global/const.dart';
import 'package:app/global/state.dart';
import 'package:app/pages/image/slide_image.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';

import 'package:app/pages/action_bar/action_bar.dart';
import 'package:app/util/time.dart';
import 'package:app/gen/pb/content/moment.model.pb.dart';

class MomentItem extends StatelessWidget {
  const MomentItem({super.key, required this.moment});

  final Moment moment;

  List<String> get _images =>
      moment.images.map((url) => BASE_STATIC_URL + url).toList(growable: false);

  @override
  Widget build(BuildContext context) {
    final user = globalState.userState.getUser(moment.userId);
    if (user == null) {
      return const SizedBox.shrink();
    }
    final created = getDateTime(
      moment.modelTime.createdAt.seconds.toInt(),
      moment.modelTime.createdAt.nanos.toInt(),
    );
    return RepaintBoundary(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Padding(
                padding: const EdgeInsets.all(10),
                child: GestureDetector(
                  onTap: () => slideImageRoute(BASE_STATIC_URL + user.avatar),
                  child: CircleAvatar(
                    radius: 22,
                    child: ClipOval(
                      child: ExtendedImage.network(
                        BASE_STATIC_URL + user.avatar,
                        width: 44,
                        height: 44,
                        fit: BoxFit.cover,
                        cache: true,
                        enableLoadState: false,
                      ),
                    ),
                  ),
                ),
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(user.name, maxLines: 1, overflow: TextOverflow.ellipsis),
                    Text('$created', style: Theme.of(context).textTheme.bodySmall),
                  ],
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(10),
                child: OutlinedButton(
                  style: OutlinedButton.styleFrom(
                    minimumSize: const Size(64, 32),
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                  ),
                  onPressed: () {},
                  child: const Text('+关注'),
                ),
              ),
            ],
          ),
          if (moment.content.isNotEmpty)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
              child: Text(
                moment.content,
                maxLines: 12,
                overflow: TextOverflow.ellipsis,
              ),
            ),
          if (_images.isNotEmpty)
            FeedImageGrid(urls: _images, onTap: slideImageRoute),
          ActionBar(moment),
        ],
      ),
    );
  }
}
