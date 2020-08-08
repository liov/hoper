import 'package:json_annotation/json_annotation.dart';

@JsonSerializable(nullable: false)
class Moment{
  int id;
  DateTime createdAt;
  String content;

}