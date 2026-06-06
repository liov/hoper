import 'dart:io';

Future<bool> hasIPv6Support(String domain) async {
  try {
    // 指定查询IPv6地址类型
    final addresses = await InternetAddress.lookup(
      domain,
      type: InternetAddressType.IPv6,
    );
    return addresses.isNotEmpty; // 存在AAAA记录返回true
  } catch (e) {
    return false; // 无解析结果或异常则不支持
  }
}
