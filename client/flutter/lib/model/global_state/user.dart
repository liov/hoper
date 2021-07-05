
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';

import '../user.dart' as $self;

class UserState {
  var users$ = Map<int, $self.User>();
  var users = Map<Int64, UserBaseInfo>();

  UserBaseInfo? getUser(Int64 id){
    return users[id];
  }
}