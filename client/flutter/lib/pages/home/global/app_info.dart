
import 'package:hive/hive.dart';

class AppInfo {

  static const _PRE = "AppInfo";
  // 版本
  static const StringVersionKey = _PRE+"VersionKey";
  // 打开次数
  static const IntOpenTimesKey = _PRE+"OpenTimesKey";

  init(Box box){
    final openTimes = box.get(IntOpenTimesKey,defaultValue:0);
    box.put(IntOpenTimesKey, openTimes+1);
  }
}
