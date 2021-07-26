import 'package:sqflite/sqflite.dart';
import 'package:sqflite/utils/utils.dart';

Future<bool> tableExists(DatabaseExecutor db, String table) async {
  var count = firstIntValue(await db.query('sqlite_master',
      columns: ['COUNT(*)'],
      where: 'type = ? AND name = ?',
      whereArgs: ['table', table]));
  return count! > 0;
}