import 'package:app/generated/protobuf/user/user.model.pb.dart' as $pb;
import 'package:json_annotation/json_annotation.dart';

part 'user.g.dart';

@JsonSerializable()
class User {
  User();
  late int id;
  late String name;
  late int gender;
  late String avatarUrl;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);

  Map<String, dynamic> toJson() => _$UserToJson(this);
}

extension UserConvert on $pb.UserBase {
  $pb.UserBase from($pb.User user){
    return $pb.UserBase(id:user.id,name:user.name,gender: user.gender,avatar: user.avatar);
  }
}