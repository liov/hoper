
import 'package:app/global/service.dart';
import 'package:get/get.dart';


class AppInfo extends GetxController{

  AppInfo():assert(isDebug = true);

  static bool isDebug = false;

  static const _PRE = "AppInfo";
  // 版本
  static const StringVersionKey = _PRE+"VersionKey";
  // 打开次数
  static const IntOpenTimesKey = _PRE+"OpenTimesKey";

  init(){
    final openTimes = globalService.box.get(IntOpenTimesKey,defaultValue:0);
    globalService.box.put(IntOpenTimesKey, openTimes+1);
  }

}
