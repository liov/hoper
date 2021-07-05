import 'package:app/generated/protobuf/content/content.model.pb.dart';
import 'package:app/model/moment/moment.dart';
import 'package:get/get.dart';

class MomentListController extends GetxController {
  var pageNo = 1;
  var pageSize = 10;
  var times = 0;
  var list = List<Moment>.empty(growable: true).obs;
  var list$ = List<Moment$>.empty(growable: true).obs;


  timesIncrement() => times++;
  pageNoIncrement() => pageNo++;
  resetList(){
    list.removeRange(0, list.length);
    pageNo = 1;
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
