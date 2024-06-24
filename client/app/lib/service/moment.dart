import 'dart:io';

import 'package:app/generated/protobuf/content/moment.service.pbgrpc.dart';
import 'package:app/model/moment/moment.service.dart';

import 'package:app/utils/dio.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:app/model/response.dart';

import '../utils/observer.dart';


class MomentClient extends Observer<CallOptions> {
  final channel = ClientChannel('grpc.hoper.xyz');

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
      return MomentListResponse$.fromJson(response.getData());
    } catch (exception) {
      return null;
    }
  }

  @override
  void update(CallOptions? options) {
    if (options != null) setOptions(options);
  }
}
