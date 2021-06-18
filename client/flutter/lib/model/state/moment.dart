import 'package:get/get.dart';

import '../moment.dart';

class MomentState extends GetxController {
  var pageNo = 1.obs;
  var pageSize = 10.obs;
  var times = 0.obs;
  var list = List<Moment>.empty(growable: true).obs;
  timesIncrement() => times++;
  pageNoIncrement() => pageNo++;
  reset(){
    list.removeRange(0, list.length);
    pageNo.value = 1;
  }
}
