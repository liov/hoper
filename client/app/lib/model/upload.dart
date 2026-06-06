import 'package:json_annotation/json_annotation.dart';

part 'upload.g.dart';

@JsonSerializable()
class UploadInfo {
  UploadInfo({required this.id, required this.url});
   int id;
   String url;

  factory UploadInfo.fromJson(Map<String, dynamic> json) => _$UploadInfoFromJson(json);

  Map<String, dynamic> toJson() => _$UploadInfoToJson(this);
}

@JsonSerializable()
class MultiUploadResp  {
  MultiUploadResp({required this.id, required this.url, required this.success});
   int id;
   String url;
   bool success;

  factory MultiUploadResp.fromJson(Map<String, dynamic> json) => _$MultiUploadRespFromJson(json);

  Map<String, dynamic> toJson() => _$MultiUploadRespToJson(this);
}