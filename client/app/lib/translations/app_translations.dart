import 'package:get/get.dart';

import 'en_US/en_us_translations.dart';
import 'vi_VN/vi_vn_translations.dart';
import 'zh_CN/zh_cn_translations.dart';

class AppTranslation extends Translations {
  @override
  Map<String, Map<String, String>> get keys => {
    'zh_CN': zhCn,
    'en_US': enUs,
    'vi_VN': viVn,
  };
}
