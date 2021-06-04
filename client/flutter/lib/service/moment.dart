import 'dart:io';

import 'package:app/model/moment.dart';
import 'package:app/util/dio.dart';
import 'package:json_annotation/json_annotation.dart';

part 'moment.g.dart';

@JsonSerializable()
class MomentListResponse {
  MomentListResponse();
  int code;
  int count;
  List<Moment> data;
  String msg;

  factory MomentListResponse.fromJson(Map<String, dynamic> json) => _$MomentListResponseFromJson(json);

  Map<String, dynamic> toJson() => _$MomentListResponseToJson(this);
}

Future<MomentListResponse> getMomentList(int pageNo,pageSize) async{
  var api = '/content?page=$pageNo&pageSize=$pageSize';

  try {
    var response = await httpClient().get(api);
    if (response.statusCode == HttpStatus.ok) {
      return MomentListResponse.fromJson(response.data);
    }
  } catch (exception) {

  }
  return null;
}