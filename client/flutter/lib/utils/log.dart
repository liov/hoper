import 'dart:developer' as developer;

import 'package:logging/logging.dart' as logging;
import 'package:logging/logging.dart';

var isLogEnable = true;

void defaultLogWriterCallback(String value, {bool isError = false}) {
  if (isError || isLogEnable) developer.log(value, name: 'Hoper');
}


logging.Logger getLogger(String name) {
  var logger = logging.Logger(name);
  logger.level = logging.Level.ALL; // defaults to Level.INFO
  logger.onRecord.listen((record) {
    print('${record.level.name}: ${record.time}: ${record.message}');
  });
  return logger;
}

abstract class Logger {
  Level get level;
  set level(Level value);
  Stream<LogRecord> get onRecord;
  void log(Level level, String message, [Object error, StackTrace stackTrace]);
  void finest(String message, [Object error, StackTrace stackTrace]);
  void finer(String message, [Object error, StackTrace stackTrace]);
  void info(String message, [Object error, StackTrace stackTrace]);
  void warning(String message, [Object error, StackTrace stackTrace]);
  void severe(String message, [Object error, StackTrace stackTrace]);
}