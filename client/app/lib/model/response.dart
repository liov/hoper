import 'dart:io';
import 'package:dio/dio.dart' as $dio;
import 'package:get/get.dart';
import 'package:json_annotation/json_annotation.dart';

part 'response.g.dart';

@JsonSerializable()
class ResponseData {
  ResponseData();
  late int code;
  late Map<String, dynamic>? details;
  String? message;

  factory ResponseData.fromJson(Map<String, dynamic> json) => _$ResponseDataFromJson(json);

  Map<String, dynamic> toJson() => _$ResponseDataToJson(this);
}

extension StringExtension on $dio.Response {
  Map<String, dynamic> getData() {
    if (this.statusCode == HttpStatus.ok) {
      final data = ResponseData.fromJson(this.data);
      if (data.code!=0) Get.rawSnackbar(message:data.msg!);
      return data.data ?? Map();
    }
    Get.rawSnackbar(message:'请求出错');
    return Map() ;
  }

  Map<String, dynamic> getData2() {
    if (this.statusCode == HttpStatus.ok) {
      final data = _$ResponseDataFromJson2(this.data);
      if (data.code!=200) Get.rawSnackbar(message:data.msg!);
      return data.data ?? Map();
    }
    Get.rawSnackbar(message:'请求出错');
    return Map() ;
  }
}

ResponseData _$ResponseDataFromJson2(Map<String, dynamic> json) {
  return ResponseData()
    ..code = json['code'] as int
    ..details = json['data'] as Map<String, dynamic>?
    ..message = json['msg'] as String;
}