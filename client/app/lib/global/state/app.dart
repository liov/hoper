import 'package:flutter/foundation.dart';

import 'package:app/global/service.dart';

class AppState {
  static bool isDebug = kDebugMode;

  static const _PRE = "AppInfo";

  // 版本
  static const StringVersionKey = "${_PRE}VersionKey";

  // 打开次数
  static const IntOpenTimesKey = "${_PRE}OpenTimesKey";

  void init() {
    final openTimes = globalService.box.get(IntOpenTimesKey, defaultValue: 0);
    globalService.box.put(IntOpenTimesKey, openTimes + 1);
  }
}
