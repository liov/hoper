import 'dart:io';

import 'package:app/generated/protobuf/content/moment.service.pbgrpc.dart';
import 'package:app/model/moment/moment.service.dart';

import 'package:app/global/dio.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:app/model/response.dart';
import 'package:app/rpc/grpc.dart';
import 'package:applib/util/observer.dart';


class MomentClient extends Observer<CallOptions> {

  late MomentServiceClient stub;

  MomentClient(Subject<CallOptions> subject) {
    setOptions(subject.options);
    subject.attach(this);
  }

  setOptions(CallOptions? options) {
    stub = MomentServiceClient(channel, options: options);
  }

  Future<MomentListResponse$?> getMomentList(int pageNo, pageSize) async {

    var api = '/v1/moment?pageNo=$pageNo&pageSize=$pageSize';

    try {
      var response = await httpClient.get(api);
      return response.getData((v) => MomentListResponse$.fromJson(v as Map<String, dynamic>));
    } catch (exception) {
      return null;
    }
  }

  @override
  void update(CallOptions? options) {
    if (options != null) setOptions(options);
  }
}
