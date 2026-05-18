import 'dart:async';
import 'dart:io';

import 'package:app/remotebrowse/rb_ice_ffi.dart';

/// Agent 全栈由 Rust 承担：信令、打洞、ffmpeg 数据面。
class RbAgentNative {
  static const _poll = Duration(milliseconds: 150);

  static Future<void> run(String signalWs, String room, String root, {int iceTimeoutMs = 15000}) async {
    if (RbIceFfi.agentRunAvailable) {
      final rc = RbIceFfi.agentRun(signalWs, room, root, iceTimeoutMs);
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
      await _runProcess(signalWs, room, root);
      return;
    }
    throw StateError('请使用 --features transport 编译的 librfv，或在 PATH 中提供 rfv-agent');
  }

  static Future<void> _runProcess(String signalWs, String room, String root) async {
    final bin = await _findBinary('rfv-agent');
    final proc = await Process.start(bin, [room, root], environment: {'RB_SIGNAL_URL': signalWs});
    final code = await proc.exitCode;
    if (code != 0) {
      throw StateError('rfv-agent 退出码 $code');
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
