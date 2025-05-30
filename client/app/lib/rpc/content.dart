import 'package:app/generated/protobuf/content/content.service.pbgrpc.dart';


import 'package:grpc/grpc.dart';


import 'package:applib/util/observer.dart';


class ContentClient extends Observer<CallOptions> {


  final channel = ClientChannel('grpc.hoper.xyz');

  late ContentServiceClient stub;

  ContentClient(Subject<CallOptions> subject){
    setOptions(subject.options);
    subject.attach(this);
  }

  setOptions(CallOptions? options){
    stub =  ContentServiceClient(channel,options:options);
  }

  @override
  void update(CallOptions? options) {
    if(options!=null) setOptions(options);
  }
}
