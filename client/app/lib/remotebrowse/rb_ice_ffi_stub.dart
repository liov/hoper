import 'dart:typed_data';

/// Web 等平台无 [dart:ffi] 时的占位实现。
class RbIceFfi {
  static bool get available => false;
  static bool get agentRunAvailable => false;

  static Object? viewerNew(int timeoutMs) => null;

  static void viewerPush(Object? h, Uint8List data) {}

  static Uint8List? viewerPollOut(Object? h) => null;

  static int viewerState(Object? h) => -1;

  static int viewerWrite(Object? h, int typ, Uint8List payload) => -1;

  static Uint8List? viewerRead(Object? h) => null;

  static void viewerClose(Object? h) {}

  static int agentRun(String signalWs, String room, String root, int timeoutMs) => -1;

  static int agentRunning() => 0;
}
