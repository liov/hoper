
import 'package:get/get.dart';

import '../../global/service.dart';
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
  bool isEnd = false;


  Future<void> newList(int userId) async {
    this.userId = userId;
    page =1;
    isEnd = false;
    list.clear();
    return getList();
  }

  Future<void> getList() async {
    globalService.logger.fine('getList');
    if(isEnd) return;
    final response = await weiboClient.getOriginalList(uid: userId, page: page, feature: feature, sinceId: sinceId);
    if (response == null) {
      isEnd = true;
      return;
    }
    if (response.list.isEmpty) {
      isEnd = true;
      return;
    }
    //sinceId = response.sinceId;
    for (var e in response.list) {
      if (e.picInfos != null){
        list.addAll(e.picInfos!.values.map((v){
          return  v.mw2000.url;
        }
        ));
      }
    }
    globalService.logger.fine('${response.list.length} ${list.length}');
    globalService.logger.fine(list);
    page++;
    update();

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
