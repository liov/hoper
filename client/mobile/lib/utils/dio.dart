import 'package:app/model/const/const.dart';
import 'package:dio/dio.dart';

Dio httpClient = Dio(BaseOptions(
    baseUrl: BASE_API_URL,
    connectTimeout: 5000,
    receiveTimeout: 3000,
    headers: {}));
