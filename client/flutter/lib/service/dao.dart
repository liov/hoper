import 'package:get/get_state_manager/src/rx_flutter/rx_disposable.dart';
import 'package:hive/hive.dart';
import 'package:sqflite/sqflite.dart';
import 'package:app/utils/dio.dart' as dio;

class Dao extends GetxService{
  final httpClient = dio.httpClient;
  late final Box box;
  late final Database db;
}