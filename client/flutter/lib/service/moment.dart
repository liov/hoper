import 'package:app/generated/protobuf/content/moment.service.pbgrpc.dart';
import 'package:app/model/moment/moment.service.dart';


import 'package:app/utils/dio.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:app/model/response.dart';

class MomentClient extends GetxService {
  late final MomentServiceClient stub;

  MomentClient() : super() {
    final channel = ClientChannel(
      'hoper.xyz',
      port: 8090,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    this.stub = MomentServiceClient(channel);
  }

  Future<MomentListResponse$?> getMomentList(int pageNo, pageSize) async {
    print(pageNo);
    var api = '/v1/moment?pageNo=$pageNo&pageSize=$pageSize';

    try {
      var response = await httpClient.get(api);
      return MomentListResponse$.fromJson(response.getData());
    } catch (exception) {
      return null;
    }
  }
}
