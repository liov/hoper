
import 'package:app/components/async/async.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'dart:io';

import '../../global/service.dart';
import 'controller.dart';

class ImageView extends StatelessWidget {
  ImageView({super.key});

  final WeiboController controller = Get.find();

  late final ScrollController _controller = ScrollController()
    ..addListener(() {
      if (_controller.position.atEdge) {
        globalService.logger.d('这里执行了吗');
        controller.getList();
      }
    }
    );
  @override
  Widget build(BuildContext context) {
    globalService.logger.d('ImageView重绘');
    final width = MediaQuery
        .of(context)
        .size
        .width;
    final height = MediaQuery
        .of(context)
        .size
        .height;
    final columns = (width / controller.picWidth).floor();
    final rows = (height / controller.picHeight).floor();
    globalService.logger.d('$columns $rows');
    final  future = controller.getList();
    return FutureBuilder<void>(
        future: future,
        builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
          return snapshot.handle() ?? GetBuilder<WeiboController>(
              builder: (_)=>RefreshIndicator(
              onRefresh: () {
                globalService.logger.d('这里执行了吗2');
                return controller.getList();
              },
              child:GridView.builder(
                gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
                  crossAxisCount: 5,
                  mainAxisSpacing: 10.0,
                  crossAxisSpacing: 10.0,
                ),
                itemCount: controller.list.length, // 数据源的长度
                itemBuilder: (BuildContext context, int index) {
                  return ExtendedImage.network(
                    controller.list[index],
                    width: controller.picWidth.toDouble(),
                    height: controller.picHeight.toDouble(),
                    fit: BoxFit.scaleDown
                  );
                },
              )
              )
          );
        });
  }
}