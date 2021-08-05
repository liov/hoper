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
      if (data.code!=0) Get.rawSnackbar(message:data.message!);
      return data.details ?? Map();
    }
    Get.rawSnackbar(message:'请求出错');
    return Map() ;
  }
}