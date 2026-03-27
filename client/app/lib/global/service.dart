
import 'package:applib/util/observer.dart';
import 'package:app/global/logger.dart';
import 'package:app/rpc/upload.dart';
import 'package:app/rpc/user.dart';
import 'package:flutter_cache_manager/flutter_cache_manager.dart';
import 'package:get/get.dart' hide Response;
import 'package:grpc/grpc.dart' hide Response;
import 'package:hive_ce/hive.dart';
import 'package:logging/logging.dart';
import 'package:path_provider/path_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:app/rpc/http.dart' as $dio;
import 'package:path/path.dart' as $path;
import 'package:sqflite_common_ffi/sqflite_ffi.dart';



final globalService = GlobalService.instance;

class GlobalService {
  GlobalService._();

  static GlobalService? _instance;

  static GlobalService get instance => _instance ??= GlobalService._();

  final logger = Logger('app');

  Subject<CallOptions> subject = Subject(
    CallOptions(timeout: const Duration(seconds: 5)),
  );
  set callOptions(CallOptions callOptions) => subject.setState(callOptions);

  late final UserGrpcClient userClient = Get.put(UserGrpcClient(subject));
  late final UploadClient uploadClient = Get.put(UploadClient(subject));

  bool _initialized = false;
  late final Box box;
  late final Database db;
  late final SharedPreferences shared;
  final cache = DefaultCacheManager();

  Future<void> init() async {
    if (_initialized) return;
    _initialized = true;
    setupLogger();
    final appDocDir = await getApplicationDocumentsDirectory();
    logger.fine(appDocDir.path);
    final appSupportDir = await getApplicationSupportDirectory();
    logger.fine(appSupportDir.path);


    $dio.httpClient.interceptors.add($dio.interceptor);
    $dio.httpProtobufClient.interceptors.add($dio.interceptor);
    boxfuture() async {
      box = await Hive.openBox('box', path: $path.join(appDocDir.path, "hive"));
    }

    dbfuture() async {
      sqfliteFfiInit();
      databaseFactory = databaseFactoryFfi;
      db = await openDatabase(
        $path.join(appDocDir.path, 'database', 'hoper.db'),
        version: 1,
        onCreate: (Database db, int version) async {
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
        },
      );
    }

    await Future.wait([boxfuture(), dbfuture()]);
    shared = await SharedPreferences.getInstance();
  }
}
