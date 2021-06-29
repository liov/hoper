
import 'package:app/generated/protobuf/content/moment.service.pbgrpc.dart';
import 'package:app/model/moment.dart';
import 'package:app/model/user.dart';
import 'package:app/service/response.dart';
import 'package:app/util/dio.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:json_annotation/json_annotation.dart';


part 'moment.g.dart';


class MomentClient extends GetxService {

  late final MomentServiceClient stub ;

  MomentClient():super(){
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

@JsonSerializable()
class MomentListResponse$ {
  MomentListResponse$();

  late List<User> users;
  late List<Moment$> list;
  late int total;

  factory MomentListResponse$.fromJson(Map<String, dynamic> json) =>
      _$MomentListResponseFromJson(json);

  Map<String, dynamic> toJson() => _$MomentListResponseToJson(this);
}
