
import 'package:app/global/state.dart';
import 'package:app/pages/moment/list/moment_list_view.dart';
import 'package:app/pages/moment/physics.dart';
import 'package:app/pages/webview/webview.dart';
import 'package:app/pages/route.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../user/login_view.dart';
import 'add/moment_add_view.dart';
import 'list_v2/moment_list_v2_view.dart';
import 'moment_controller.dart';



class MomentView extends StatefulWidget{
  const MomentView({super.key});

  @override
  State<StatefulWidget> createState() => _MomentState();

}

class _MomentState extends State<MomentView> with AutomaticKeepAliveClientMixin {
  final MomentController controller = Get.find();


  @override
  Widget build(BuildContext context) {
    globalService.logger.d("MomentView重绘");
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
        actions: [
          IconButton(icon: const Icon(Icons.add),onPressed: (){
            if (globalState.authState.userAuth!=null) {
              Get.toNamed(Routes.MOMENT_ADD);
            } else {
              Get.to(()=>LoginView());
            }
          },)
        ],
      ),
      body: TabBarView(
        //physics:PageViewTabClampingScrollPhysics(controller:controller.homeController.to),
        controller: controller.tabController,
        children: controller.tabValues.map((f) {
          globalService.logger.d(f);
          if (f == "推荐") return Text('TODO');
          if (f == "刚刚") return MomentListV2View(tag:'newest');
          if (f == "关注") return Text('TODO');
          return Text(f);
        }).toList(),
      ),
    );
  }

  @override
  bool get wantKeepAlive => true;
}
