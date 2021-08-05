import 'package:get/get_state_manager/src/rx_flutter/rx_disposable.dart';
import 'package:hive/hive.dart';
import 'package:sqflite/sqflite.dart';


class Dao extends GetxService{
  late final Box box;
  late final Database db;
}