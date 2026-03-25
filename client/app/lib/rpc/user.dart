import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/user/user.service.pbgrpc.dart';

import 'package:app/rpc/http.dart';
import 'package:dio/dio.dart' hide Headers;
import 'package:grpc/grpc.dart';
import 'package:app/model/response.dart';
import 'package:app/rpc/grpc.dart';
import 'package:applib/util/observer.dart';
import 'package:app/rpc/http.dart';


class UserGrpcClient extends Observer<CallOptions> {
  late UserServiceClient stub;

  UserGrpcClient(Subject<CallOptions> subject) {
    setOptions(subject.options);
    subject.attach(this);
  }

  void setOptions(CallOptions? options) {
    stub = UserServiceClient(channel, options: options);
  }

  @override
  update(CallOptions? options) {
    if (options != null) setOptions(options);
  }

}


class UserClient {

  Future<User> login(LoginReq request) async {
    final response = await httpClient.post('/login', data: request.writeToJsonMap());
    return User.fromJson(response.data);
  }


  Future<User> loginPB(LoginReq request) async {
    final response = await httpProtobufClient.post('/login', data: request.writeToBuffer());
    return User.fromBuffer(response.data);
  }
}