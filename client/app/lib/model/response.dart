import 'dart:io';
import 'package:app/global/service.dart';
import 'package:app/model/weibo/weibo.dart';
import 'package:dio/dio.dart' as $dio;
import 'package:get/get.dart';
import 'package:json_annotation/json_annotation.dart';

part 'response.g.dart';

@JsonSerializable(genericArgumentFactories: true)
class ResponseData<T> {
  ResponseData({required this.code, this.data, this.msg});
  final  int code;
  final  T? data;
  final String? msg;

  factory ResponseData.fromJson(Map<String, dynamic> json,T Function(dynamic json) fromJsonT) => _$ResponseDataFromJson(json,fromJsonT);

  Map<String, dynamic> toJson(T Function(T) toJsonT) => _$ResponseDataToJson(this,toJsonT);
}

extension Extension<T> on $dio.Response {
  T? getData(T Function(Object? json) fromJsonT) {
    if (statusCode == HttpStatus.ok) {
      final data = ResponseData.fromJson(this.data,fromJsonT);
      if (data.code!=0) {
        globalService.logger.e('请求出错 ${data.msg}');
        Get.rawSnackbar(message: data.msg!);
      }
        return data.data;
    }
    globalService.logger.e('请求出错 $statusCode');
    Get.rawSnackbar(message:'请求出错');
    return null;
  }

  T? getBaoyuData(T Function(Object? json) fromJsonT) {
    if (statusCode == HttpStatus.ok) {
      final data = ResponseData.fromJson(this.data,fromJsonT);
      if (data.code!=200) {
        globalService.logger.e('请求出错 ${data.msg}');
        Get.rawSnackbar(message:data.msg!);
      };
      return data.data;
    }
    globalService.logger.e('请求出错 $statusCode');
    Get.rawSnackbar(message:'请求出错');
    return null;
  }

  T? getWeiboData(T Function(Object? json) fromJsonT) {
    if (statusCode == HttpStatus.ok) {
      final data = WeiboResponse.fromJson(this.data as Map<String, dynamic>,fromJsonT);
      globalService.logger.d(data);
      if (data.ok!=1) {
        globalService.logger.e('请求出错 ${data.msg}');
        Get.rawSnackbar(message: data.msg!);
      }
      return data.data;
    }
    globalService.logger.e('请求出错 $statusCode');
    Get.rawSnackbar(message:'请求出错');
    return null;
  }
}
