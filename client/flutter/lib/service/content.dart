import 'package:app/generated/protobuf/content/content.service.pbgrpc.dart';


import 'package:grpc/grpc.dart';


import '../utils/observer.dart';


class ContentClient extends Observer<CallOptions> {


  final channel = ClientChannel(
    'hoper.xyz',
    port: 8090,
    options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
  );

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
