import 'package:app/global/const.dart';
import 'package:dio/dio.dart';

Dio httpClient = Dio(BaseOptions(
    baseUrl: BASE_API_URL,
    connectTimeout: const Duration(seconds: 5),
    receiveTimeout: const Duration(seconds: 3),
    headers: {}));
