import 'package:json_annotation/json_annotation.dart';
part 'webview_message.g.dart';

@JsonSerializable()
class WebviewMessage{
  WebviewMessage({required this.method, required this.params});
  final String method;
  final List<dynamic> params;

  factory WebviewMessage.fromJson(Map<String, dynamic> json) => _$WebviewMessageFromJson(json);

  Map<String, dynamic> toJson() => _$WebviewMessageToJson(this);
}
