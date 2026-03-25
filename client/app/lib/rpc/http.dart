import 'package:app/global/const.dart';

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
    headers: {
      r'User-Agent': 'app/1.0.0',
      r'Content-Type': 'application/x-protobuf',
      r'accept': 'application/x-protobuf',
    },
  ),
);
