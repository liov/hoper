import 'dart:async';

import 'package:app/remotebrowse/direct_dialer.dart';
import 'package:app/remotebrowse/signal_session.dart';
import 'package:app/remotebrowse/wire_transport.dart';

class RbRoomDialer {
  static const roomTimeout = Duration(seconds: 5);
  static const quicPort = 19092;

  static Future<RbWireTransport?> tryConnect(RbSignalSession sig) async {
    try {
      final eps = await sig.waitPeerEndpoints(timeout: roomTimeout);
      for (final ep in eps.items) {
        if (ep.host.isEmpty) {
          continue;
        }
        final t = await RbDirectDialer.connectManual(ep.host, quicPort);
        if (t != null) {
          return t;
        }
      }
    } catch (_) {}
    return null;
  }
}
