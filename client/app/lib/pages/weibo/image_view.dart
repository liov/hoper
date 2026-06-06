import 'dart:ui';

import 'package:app/util/async.dart';
import 'package:app/providers/providers.dart';
import 'package:app/providers/weibo_notifier.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:app/global/service.dart';
import 'package:app/pages/image/slide_image.dart';

class ImageView extends ConsumerWidget {
  const ImageView({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return const _ImageViewBody();
  }
}

class _ImageViewBody extends ConsumerStatefulWidget {
  const _ImageViewBody();

  @override
  ConsumerState<_ImageViewBody> createState() => _ImageViewBodyState();
}

class _ImageViewBodyState extends ConsumerState<_ImageViewBody> {
  Future<void>? _future;
  ScrollController? _controller;
  late final TextEditingController _searchController;
  var _initialized = false;

  @override
  void initState() {
    super.initState();
    _searchController = TextEditingController();
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (_initialized) {
      return;
    }
    _initialized = true;
    final weibo = ref.read(weiboProvider.notifier);
    _future = weibo.getList();
    _controller = ScrollController();
  }

  @override
  void dispose() {
    _controller?.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final weiboState = ref.watch(weiboProvider);
    final weibo = ref.read(weiboProvider.notifier);
    globalService.logger.fine('ImageView重绘');
    final size = MediaQuery.sizeOf(context);
    final width = size.width;
    final height = size.height;
    final columns = (width / weiboState.picWidth).floor();
    final rows = (height / weiboState.picHeight).floor();
    globalService.logger.fine('$columns $rows');

    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, result) {
        if (didPop) {
          return;
        }
        if (Theme.of(context).platform == TargetPlatform.android) {
          SystemNavigator.pop();
        }
        globalService.logger.fine(result);
      },
      child: Scaffold(
        appBar: AppBar(
          title: TextField(
            controller: _searchController,
            decoration: const InputDecoration(
              hintText: '复制分享链接到此处',
              border: InputBorder.none,
              hintStyle: TextStyle(color: Colors.white70),
            ),
            style: const TextStyle(color: Colors.white),
            autofocus: false,
          ),
          actions: [
            IconButton(
              icon: const Icon(Icons.search),
              onPressed: () {
                final idStr = _searchController.text.split('/').last.split('#').first;
                globalService.logger.fine(idStr);
                weibo.newList(int.parse(_searchController.text));
              },
            ),
          ],
        ),
        body: FutureBuilder<void>(
          future: _future,
          builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
            return snapshot.handle() ?? _refreshIndicator(context, weiboState, weibo);
          },
        ),
      ),
    );
  }

  Widget _refreshIndicator(BuildContext context, WeiboState weiboState, Weibo weibo) {
    return RefreshIndicator(
      onRefresh: () {
        globalService.logger.fine('这里执行了吗2');
        return weibo.getList();
      },
      child: ScrollConfiguration(
        behavior: ScrollConfiguration.of(context).copyWith(
          dragDevices: {PointerDeviceKind.touch, PointerDeviceKind.mouse},
        ),
        child: NotificationListener<ScrollEndNotification>(
          onNotification: (n) {
            if (n.metrics.pixels >= n.metrics.maxScrollExtent) {
              weibo.getList();
            }
            return false;
          },
          child: GridView.builder(
            gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
              maxCrossAxisExtent: 300,
              mainAxisSpacing: 10.0,
              crossAxisSpacing: 10.0,
            ),
            itemCount: weiboState.list.length,
            itemBuilder: (BuildContext context, int index) {
              return GestureDetector(
                child: Hero(
                  tag: weiboState.list[index],
                  child: ExtendedImage.network(
                    weiboState.list[index],
                    width: weiboState.picWidth.toDouble(),
                    height: weiboState.picHeight.toDouble(),
                    fit: BoxFit.scaleDown,
                  ),
                ),
                onTap: () {
                  slideImageRoute(weiboState.list[index]);
                },
              );
            },
            controller: _controller,
            physics: const AlwaysScrollableScrollPhysics(),
          ),
        ),
      ),
    );
  }
}
