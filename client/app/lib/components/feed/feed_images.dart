import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';

class FeedImageGrid extends StatelessWidget {
  const FeedImageGrid({super.key, required this.urls, this.onTap});

  final List<String> urls;
  final void Function(String url)? onTap;

  @override
  Widget build(BuildContext context) {
    if (urls.isEmpty) {
      return const SizedBox.shrink();
    }
    final width = MediaQuery.sizeOf(context).width;
    final side = ((width - 32 - 8) / 3).floorToDouble();
    final show = urls.length > 9 ? urls.sublist(0, 9) : urls;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      child: Wrap(
        spacing: 4,
        runSpacing: 4,
        children: [
          for (final url in show)
            RepaintBoundary(
              child: GestureDetector(
                onTap: onTap == null ? null : () => onTap!(url),
                child: ExtendedImage.network(
                  url,
                  width: side,
                  height: side,
                  fit: BoxFit.cover,
                  cache: true,
                  enableLoadState: false,
                  gaplessPlayback: true,
                ),
              ),
            ),
        ],
      ),
    );
  }
}
