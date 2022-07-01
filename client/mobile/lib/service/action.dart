import 'package:app/generated/protobuf/content/action.service.pbgrpc.dart';


import 'package:grpc/grpc.dart';


import '../utils/observer.dart';


class ActionClient extends Observer<CallOptions> {


  final channel = ClientChannel('grpc.hoper.xyz',);

  late ActionServiceClient stub;

  ActionClient(Subject<CallOptions> subject){
    setOptions(subject.options);
    subject.attach(this);
  }

  setOptions(CallOptions? options){
    stub =  ActionServiceClient(channel,options:options);
  }

  @override
  void update(CallOptions? options) {
    if(options!=null) setOptions(options);
  }
}
