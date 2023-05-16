import 'package:pigeon/pigeon_lib.dart';

class StringMessage {
  String? route;
}

@HostApi()
abstract class RouteApi {
  void initialize();
  void toNative(StringMessage msg);
}