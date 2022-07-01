import 'dart:developer' as developer;

import 'package:logging/logging.dart';

var isLogEnable = true;

void defaultLogWriterCallback(String value, {bool isError = false}) {
  if (isError || isLogEnable) developer.log(value, name: 'Hoper');
}


Logger getLogger(String name) {
  var logger = Logger(name);
  logger.level = Level.ALL; // defaults to Level.INFO
  logger.onRecord.listen((record) {
    print('${record.level.name}: ${record.time}: ${record.message}');
  });
  return logger;
}