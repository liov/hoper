import 'package:app/model/const/const.dart';
import 'package:dio/dio.dart';

Dio httpClient = Dio(BaseOptions(
    baseUrl: baseApiUrl,
    connectTimeout: 5000,
    receiveTimeout: 3000,
    headers: {"Authorization": ""}));
