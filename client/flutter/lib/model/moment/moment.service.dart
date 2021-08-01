import 'package:app/model/moment/moment.dart';
import 'package:app/model/user.dart';

import 'package:json_annotation/json_annotation.dart';

part 'moment.service.g.dart';

@JsonSerializable()
class MomentListResponse$ {
  MomentListResponse$();

  late List<User> users;
  late List<Moment$> list;
  late int total;

  factory MomentListResponse$.fromJson(Map<String, dynamic> json) =>
      _$MomentListResponse$FromJson(json);

  Map<String, dynamic> toJson() => _$MomentListResponse$ToJson(this);
}
