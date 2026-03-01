
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';

import 'package:app/model/user.dart' as $self;

class UserState {
  var users$ = <int, $self.User>{};
  final _users = <Int64, UserBase>{};

  UserBase? getUser(Int64 id){
    return _users[id];
  }

  appendUsers(List<UserBase> users){
    for (var e in users) {
      _users[e.id] = e;
    }
  }
  append(UserBase? user){
    if (user!=null) {
      _users[user.id] = user;
    }
  }
}