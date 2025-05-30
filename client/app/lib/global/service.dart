
import 'package:applib/util/observer.dart';
import 'package:app/rpc/upload.dart';
import 'package:app/rpc/user.dart';
import 'package:dio/dio.dart';
import 'package:flutter_cache_manager/flutter_cache_manager.dart';
import 'package:get/get.dart' hide Response;
import 'package:grpc/grpc.dart' hide Response;
import 'package:hive_ce/hive.dart';
import 'package:logger/logger.dart';
import 'package:path_provider/path_provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:sqflite/sqflite.dart';
import 'package:app/global/dio.dart' as $dio;
import 'package:path/path.dart' as $path;
import 'package:sqflite_common_ffi/sqflite_ffi.dart';

final globalService = GlobalService.instance;

class GlobalService{
  GlobalService._();

  static GlobalService? _instance;

  static GlobalService get instance => _instance ??= GlobalService._();

  final logger = Logger(printer: HybridPrinter(PrettyPrinter(), debug: SimplePrinter()),level:Level.verbose);

  Subject<CallOptions> subject = Subject(CallOptions(timeout: const Duration(seconds: 5)));
  set callOptions(CallOptions callOptions)=> subject.setState(callOptions);

  late final UserClient userClient = Get.put(UserClient(subject));
  late final UploadClient uploadClient = Get.put(UploadClient(subject));

  late final Box box;
  late final Database db;
  late final SharedPreferences shared;
  final cache = DefaultCacheManager();



  init() async {

    final appDocDir = await getApplicationDocumentsDirectory();
    logger.d(appDocDir.path);
    final appSupportDir = await getApplicationSupportDirectory();
    logger.d(appSupportDir.path);
    //final dbpath = await getDatabasesPath();
    $dio.httpClient.interceptors.add(
      InterceptorsWrapper(
        onRequest: (RequestOptions options, RequestInterceptorHandler handler) {
// Do something before request is sent.
// If you want to resolve the request with custom data,
// you can resolve a `Response` using `handler.resolve(response)`.
// If you want to reject the request with a error message,
// you can reject with a `DioException` using `handler.reject(dioError)`.
          return handler.next(options);
        },
        onResponse: (Response response, ResponseInterceptorHandler handler) {
// Do something with response data.
// If you want to reject the request with a error message,
// you can reject a `DioException` object using `handler.reject(dioError)`.
          return handler.next(response);
        },
        onError: (DioException error, ErrorInterceptorHandler handler) {
// Do something with response error.
// If you want to resolve the request with some custom data,
// you can resolve a `Response` object using `handler.resolve(response)`.
          return handler.next(error);
        },
      ),
    );
    boxfuture() async {
      box = await Hive.openBox('box', path: $path.join(appDocDir.path, "hive"));
    }

    dbfuture() async {
      databaseFactory = databaseFactoryFfi;
      db = await openDatabase(
          $path.join(appDocDir.path,'database', 'hoper.db'),
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

    }
    await Future.wait([boxfuture(), dbfuture()]);
    shared = await SharedPreferences.getInstance();
  }
}