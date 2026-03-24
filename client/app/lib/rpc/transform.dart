import 'package:dio/dio.dart';


class HybridTransformer extends BackgroundTransformer {
  @override
  Future<dynamic> transformResponse(
    RequestOptions options,
    ResponseBody responseBody,
  ) async {
    final contentType = responseBody.headers['content-type']?.firstOrNull ?? '';

    if (contentType.contains('application/x-protobuf') ||
        contentType.contains('application/protobuf')) {
      // 返回原始字节，调用方手动 fromBuffer()
      final bytes = <int>[];
      await for (final chunk in responseBody.stream) {
        bytes.addAll(chunk);
      }
      return bytes;
    }

    // 默认走 JSON 路径
    return super.transformResponse(options, responseBody);
  }
}
