
import 'package:app/utils/observer.dart';
import 'package:app/service/upload.dart';
import 'package:app/service/user.dart';
import 'package:dio/dio.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:hive/hive.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'package:app/utils/dio.dart' as $dio;
import 'package:path/path.dart' as $path;

final globalService = GlobalService.instance;

class GlobalService{
  GlobalService._();

  static GlobalService? _instance;

  static GlobalService get instance => _instance ??= GlobalService._();

  Subject<CallOptions> subject = Subject();
  late final UserClient userClient = Get.put(UserClient(subject));
  late final UploadClient uploadClient = Get.put(UploadClient(subject));
  late final Dio httpClient = $dio.httpClient;



  late final Box box;
  late final Database db;

  init() async {
    final appDocDir = await getApplicationDocumentsDirectory();

    final boxfuture = () async {
      box = await Hive.openBox('box', path: $path.join(appDocDir.path, "hive"));
    };

    final dbfuture = () async {
      db = await openDatabase(
          $path.join(appDocDir.path, 'databases', 'hoper.db'),
          version: 1, onCreate: (Database db, int version) async {
        // When creating the db, create the table
        await db.execute(
            'CREATE TABLE hoper (id INTEGER PRIMARY KEY, name TEXT, value INTEGER, num REAL)');
      });

    };
    await Future.wait([boxfuture(), dbfuture()]);

  }
}