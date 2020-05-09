import 'dart:async';

import 'package:grpc/grpc.dart';

import 'package:hello/src/generated/helloworld.pb.dart';
import 'package:hello/src/generated/helloworld.pbgrpc.dart';

class GreeterService extends GreeterServiceBase {
  @override
  Future<HelloReply> sayHello(ServiceCall call, HelloRequest request) async {
    return HelloReply()..message = 'Dart, ${request.name}!';
  }
}

Future<void> main(List<String> args) async {
  final server = Server([GreeterService()]);
  await server.serve(port: 50051);
  print('Server listening on port ${server.port}...');
}