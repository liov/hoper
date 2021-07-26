import 'package:app/model/const/const.dart';
import 'package:dio/dio.dart';
import 'package:hive/hive.dart';
import 'package:sqflite/sqflite.dart';
import 'package:app/utils/dio.dart' as dio;

final dao = Dao();

class Dao {
  final httpClient = dio.httpClient;
  late final Box box;
  late final Database db;
}