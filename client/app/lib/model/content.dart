import 'package:json_annotation/json_annotation.dart';

part 'content.g.dart';

@JsonSerializable()
class TinyTag{
  TinyTag();
  late int id;
  late String name;
  late TagType type;

  factory TinyTag.fromJson(Map<String, dynamic> json) => _$TinyTagFromJson(json);

  Map<String, dynamic> toJson() => _$TinyTagToJson(this);
}


enum TagType{
  @JsonValue(0)
  Placeholder,
  @JsonValue(1)
  Content,
  @JsonValue(2)
  Mood,
  @JsonValue(3)
  Weather,
  @JsonValue(4)
  Location,
  @JsonValue(5)
  Category,
  @JsonValue(6)
  Notice
}