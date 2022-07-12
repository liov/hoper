import 'package:flutter/material.dart';

import 'package:get/get.dart';

import '../../global/global_service.dart';
import 'home_controller.dart';

class DashboardView extends GetView<HomeController> {
  @override
  Widget build(BuildContext context) {
    globalService.logger.d('DashboardView重绘');
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              'DashboardView is working',
              style: TextStyle(fontSize: 20),
            ),
            Obx(() => Text('Time: ${controller.now.toString()}',style: TextStyle(fontSize: 20),)),
          ],
        ),
      ),
    );
  }
}
