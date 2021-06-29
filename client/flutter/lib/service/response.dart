import 'dart:io';
import 'package:dio/dio.dart';
import 'package:json_annotation/json_annotation.dart';

part 'response.g.dart';

@JsonSerializable()
class ResponseData$ {
  ResponseData$();
  late int code;
  late Map<String, dynamic> details;
  late String message;

  factory ResponseData$.fromJson(Map<String, dynamic> json) => _$ResponseDataFromJson(json);

  Map<String, dynamic> toJson() => _$ResponseDataToJson(this);
}

extension StringExtension on Response {
  Map<String, dynamic> getData() {
    if (this.statusCode == HttpStatus.ok) {
      final data = ResponseData$.fromJson(this.data);
      return data.details;
    }
    return Map() ;
  }
}