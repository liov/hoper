import 'package:app/util/dio.dart';
import 'package:flutter/cupertino.dart';

import '../user.dart';

class UserInfo extends ChangeNotifier {
  User user;
  String cookie;
  bool isActive;
}
