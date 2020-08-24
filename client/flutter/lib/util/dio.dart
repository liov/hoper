
import 'package:dio/dio.dart';

Dio _dio;

_getHttpClient(){
  // 或者通过传递一个 `options`来创建dio实例
  BaseOptions options = BaseOptions(
      baseUrl: "https://hoper.xyz/api",
      connectTimeout: 5000,
      receiveTimeout: 3000,
      headers:{"Authorization": ""}
  );
  _dio = Dio(options);
}



Dio httpClient(){
  if(_dio==null) _getHttpClient();
  return _dio;
}