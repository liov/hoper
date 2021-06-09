import 'package:json_annotation/json_annotation.dart';

part 'content.g.dart';

@JsonSerializable()
class TinyTag{
  TinyTag();
  late int id;
  late String name;
  late int type;

  factory TinyTag.fromJson(Map<String, dynamic> json) => _$TinyTagFromJson(json);

  Map<String, dynamic> toJson() => _$TinyTagToJson(this);
}

enum TagType{
  Placeholder,
  Content,
  Mood,
  Weather,
  Location,
  Category,
  Notice
}