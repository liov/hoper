
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';

import 'package:app/model/user.dart' as $self;

class UserState {
  var users$ = Map<int, $self.User>();
  var _users = Map<Int64, UserBaseInfo>();

  UserBaseInfo? getUser(Int64 id){
    return _users[id];
  }

  appendUsers(List<UserBaseInfo> users){
    users.forEach((e) => _users[e.id] = e);
  }
}