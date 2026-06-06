import 'dart:developer' as developer;
import 'dart:io';
import 'package:logging/logging.dart';

const _reset = '\x1B[0m';
const _levelColors = {
  'FINEST': '\x1B[90m', // 灰
  'FINER': '\x1B[90m',
  'FINE': '\x1B[36m', // 青
  'CONFIG': '\x1B[34m', // 蓝
  'INFO': '\x1B[32m', // 绿
  'WARNING': '\x1B[33m', // 黄
  'SEVERE': '\x1B[31m', // 红
  'SHOUT': '\x1B[1;31m', // 粗体红
};

// 只保留属于 package:app/ 的帧，过滤掉框架内部噪音
bool _stackUseful(StackTrace? st) {
  if (st == null) {
    return false;
  }
  final raw = st.toString().trim();
  if (raw.isEmpty || raw == '_StringStackTrace ()' || raw == '_StringStackTrace') {
    return false;
  }
  return true;
}

String _trimStack(StackTrace? st) {
  if (!_stackUseful(st)) {
    return '';
  }
  final lines = st!
      .toString()
      .split('\n')
      .where((l) => l.contains('package:app/'))
      .toList();
  return lines.isEmpty
      ? st.toString().split('\n').take(5).join('\n')
      : lines.join('\n');
}

// 找调用栈里第一个属于本 app 的帧，返回 "path/file.dart:line"
String _caller() {
  // package:app/pages/home/home_view.dart:47:5
  final re = RegExp(r'package:app/(.+\.dart):(\d+):\d+\)');
  for (final frame in StackTrace.current.toString().split('\n')) {
    if (frame.contains('global/logger.dart')) continue;
    final m = re.firstMatch(frame);
    if (m != null) return '${m.group(1)}:${m.group(2)}';
  }
  return '';
}

StackTrace? _outStackTrace(LogRecord r) {
  if (_stackUseful(r.stackTrace)) {
    return r.stackTrace;
  }
  final trimmed = _trimStack(r.stackTrace);
  if (trimmed.isEmpty) {
    return null;
  }
  return StackTrace.fromString(trimmed);
}

void _logRecord(LogRecord r) {
  final time = r.time.toIso8601String();
  final caller = _caller();
  final prefix = '$time [${r.level.name}] $caller';
  final stack = _trimStack(r.stackTrace);
  final outStack = _outStackTrace(r);
  if (stdout.supportsAnsiEscapes) {
    final color = _levelColors[r.level.name] ?? '';
    final msg = r.error != null
        ? '${r.message}\n${r.error}${stack.isNotEmpty ? '\n$stack' : ''}'
        : r.message;
    stderr.writeln('$color$prefix: $msg$_reset');
  } else {
    final body = r.error != null ? '${r.message}\n${r.error}' : r.message;
    if (outStack != null) {
      developer.log('$prefix: $body', time: r.time, level: r.level.value, name: r.loggerName, stackTrace: outStack);
    } else {
      developer.log('$prefix: $body', time: r.time, level: r.level.value, name: r.loggerName);
    }
  }
}

void setupLogger() {
  Logger.root.level = Level.ALL;
  Logger.root.onRecord.listen(_logRecord);
  recordStackTraceAtLevel = Level.OFF;
}

/// runZonedGuarded 未捕获异常；空堆栈不重复打印 `_StringStackTrace ()`。
void logZoneError(Object error, StackTrace stack) {
  final useful = _stackUseful(stack);
  Logger('app').warning(
    useful ? '未捕获异常' : '未捕获异常（无堆栈）',
    error,
    useful ? stack : null,
  );
}

final logger = Logger('app');
