import 'package:app/generated/protobuf/content/content.service.pbgrpc.dart';
import 'package:app/rpc/grpc.dart';

import 'package:grpc/grpc.dart';

import 'package:applib/util/observer.dart';

class ContentGrpcClient extends Observer<CallOptions> {
  late ContentServiceClient stub;

  ContentGrpcClient(Subject<CallOptions> subject) {
    setOptions(subject.options);
    subject.attach(this);
  }

  void setOptions(CallOptions? options) {
    stub = ContentServiceClient(channel, options: options);
  }

  @override
  void update(CallOptions? options) {
    if (options != null) setOptions(options);
  }
}
