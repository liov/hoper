import 'en_US/en_us_translations.dart';
import 'vi_VN/vi_vn_translations.dart';
import 'zh_CN/zh_cn_translations.dart';

/// 文案表（原 GetX Translations，保留供后续接入 MaterialApp 本地化）。
abstract final class AppTranslation {
  static final Map<String, Map<String, String>> keys = {
    'zh_CN': zhCn,
    'en_US': enUs,
    'vi_VN': viVn,
  };
}
