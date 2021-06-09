
import 'package:app/model/moment.dart';
import 'package:app/model/user.dart';
import 'package:app/service/response.dart';
import 'package:app/util/dio.dart';
import 'package:json_annotation/json_annotation.dart';

part 'moment.g.dart';

@JsonSerializable()
class MomentListResponse {
  MomentListResponse();

  late List<User> users;
  late List<Moment> list;
  late int total;

  factory MomentListResponse.fromJson(Map<String, dynamic> json) =>
      _$MomentListResponseFromJson(json);

  Map<String, dynamic> toJson() => _$MomentListResponseToJson(this);
}

Future<MomentListResponse?> getMomentList(int pageNo, pageSize) async {
  var api = '/v1/moment?page=$pageNo&pageSize=$pageSize';

  try {
    var response = await httpClient.get(api);
    return MomentListResponse.fromJson(response.getData());
  } catch (exception) {
    print(exception);
  }
  return null;
}
