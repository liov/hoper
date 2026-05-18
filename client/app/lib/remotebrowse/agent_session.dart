import 'package:app/remotebrowse/rb_agent_native.dart';

class RbAgentSession {
  static Future<void> run(Uri signalWs, String room, String root) {
    return RbAgentNative.run(signalWs.toString(), room, root);
  }
}
