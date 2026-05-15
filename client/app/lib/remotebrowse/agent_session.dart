import 'dart:async';
import 'dart:io';

import 'package:app/remotebrowse/agent_wire.dart';
import 'package:app/remotebrowse/direct_dialer.dart';
import 'package:app/remotebrowse/direct_transport.dart';
import 'package:app/remotebrowse/ice_dialer.dart';
import 'package:app/remotebrowse/relay_transport.dart';
import 'package:app/remotebrowse/room_dialer.dart';
import 'package:app/remotebrowse/signal_session.dart';
import 'package:app/remotebrowse/wire_codec.dart';
import 'package:app/remotebrowse/wire_transport.dart';

class RbAgentSession {
  static const connectTimeout = Duration(seconds: 12);

  static Future<void> run(Uri signalWs, String room, String root) async {
    final sig = await RbSignalSession.connect(signalWs);
    ServerSocket? listener;
    try {
      listener = await ServerSocket.bind(InternetAddress.anyIPv4, RbDirectDialer.defaultPort);
      await sig.registerAgent(room);
      await sig.sendPeerEndpoints(await RbDirectDialer.gatherEndpoints(listener.port));
      final transport = await _pickTransport(sig, listener);
      await sig.close();
      await RbAgentWire.serve(transport, root);
    } finally {
      await listener?.close();
      await sig.close();
    }
  }

  static Future<RbWireTransport> _pickTransport(RbSignalSession sig, ServerSocket listener) async {
    final direct = await _pickDirect(sig, listener);
    if (direct != null) {
      return direct;
    }
    final roomLink = await RbRoomDialer.tryConnect(sig);
    if (roomLink != null) {
      return roomLink;
    }
    final ice = await _pickIce(sig);
    if (ice != null) {
      return ice;
    }
    final tok = await sig.waitRelayToken();
    return RbRelayTransport.connect(tok.relayHost, tok.relayPort, tok.sessionId, rbRoleAgent);
  }

  static Future<RbWireTransport?> _pickDirect(RbSignalSession sig, ServerSocket listener) async {
    final c = Completer<RbWireTransport?>();
    unawaited(_acceptOnce(listener, c));
    unawaited(_dialPeer(sig, c));
    try {
      return await c.future.timeout(RbDirectDialer.directTimeout);
    } catch (_) {
      return null;
    }
  }

  static Future<void> _acceptOnce(ServerSocket listener, Completer<RbWireTransport?> c) async {
    try {
      final sock = await listener.first.timeout(RbDirectDialer.directTimeout);
      if (!c.isCompleted) {
        c.complete(RbDirectTransport.fromSocket(sock));
      }
    } catch (_) {}
  }

  static Future<void> _dialPeer(RbSignalSession sig, Completer<RbWireTransport?> c) async {
    try {
      final eps = await sig.waitPeerEndpoints(timeout: RbDirectDialer.directTimeout);
      for (final ep in eps.items) {
        if (c.isCompleted || ep.host.isEmpty || ep.port == 0) {
          continue;
        }
        final t = await RbDirectDialer.connectManual(ep.host, ep.port);
        if (t != null && !c.isCompleted) {
          c.complete(t);
          return;
        }
      }
    } catch (_) {}
  }

  static Future<RbWireTransport?> _pickIce(RbSignalSession sig) async {
    final c = Completer<RbWireTransport?>();
    unawaited(RbIceDialer.tryConnect(sig).then((v) {
      if (!c.isCompleted) {
        c.complete(v);
      }
    }));
    final t = Timer(connectTimeout, () {
      if (!c.isCompleted) {
        c.complete(null);
      }
    });
    final v = await c.future;
    t.cancel();
    return v;
  }
}
