import 'dart:async';
import 'dart:io';

import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_ice_ffi_io.dart';

/// Agent：信令 + 打洞；列举/缩略图路径由 Viewer 在 wire 里指定。
class RbAgentNative {
  static const _poll = Duration(milliseconds: 150);

  static Future<void> run(String signalWs, String room, {String? sandbox, int iceTimeoutMs = 15000}) async {
    final sb = sandbox ?? '';
    rbLog.info('agent run signal=$signalWs room=$room sandbox=$sb iceMs=$iceTimeoutMs');
    if (RbIceFfi.agentRunAvailable) {
      await _runFfiForever(signalWs, room, sb, iceTimeoutMs);
      return;
    }
    if (Platform.isWindows || Platform.isLinux || Platform.isMacOS) {
      rbLog.info('agent spawn rb process');
      await _runProcessForever(signalWs, room, sb);
      return;
    }
    throw StateError('移动端请用 PC 运行 rb Agent；Viewer 可正常浏览远端相册');
  }

  /// Agent 异常退出时自动重启，仅在外层取消/页面销毁时结束。
  static Future<void> _runFfiForever(String signalWs, String room, String sb, int iceTimeoutMs) async {
    var backoff = const Duration(seconds: 1);
    while (true) {
      if (RbIceFfi.agentRunning() == 0) {
        final rc = RbIceFfi.agentRun(signalWs, room, sb, iceTimeoutMs);
        if (rc == -2) {
          await Future<void>.delayed(_poll);
          continue;
        }
        if (rc != 0) {
          throw StateError('rb_agent_run 失败: $rc');
        }
      }
      while (RbIceFfi.agentRunning() != 0) {
        await Future<void>.delayed(_poll);
      }
      rbLog.warning('agent ffi stopped, restart in ${backoff.inSeconds}s');
      await Future<void>.delayed(backoff);
      backoff = Duration(seconds: (backoff.inSeconds * 2).clamp(1, 30));
    }
  }

  /// 桌面：`rb <房间码>`；进程退出则退避重启，不结束共享会话。
  static Future<void> _runProcessForever(String signalWs, String room, String sandbox) async {
    final bin = await _findBinary('rb');
    final env = Map<String, String>.from(Platform.environment);
    env['RB_SIGNAL_URL'] = signalWs;
    env.remove('RB_AGENT_SANDBOX');
    if (sandbox.isNotEmpty) {
      env['RB_AGENT_SANDBOX'] = sandbox;
    }
    var backoff = const Duration(seconds: 1);
    while (true) {
      rbLog.info('agent exec $bin room=$room RB_SIGNAL_URL=$signalWs');
      final proc = await Process.start(bin, [room], environment: env);
      final code = await proc.exitCode;
      rbLog.warning('agent process exit code=$code, restart in ${backoff.inSeconds}s');
      await Future<void>.delayed(backoff);
      backoff = Duration(seconds: (backoff.inSeconds * 2).clamp(1, 30));
    }
  }

  static Future<String> _findBinary(String name) async {
    final dir = Platform.environment['RB_BIN_DIR'];
    if (dir != null && dir.isNotEmpty) {
      final p = '$dir${Platform.pathSeparator}$name';
      if (await File(p).exists()) {
        return p;
      }
    }
    return name;
  }
}
