import 'package:app/generated/protobuf/content/action.service.pbgrpc.dart';
import 'package:app/rpc/grpc.dart';

import 'package:grpc/grpc.dart';


import 'package:applib/util/observer.dart';


class ActionClient extends Observer<CallOptions> {


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
