import 'package:app/global/const.dart';
import 'package:app/global/service.dart';
import 'package:dio/dio.dart';

Dio httpClient = Dio(
  BaseOptions(
    baseUrl: BASE_API_URL,
    connectTimeout: const Duration(seconds: 5),
    receiveTimeout: const Duration(seconds: 3),
    headers: {
      r'User-Agent': 'app/1.0.0',
      r'Content-Type': 'application/json',
      r'accept': 'application/json',
    },
  ),
);

Dio httpProtobufClient = Dio(
  BaseOptions(
    baseUrl: BASE_API_URL,
    connectTimeout: const Duration(seconds: 5),
    receiveTimeout: const Duration(seconds: 3),
    responseType: ResponseType.bytes,
    contentType: 'application/x-protobuf',
    headers: {
      r'User-Agent': 'app/1.0.0',
      r'Content-Type': 'application/x-protobuf',
      r'accept': 'application/x-protobuf',
    },
  ),
);

final interceptor = InterceptorsWrapper(
  onRequest: (RequestOptions options, RequestInterceptorHandler handler) {
    globalService.logger.fine('${options.method} ${options.uri}');
    return handler.next(options);
  },
  onResponse: (Response response, ResponseInterceptorHandler handler) {
    globalService.logger.fine(
      '${response.statusCode} ${response.requestOptions.uri}',
    );
    return handler.next(response);
  },
  onError: (DioException error, ErrorInterceptorHandler handler) {
    globalService.logger.severe(
      'dio error ${error.requestOptions.uri}',
      error,
      error.stackTrace,
    );
    return handler.resolve(error.response ?? Response<dynamic>(requestOptions: error.requestOptions, data: null));
  },
);
