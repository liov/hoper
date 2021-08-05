import 'package:app/generated/protobuf/content/content.model.pb.dart';
import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/global/global_controller.dart';
import 'package:app/service/moment.dart';
import 'package:app/utils/multi_entity.dart';
import 'package:get/get.dart';

// 相当于多个controller,实验性，不要这么用
class MomentListController extends GetxController with MultiEntity<ListState>{
  final MomentClient momentClient = Get.put(MomentClient());

  newList(String tag) async{
    if (getEntity(tag)!=null) return;
    final list = ListState(tag);
    await list.grpcGetList(momentClient, globalController);
    entityMap[tag] = list;
  }

  pullList(String tag) async{
    await entityMap[tag]?.grpcGetList(momentClient, globalController);
    update([tag]);
  }

  resetList(String tag) async{
    entityMap[tag]?.resetList();
    await pullList(tag);
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

class ListState {
  ListState(String tag){
    req = MomentListReq(pageNo:1,pageSize:10);
  }


  late final MomentListReq req;
  var times = 0;
  var list = List<Moment>.empty(growable: true);


  resetList(){
    list.removeRange(0, list.length);
    req.pageNo = 1;
  }

  grpcGetList(MomentClient momentClient,GlobalController globalController) async {
    var response = await momentClient.stub.list(req);
    if (response.list.isEmpty) return;
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    globalController.userState.appendUsers(response.users);
    list.addAll(response.list);
    times++;
    req.pageNo++;
  }

}