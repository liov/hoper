import 'package:app/pages/webview/webview.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'moment_controller.dart';
import 'list/moment_list_view.dart';
import 'moment_state.dart';

class MomentView extends GetView<MomentController> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: TabBar(
          isScrollable: true,
          tabs: controller.state.tabValues.map((choice) {
            return Tab(
              text: choice,
            );
          }).toList(),
          controller: controller.ac,
        ),
      ),
      body: TabBarView(
        controller: controller.ac,
        children: controller.state.tabValues.map((f) {
          if (f == "推荐") return MomentListView();
          return Center(child: Text(f),);
        }).toList(),
      ),
      floatingActionButton: FloatingActionButton(
        heroTag: 'login',
        onPressed: () => Get.to(WebViewExample()),
        tooltip: 'ToBrowser',
        child: Icon(Icons.send),
      ),
    );
  }
}
