import 'dart:ui';

import 'package:applib/util/async.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';

import 'package:app/global/service.dart';
import 'package:app/pages/image/slide_image.dart';
import 'controller.dart';

class ImageView extends StatelessWidget {
  ImageView({super.key}) : super();

  final WeiboController controller = Get.find();
  late final Future<void> future = controller.getList();
  late final ScrollController _controller = ScrollController();
  final TextEditingController _searchController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    globalService.logger.fine('ImageView重绘');
    final size = MediaQuery.sizeOf(context);
    final width = size.width;
    final height = size.height;
    final columns = (width / controller.picWidth).floor();
    final rows = (height / controller.picHeight).floor();
    globalService.logger.fine('$columns $rows');

    return PopScope(
        canPop:false,
      onPopInvokedWithResult: (didPop, result)  {
        if (didPop) {
          // 页面已经被 pop 了，不需要额外处理
          return;
        }
        if (Theme.of(context).platform == TargetPlatform.android) {
          // Android：尝试让 App 进入后台
          SystemNavigator.pop(); // 进入后台
        } else {}
        globalService.logger.fine(result);
      },
      child: Scaffold(
        appBar: AppBar(
          title: TextField(
            controller: _searchController,
            decoration: InputDecoration(
              hintText: '复制分享链接到此处',
              border: InputBorder.none,
              hintStyle: const TextStyle(color: Colors.white70),
            ),
            style: const TextStyle(color: Colors.white),
            autofocus: false,
          ),
          actions: [
            IconButton(
              icon: const Icon(Icons.search),
              onPressed: () {
                final idStr =
                    _searchController.text.split('/').last.split('#').first;
                globalService.logger.fine(idStr);
                controller.newList(int.parse(_searchController.text));
              },
            ),
          ],
        ),
        body: FutureBuilder<void>(
          future: future,
          builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
            return snapshot.handle() ??
                GetBuilder<WeiboController>(
                  builder: (_) => refreshIndicator(context),
                  dispose: (state) => _controller.dispose(),
                );
          },
        ),
      ),
    );
  }

  Widget refreshIndicator(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () {
        globalService.logger.fine('这里执行了吗2');
        return controller.getList();
      },
      child: ScrollConfiguration(
        behavior: ScrollConfiguration.of(context).copyWith(
          dragDevices: {PointerDeviceKind.touch, PointerDeviceKind.mouse},
        ),
        child: NotificationListener<ScrollEndNotification>(
          onNotification: (n) {
            if (n.metrics.pixels >= n.metrics.maxScrollExtent) {
              controller.getList();
            }
            return false;
          },
          child: GridView.builder(
          gridDelegate: SliverGridDelegateWithMaxCrossAxisExtent(
            maxCrossAxisExtent: 300,
            mainAxisSpacing: 10.0,
            crossAxisSpacing: 10.0,
          ),
          itemCount: controller.list.length,
          // 数据源的长度
          itemBuilder: (BuildContext context, int index) {
            return GestureDetector(
              child: Hero(
                tag: controller.list[index],
                child: ExtendedImage.network(
                  controller.list[index],
                  width: controller.picWidth.toDouble(),
                  height: controller.picHeight.toDouble(),
                  fit: BoxFit.scaleDown,
                ),
              ),
              onTap: () {
                slideImageRoute(controller.list[index]);
              },
            );
          },
          controller: _controller,
          physics: AlwaysScrollableScrollPhysics(),
        ),
        ),
      ),
    );
  }
}
