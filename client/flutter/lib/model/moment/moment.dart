import 'package:app/generated/user/user.model.pb.dart';
import 'package:json_annotation/json_annotation.dart';

@JsonSerializable(nullable: false)
class Moment{
  int id;
  DateTime createdAt;
  String content;
  String imageUrl;
  Mood mood;
  List<Tag> tags;

}

@JsonSerializable(nullable: false)
class Mood{
  String name;
  String description;
  String expressionURL;
  int status;
}

@JsonSerializable(nullable: false)
class Tag{
  String name;
  String description;
  int status;
}

@JsonSerializable(nullable: false)
class Category{
  int id;
  String name;
  int parentId;
  int sequence;
  int status;
}

@JsonSerializable(nullable: false)
class MomentComment{
  int id;
  DateTime createdAt;
  User user;
}