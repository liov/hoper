
// md5 加密
import 'dart:convert';
import 'package:crypto/crypto.dart';
import 'dart:io';

Future<String> calculateFileMD5(String filePath) async {
  File file = File(filePath);
  List<int> bytes = await file.readAsBytes();
  Digest digest = md5.convert(bytes);
  return digest.toString();
}

String MD5(String data) {
  // 将字符串转换为 UTF-8 编码的字节数组
  List<int> bytes = utf8.encode(data);
  // 生成 MD5 摘要
  Digest digest = md5.convert(bytes);
  // 转换为十六进制字符串（小写）
  return digest.toString();
}