import 'package:json_annotation/json_annotation.dart';

part 'upload.g.dart';

@JsonSerializable()
class UploadInfo {
  UploadInfo();
  late int id;
  late String url;

  factory UploadInfo.fromJson(Map<String, dynamic> json) => _$UploadInfoFromJson(json);

  Map<String, dynamic> toJson() => _$UploadInfoToJson(this);
}

@JsonSerializable()
class MultiUploadRep  {
  MultiUploadRep();
  late int id;
  late String url;
  late bool success;

  factory MultiUploadRep.fromJson(Map<String, dynamic> json) => _$MultiUploadRepFromJson(json);

  Map<String, dynamic> toJson() => _$MultiUploadRepToJson(this);
}