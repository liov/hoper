import 'package:pigeon/pigeon_lib.dart';

class StringMessage {
  String? route;
}

@HostApi()
abstract class RouteApi {
  void initialize();
  void toNative(StringMessage msg);
}

void configurePigeon(PigeonOptions opts) {
  opts.dartOut = 'lib/generated/pigeons/route.dart';
  opts.objcHeaderOut = 'ios/Classes/route.h';
  opts.objcSourceOut = 'ios/Classes/route.m';
  opts.objcOptions!.prefix = 'FLT';
  opts.javaOut =
  'android/app/src/main/java/io/flutter/plugins/Route.java';
  opts.javaOptions!.package = 'io.flutter.plugins';
}