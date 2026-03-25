import 'package:app/generated/protobuf/content/moment.service.pbgrpc.dart';
import 'package:app/model/moment.dart';
import 'package:dio/dio.dart';
import 'package:grpc/grpc.dart';
import 'package:app/rpc/grpc.dart';
import 'package:applib/util/observer.dart';
import 'package:app/rpc/http.dart';

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

class MomentClient {
  Future<MomentListResp> getMomentList(MomentListReq message) async {
    final response = await httpClient.get(
      '/moment',
      queryParameters: message.writeToJsonMap(),
    );
    return MomentListResp.fromJson(response.data);
  }

  Future<MomentListResp> getMomentListPB(MomentListReq message) async {
    final response = await httpProtobufClient.get(
      '/moment',
      queryParameters: message.writeToJsonMap(),
    );
    return MomentListResp.fromJson(response.data);
  }
}
