import 'package:dio/dio.dart';

Dio httpClient = Dio(BaseOptions(
    baseUrl: "https://hoper.xyz/api",
    connectTimeout: 5000,
    receiveTimeout: 3000,
    headers: {"Authorization": ""}));
