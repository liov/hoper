
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';

import '../user.dart' as $self;

class UserState extends GetxController {
  var usersS = Map<int, $self.User>().obs;
  var users = Map<Int64, UserBaseInfo>().obs;
}