import 'dart:async';
import 'dart:io';

import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:app/remotebrowse/direct_transport.dart';
import 'package:app/remotebrowse/signal_session.dart';
import 'package:app/remotebrowse/wire_transport.dart';

class RbDirectDialer {
  static const directTimeout = Duration(seconds: 5);
  static const defaultPort = 19091;

  static Future<RbWireTransport?> tryConnect(RbSignalSession sig, {String? manualHost, int? manualPort}) async {
    final ln = await _listen();
    try {
      await sig.sendPeerEndpoints(await _gatherEndpoints(ln?.$2 ?? defaultPort));
      final c = Completer<RbWireTransport?>();
      if (manualHost != null && manualHost.isNotEmpty) {
        unawaited(_dialHost(manualHost, manualPort ?? defaultPort, c));
      }
      if (ln != null) {
        unawaited(_acceptOnce(ln.$1, c));
      }
      unawaited(_dialPeer(sig, c));
      try {
        return await c.future.timeout(directTimeout);
      } catch (_) {
        return null;
      }
    } finally {
      await ln?.$1.close();
    }
  }

  static Future<RbWireTransport?> connectManual(String host, int port) async {
    try {
      return await RbDirectTransport.connect(host, port).timeout(directTimeout);
    } catch (_) {
      return null;
    }
  }

  static Future<(ServerSocket, int)?> _listen() async {
    try {
      final ln = await ServerSocket.bind(InternetAddress.anyIPv4, defaultPort);
      return (ln, ln.port);
    } catch (_) {
      return null;
    }
  }

  static Future<PeerEndpoints> _gatherEndpoints(int port) async {
    final items = <PeerEndpoint>[];
    final ifaces = await NetworkInterface.list(type: InternetAddressType.IPv4, includeLoopback: false);
    for (final ni in ifaces) {
      for (final addr in ni.addresses) {
        if (addr.isLoopback) {
          continue;
        }
        items.add(PeerEndpoint(host: addr.address, port: port));
      }
    }
    return PeerEndpoints(items: items);
  }

  static Future<void> _dialHost(String host, int port, Completer<RbWireTransport?> c) async {
    final t = await connectManual(host, port);
    if (t != null && !c.isCompleted) {
      c.complete(t);
    }
  }

  static Future<void> _acceptOnce(ServerSocket ln, Completer<RbWireTransport?> c) async {
    try {
      final sock = await ln.first.timeout(directTimeout);
      if (!c.isCompleted) {
        c.complete(RbDirectTransport.fromSocket(sock));
      }
    } catch (_) {}
  }

  static Future<void> _dialPeer(RbSignalSession sig, Completer<RbWireTransport?> c) async {
    try {
      final eps = await sig.waitPeerEndpoints(timeout: directTimeout);
      for (final ep in eps.items) {
        if (c.isCompleted || ep.host.isEmpty || ep.port == 0) {
          continue;
        }
        final t = await connectManual(ep.host, ep.port);
        if (t != null && !c.isCompleted) {
          c.complete(t);
          return;
        }
      }
    } catch (_) {}
  }
}
