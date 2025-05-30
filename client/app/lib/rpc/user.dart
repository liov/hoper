import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/user/user.service.pbgrpc.dart';


import 'package:app/global/dio.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:app/model/response.dart';

import 'package:applib/util/observer.dart';


class UserClient extends Observer<CallOptions> {


  final channel = ClientChannel('grpc.hoper.xyz');


  late  UserServiceClient stub ;

  UserClient(Subject<CallOptions> subject){
    setOptions(subject.options);
    subject.attach(this);
  }

  setOptions(CallOptions? options){
    stub =  UserServiceClient(channel,options:options);
  }


  @override
  update(CallOptions? options){
    if(options!=null) setOptions(options);
  }

  Future<User?> Login(String account, password) async {

    var api = '/v1/login';
    try {
      var response = await httpClient.post(api,data:{'input': account, 'password': password, 'vCode':'5678'});
      return User.create()..mergeFromJsonMap(response.getData((v) => v as Map<String, dynamic>));
    } catch (exception) {
      return null;
    }
  }

}

