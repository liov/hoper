import 'dart:async';

import 'package:app/remotebrowse/direct_dialer.dart';
import 'package:app/remotebrowse/ice_dialer.dart';
import 'package:app/remotebrowse/relay_transport.dart';
import 'package:app/remotebrowse/room_dialer.dart';
import 'package:app/remotebrowse/signal_session.dart';
import 'package:app/remotebrowse/wire_codec.dart';
import 'package:app/remotebrowse/wire_session.dart';
import 'package:app/remotebrowse/wire_transport.dart';

class RbViewerSession {
  static const iceTimeout = Duration(seconds: 12);
  static const directPort = 19091;

  static Future<RbWireSession> connect(Uri signalWs, String room, {String? directHost, int? directPort}) async {
    final sig = await RbSignalSession.connect(signalWs);
    try {
      await sig.registerViewer(room);
      final direct = await RbDirectDialer.tryConnect(sig, manualHost: directHost, manualPort: directPort);
      if (direct != null) {
        await sig.close();
        return RbWireSession(direct);
      }
      final roomLink = await RbRoomDialer.tryConnect(sig);
      if (roomLink != null) {
        await sig.close();
        return RbWireSession(roomLink);
      }
      final ice = await _pickIce(sig);
      if (ice != null) {
        await sig.close();
        return RbWireSession(ice);
      }
      final tok = await sig.waitRelayToken();
      await sig.close();
      final relay = await RbRelayTransport.connect(tok.relayHost, tok.relayPort, tok.sessionId, rbRoleViewer);
      return RbWireSession(relay);
    } catch (e) {
      await sig.close();
      rethrow;
    }
  }

  static Future<RbWireSession> connectDirect(String host, int port) async {
    final direct = await RbDirectDialer.connectManual(host, port);
    if (direct == null) {
      throw StateError('direct connect failed');
    }
    return RbWireSession(direct);
  }

  static Future<RbWireTransport?> _pickIce(RbSignalSession sig) async {
    final c = Completer<RbWireTransport?>();
    unawaited(RbIceDialer.tryConnect(sig).then((v) {
      if (!c.isCompleted) {
        c.complete(v);
      }
    }));
    final t = Timer(iceTimeout, () {
      if (!c.isCompleted) {
        c.complete(null);
      }
    });
    final v = await c.future;
    t.cancel();
    return v;
  }
}
