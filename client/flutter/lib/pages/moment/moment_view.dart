import 'package:app/pages/dynamic/dynamic.dart';
import 'package:app/pages/moment/physics.dart';
import 'package:app/pages/webview/webview.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'list_v2/moment_list_v2_view.dart';
import 'moment_controller.dart';



class MomentView extends StatefulWidget{
  @override
  State<StatefulWidget> createState() => _MomentState();

}

class _MomentState extends State<MomentView> with AutomaticKeepAliveClientMixin {
  final MomentController controller = Get.find();
  @override
  Widget build(BuildContext context) {
    super.build(context);
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: TabBar(
          isScrollable: true,
          tabs: controller.tabValues.map((choice) {
            return Tab(text: choice,);
          }).toList(),
          controller: controller.tabController,
        ),
      ),
      body: TabBarView(
        physics:PageViewTabClampingScrollPhysics(controller:controller.homeController.to),
        controller: controller.tabController,
        children: controller.tabValues.map((f) {
          if (f == "推荐") return MomentListV2View();
          if (f == "刚刚") return MomentListV2View(tag:'newest');
          if (f == "关注") return MomentListV2View(tag:'follow');
          return Container(child:Text(f));
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

  @override
  bool get wantKeepAlive => true;
}
