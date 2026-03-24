import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/user/user.service.pbgrpc.dart';

import 'package:app/global/dio.dart';
import 'package:dio/dio.dart' hide Headers;
import 'package:grpc/grpc.dart';
import 'package:app/model/response.dart';
import 'package:app/rpc/grpc.dart';
import 'package:applib/util/observer.dart';
import 'package:retrofit/retrofit.dart';

part 'user.g.dart';

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

@RestApi(baseUrl: 'http://api.hoper.xyz/api', headers: {'User-Agent': 'app/1.0.0', r'Content-Type': 'application/x-protobuf'})
abstract class UserClient {
  factory UserClient(
    Dio dio, {
    String? baseUrl,
    ParseErrorLogger? errorLogger,
  }) = _UserClient;

  @POST('/v1/login')
  @Headers(<String, String>{r'accept': 'application/json', r'Content-Type': 'application/json'})
  @DioResponseType(ResponseType.json)
  Future<User> login(@Body() LoginReq request);

  @POST('/v1/login')
  Future<User> loginPB(@Body() LoginReq request);
}