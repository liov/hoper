import 'package:sqflite/sqflite.dart';
import 'package:sqflite/utils/utils.dart';

Future<bool> tableExists(DatabaseExecutor db, String table) async {
  print('tableExists');
  var count = firstIntValue(await db.query('sqlite_master',
      columns: ['COUNT(*)'],
      where: 'type = table AND name = ?',
      whereArgs: [table]).onError<Exception>((e, _) => throw e));
  return count! > 0;
}