
import 'package:app/utils/observer.dart';
import 'package:app/service/upload.dart';
import 'package:app/service/user.dart';
import 'package:dio/dio.dart';
import 'package:flutter_cache_manager/flutter_cache_manager.dart';
import 'package:get/get.dart';
import 'package:grpc/grpc.dart';
import 'package:hive/hive.dart';
import 'package:logging/logging.dart';
import 'package:path_provider/path_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:sqflite/sqflite.dart';
import 'package:app/utils/dio.dart' as $dio;
import 'package:path/path.dart' as $path;

final globalService = GlobalService.instance;

class GlobalService{
  GlobalService._();

  static GlobalService? _instance;

  static GlobalService get instance => _instance ??= GlobalService._();

  Subject<CallOptions> subject = Subject(CallOptions(timeout: Duration(seconds: 5)));
  set callOptions(CallOptions callOptions)=> subject.setState(callOptions);

  late final UserClient userClient = Get.put(UserClient(subject));
  late final UploadClient uploadClient = Get.put(UploadClient(subject));
  late final Dio httpClient = $dio.httpClient;



  late final Box box;
  late final Database db;
  late final SharedPreferences shared;
  final cache = DefaultCacheManager();
  final log = Logger.root;


  init() async {
    log.level = Level.ALL; // defaults to Level.INFO
    log.onRecord.listen((record) {
      print('${record.level.name}: ${record.time}: ${record.message}');
    });

    final appDocDir = await getApplicationDocumentsDirectory();

    final boxfuture = () async {
      box = await Hive.openBox('box', path: $path.join(appDocDir.path, "hive"));
    };

    final dbfuture = () async {
      db = await openDatabase(
          $path.join(appDocDir.path, 'databases', 'hoper.db'),
          version: 1, onCreate: (Database db, int version) async {

        // When creating the db, create the table
            await db.execute('''
              create table if not exists hoper (
                id integer primary key autoincrement,
                name text,
                value text,
                num REAL,
                created_at text,
                updated_at text
              )
            ''');

      });

    };
    await Future.wait([boxfuture(), dbfuture()]);
    shared = await SharedPreferences.getInstance();
  }
}