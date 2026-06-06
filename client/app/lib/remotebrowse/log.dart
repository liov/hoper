import 'dart:developer' as developer;
import 'dart:io';

import 'package:logging/logging.dart';

bool _stackUseful(StackTrace? st) {
  if (st == null) {
    return false;
  }
  final raw = st.toString().trim();
  return raw.isNotEmpty && raw != '_StringStackTrace ()' && raw != '_StringStackTrace';
}

void _emitRbLog(LogRecord r) {
  final ts = r.time.toIso8601String();
  final buf = StringBuffer('[$ts] [${r.level.name}] ${r.message}');
  if (r.error != null) {
    buf.write('\n  ${r.error}');
  }
  if (_stackUseful(r.stackTrace)) {
    buf.write('\n${r.stackTrace}');
  }
  final line = buf.toString();
  if (stdout.supportsAnsiEscapes) {
    stderr.writeln('[remotebrowse] $line');
    return;
  }
  developer.log(line, name: 'remotebrowse', time: r.time, level: r.level.value);
}

/// 独立 Logger，不向 [Logger.root] 冒泡，避免 IDE/重复监听多打一行空堆栈。
final Logger rbLog = Logger.detached('remotebrowse')
  ..level = Level.ALL
  ..onRecord.listen(_emitRbLog);
