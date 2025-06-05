import 'dart:ui';

import 'package:applib/util/async.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'dart:io';

import '../../global/service.dart';
import '../image/slide_image.dart';
import 'controller.dart';

class ImageView extends StatelessWidget {
  ImageView({super.key}) : super();

  final WeiboController controller = Get.find();
  late final Future<void> future = controller.getList();
  late final ScrollController _controller =
      ScrollController()..addListener(() {
        if (_controller.position.pixels >=
            _controller.position.maxScrollExtent) {
          globalService.logger.d('这里执行了吗');
          controller.getList();
        }
      });
  final TextEditingController _searchController = TextEditingController();
  @override
  Widget build(BuildContext context) {
    globalService.logger.d('ImageView重绘');
    final width = MediaQuery.of(context).size.width;
    final height = MediaQuery.of(context).size.height;
    final columns = (width / controller.picWidth).floor();
    final rows = (height / controller.picHeight).floor();
    globalService.logger.d('$columns $rows');

    return  Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _searchController,
          decoration: InputDecoration(
            hintText: '复制分享链接到此处',
            border: InputBorder.none,
            hintStyle: const TextStyle(color: Colors.white70),
          ),
          style: const TextStyle(color: Colors.white),
          autofocus: true,
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.close),
            onPressed: () {
              final idStr = _searchController.text.split('/').last.split('#').first;
              globalService.logger.d(idStr);
              controller.newList(int.parse(_searchController.text));
            },
          ),
        ],
      ),
      body: FutureBuilder<void>(
          future: future,
          builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
            return snapshot.handle() ?? GetBuilder<WeiboController>(
              builder: (_) => refreshIndicator(context),
              dispose: (state) => _controller.dispose(),
            );
          }
    ),
    );
  }

  Widget refreshIndicator(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () {
        globalService.logger.d('这里执行了吗2');
        return controller.getList();
      },
      child: ScrollConfiguration(
        behavior: ScrollConfiguration.of(context).copyWith(
          dragDevices: {PointerDeviceKind.touch, PointerDeviceKind.mouse},
        ),
        child: GridView.builder(
          gridDelegate: SliverGridDelegateWithMaxCrossAxisExtent(
            maxCrossAxisExtent: 300,
            mainAxisSpacing: 10.0,
            crossAxisSpacing: 10.0,
          ),
          itemCount: controller.list.length,
          // 数据源的长度
          itemBuilder: (BuildContext context, int index) {
            return GestureDetector(child:Hero(
            tag: controller.list[index],
                child:ExtendedImage.network(
              controller.list[index],
              width: controller.picWidth.toDouble(),
              height: controller.picHeight.toDouble(),
              fit: BoxFit.scaleDown,
            )),
              onTap: () {
                slideImageRoute(controller.list[index]);
              },
            );
          },
          controller: _controller,
          physics: AlwaysScrollableScrollPhysics(),
        ),
      ),
    );
  }
}
