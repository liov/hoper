import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/user/user.service.pbgrpc.dart';


import 'package:app/utils/dio.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:app/model/response.dart';


class UserClient extends GetxService {
  late final UserServiceClient stub;

  UserClient() : super() {
    final channel = ClientChannel(
      'hoper.xyz',
      port: 8090,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    this.stub = UserServiceClient(channel);
  }

  Future<User?> Login(String account, password) async {

    var api = '/v1/login';
    try {
      var response = await httpClient.post(api,data:{'input': account, 'password': password});
      return User.create()..mergeFromJsonMap(response.getData());
    } catch (exception) {
      return null;
    }
  }
}
