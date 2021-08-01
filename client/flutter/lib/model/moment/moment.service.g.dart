// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'moment.service.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

MomentListResponse$ _$MomentListResponse$FromJson(Map<String, dynamic> json) {
  return MomentListResponse$()
    ..users = (json['users'] as List<dynamic>)
        .map((e) => User.fromJson(e as Map<String, dynamic>))
        .toList()
    ..list = (json['list'] as List<dynamic>)
        .map((e) => Moment$.fromJson(e as Map<String, dynamic>))
        .toList()
    ..total = json['total'] as int;
}

Map<String, dynamic> _$MomentListResponse$ToJson(MomentListResponse$ instance) =>
    <String, dynamic>{
      'users': instance.users,
      'list': instance.list,
      'total': instance.total,
    };
