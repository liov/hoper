import 'package:app/model/moment.dart';
import 'package:json_annotation/json_annotation.dart';

part 'moment.g.dart';

@JsonSerializable(nullable: false)
class MomentListResponse {
  MomentListResponse();
  int code;
  int count;
  List<Moment> data;
  String msg;

  factory MomentListResponse.fromJson(Map<String, dynamic> json) => _$MomentListResponseFromJson(json);

  Map<String, dynamic> toJson() => _$MomentListResponseToJson(this);
}