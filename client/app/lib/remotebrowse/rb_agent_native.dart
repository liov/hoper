import 'dart:async';
import 'dart:io';

import 'package:app/remotebrowse/rb_ice_ffi.dart';

/// Agent：信令 + 打洞；列举/缩略图路径由 Viewer 在 wire 里指定。
class RbAgentNative {
  static const _poll = Duration(milliseconds: 150);

  static Future<void> run(String signalWs, String room, {String? sandbox, int iceTimeoutMs = 15000}) async {
    final sb = sandbox ?? '';
    if (RbIceFfi.agentRunAvailable) {
      final rc = RbIceFfi.agentRun(signalWs, room, sb, iceTimeoutMs);
      if (rc == -2) {
        throw StateError('Agent 已在运行');
      }
      if (rc != 0) {
        throw StateError('rb_agent_run 失败: $rc');
      }
      while (RbIceFfi.agentRunning() != 0) {
        await Future<void>.delayed(_poll);
      }
      return;
    }
    if (Platform.isWindows || Platform.isLinux || Platform.isMacOS) {
      await _runProcess(signalWs, room, sb);
      return;
    }
    throw StateError('请使用 --features transport 编译的 librfv');
  }

  /// 桌面：`rfv <房间码>`，与子命令无关。
  static Future<void> _runProcess(String signalWs, String room, String sandbox) async {
    final bin = await _findBinary('rfv');
    final env = <String, String>{
      'RB_SIGNAL_URL': signalWs,
      if (sandbox.isNotEmpty) 'RB_AGENT_SANDBOX': sandbox,
    };
    final proc = await Process.start(bin, [room], environment: env);
    final code = await proc.exitCode;
    if (code != 0) {
      throw StateError('rfv 退出码 $code');
    }
  }

  static Future<String> _findBinary(String name) async {
    final dir = Platform.environment['RFV_BIN_DIR'];
    if (dir != null && dir.isNotEmpty) {
      final p = '$dir${Platform.pathSeparator}$name';
      if (await File(p).exists()) {
        return p;
      }
    }
    return name;
  }
}
