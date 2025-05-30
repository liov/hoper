import 'package:app/generated/protobuf/content/content.model.pb.dart';
import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/global/state.dart';
import 'package:app/rpc/moment.dart';
import 'package:applib/util/multi_entity.dart';
import 'package:get/get.dart';

import 'package:app/generated/protobuf/content/moment.model.pb.dart';

// 相当于多个controller,实验性，不要这么用
class MomentListController extends GetxController with MultiEntity<ListState>{
  final MomentClient momentClient = Get.find();

  Future<void> newList(String tag) async{
    if (getEntity(tag)!=null) return;
    final list = ListState(tag);
    await list.grpcGetList(momentClient);
    entityMap[tag] = list;
  }

  Future<void> pullList(String tag) async{
    await entityMap[tag]?.grpcGetList(momentClient);
    update([tag]);
  }

  Future<void> resetList(String tag) async{
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

  grpcGetList(MomentClient momentClient) async {
    var response = await momentClient.stub.list(req);
    if (response.list.isEmpty) return;
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    globalState.userState.appendUsers(response.users);
    list.addAll(response.list);
    times++;
    req.pageNo++;
  }

}