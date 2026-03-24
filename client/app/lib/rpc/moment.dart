import 'package:app/generated/protobuf/content/moment.service.pbgrpc.dart';
import 'package:app/model/moment/moment.service.dart';

import 'package:dio/dio.dart' hide Headers;
import 'package:grpc/grpc.dart';
import 'package:app/rpc/grpc.dart';
import 'package:applib/util/observer.dart';
import 'package:retrofit/retrofit.dart';

part 'moment.g.dart';

class MomentGrpcClient extends Observer<CallOptions> {
  late MomentServiceClient stub;

  MomentGrpcClient(Subject<CallOptions> subject) {
    setOptions(subject.options);
    subject.attach(this);
  }

  void setOptions(CallOptions? options) {
    stub = MomentServiceClient(channel, options: options);
  }

  @override
  void update(CallOptions? options) {
    if (options != null) setOptions(options);
  }
}


@RestApi(baseUrl: 'http://api.hoper.xyz/api')
abstract class MomentClient {
  factory MomentClient(
    Dio dio, {
    String? baseUrl,
    ParseErrorLogger? errorLogger,
  }) = _MomentClient;


  @GET('/moment')
  Future<MomentListResponse$?> getMomentList({@Queries() MomentListReq message});

  @GET('/moment')
  @Headers(<String, String>{'accept': 'application/x-protobuf'})
  @DioResponseType(ResponseType.bytes)
  Future<MomentListResp> getMomentListPB(@Body() MomentListReq message);
}