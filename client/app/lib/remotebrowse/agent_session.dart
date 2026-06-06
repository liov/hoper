import 'package:app/remotebrowse/rb_agent_native.dart';

class RbAgentSession {
  /// 仅注册房间并等待 P2P；要浏览的目录由「浏览」端填写并下发。
  static Future<void> run(Uri signalWs, String room, {String? sandbox}) {
    return RbAgentNative.run(signalWs.toString(), room, sandbox: sandbox);
  }
}
