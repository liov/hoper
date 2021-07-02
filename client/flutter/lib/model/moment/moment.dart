import 'package:json_annotation/json_annotation.dart';

import '../content.dart';

part 'moment.g.dart';

@JsonSerializable()
class Moment$ {
  Moment$();
  late int id;
  late DateTime createdAt;
  late String content;
  late String? images;
  late TinyTag? mood;
  late TinyTag? weather;
  late List<TinyTag>? tags;

  late int userId;
  late int likeId;
  late int unlikeId;
  late bool collect;
  late int sequence;
  late int anonymous;
  late int permission;

  factory Moment$.fromJson(Map<String, dynamic> json) => _$MomentFromJson(json);

  Map<String, dynamic> toJson() => _$MomentToJson(this);
}
