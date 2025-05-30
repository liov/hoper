
import 'package:app/model/response.dart';
import 'package:app/global/dio.dart';
import 'package:dio/dio.dart';
import 'package:intl/intl.dart';

class BaoyuClient  {
  static Future<String> signup(String prefix) async {
    final date = DateFormat('yyyyMMdd').format(DateTime.now());
    final phone = prefix + date;
    try {
      final response = await httpClient.post("https://api.xiaoyoucaip2p.com/api/appsign",
        data: {"password":"zxcvbnm1","mobile":phone,"source":"ios","pid":"16530915"},
      options: Options(headers:{"Cookie":"laravel_session=KN4BBcHUzxw79Quey7yY2DNtgg7zzuNhmfzSTF8V","user-agent":"96bao yu/1.0 (iPhone; iOS 14.6; Scale/3.00)","accept-language":"zh-Hans-CN;q=1"}));
      final data = response.getBaoyuData((v) => v as String);
      print(response.data);
      return data["access_token"];
    } catch (e) {
      print(e);
      return '';
    }
  }
}