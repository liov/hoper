
import 'package:get/get.dart';
import 'package:applib/util/multi_entity.dart';

import '../../global/service.dart';
import '../../model/weibo/weibo.dart';
import '../../rpc/weibo.dart';

class WeiboController extends GetxController{
  final WeiboClient weiboClient = Get.find();
  int userId = 0;
  int page = 1;
  int feature = 0;
  String sinceId = '';
  List<String> list = [];
  int picWidth = 300;
  int picHeight = 300;
  Future<void> getList() async {
    globalService.logger.d('getList');
    try{
    final response = await weiboClient.getList(uid: userId, page: page, feature: feature, sinceId: sinceId);

    for (var (e as Map<String,dynamic>)in response?['list'] as List<dynamic>) {
      if (e['pic_infos'] != null){
        list.addAll((e['pic_infos']as Map<String,dynamic>).values.map((v){
          final picInfo = (v as Map<String,dynamic>)['mw2000'] as Map<String,dynamic>;
          return picInfo['url'] as String;
        }
        ));
      }
    }
    globalService.logger.d(list);
    update();
    }catch(e){
      globalService.logger.e(e);
    }
  }


  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }
}
