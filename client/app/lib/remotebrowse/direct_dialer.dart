import 'dart:async';
import 'dart:io';

import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:app/remotebrowse/connect_cancel.dart';
import 'package:app/remotebrowse/rb_connect_endpoint.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/signal_session.dart';

class RbDirectDialer {
  static const directTimeout = Duration(seconds: 20);
  static const defaultPort = 19091;

  static Future<Socket?> tryConnect(RbSignalSession sig, {String? manualHost, int? manualPort, RbConnectCancel? cancel}) async {
    final ln = await _listen();
    if (ln != null) {
      rbLog.fine('direct listen port=${ln.$2}');
      final eps = await gatherEndpoints(ln.$2);
      rbLog.fine('direct peer_endpoints n=${eps.items.length} ${eps.items.map((e) => "${e.host}:${e.port}").join(", ")}');
      await sig.sendPeerEndpoints(eps);
    } else {
      rbLog.fine('direct listen unavailable, outbound dial agent only (no peer_endpoints)');
    }
    try {
      final c = Completer<Socket?>();
      if (manualHost != null && manualHost.isNotEmpty) {
        unawaited(_dialHost(manualHost, manualPort ?? defaultPort, c));
      }
      if (ln != null) {
        unawaited(_acceptOnce(ln.$1, c));
      }
      unawaited(_dialPeer(sig, c, cancel));
      try {
        return await c.future.timeout(directTimeout);
      } on RbConnectCancelled {
        rethrow;
      } catch (_) {
        cancel?.check();
        rbLog.fine('direct pick timeout');
        return null;
      }
    } finally {
      await ln?.$1.close();
    }
  }

  static Future<Socket?> connectManual(String host, int port) async {
    try {
      final h = rbSocketConnectHost(host);
      rbLog.fine('direct dial $h:$port');
      return await Socket.connect(h, port, timeout: directTimeout);
    } catch (e) {
      rbLog.fine('direct dial failed $host:$port $e');
      return null;
    }
  }

  static Future<(ServerSocket, int)?> _listen() async {
    try {
      final ln = await ServerSocket.bind(InternetAddress.anyIPv4, defaultPort);
      return (ln, ln.port);
    } catch (_) {
      rbLog.fine('direct bind :$defaultPort failed, try ephemeral port');
    }
    try {
      final ln = await ServerSocket.bind(InternetAddress.anyIPv4, 0);
      return (ln, ln.port);
    } catch (e) {
      rbLog.fine('direct listen failed: $e');
      return null;
    }
  }

  static Future<PeerEndpoints> gatherEndpoints(int port) async {
    final items = <PeerEndpoint>[];
    for (final type in [InternetAddressType.IPv4, InternetAddressType.IPv6]) {
      final ifaces = await NetworkInterface.list(type: type, includeLoopback: false);
      for (final ni in ifaces) {
        for (final addr in ni.addresses) {
          if (addr.isLoopback) {
            continue;
          }
          final host = addr.type == InternetAddressType.IPv6 ? '[${addr.address}]' : addr.address;
          items.add(PeerEndpoint(host: host, port: port));
        }
      }
    }
    return PeerEndpoints(items: items);
  }

  static Future<void> _dialHost(String host, int port, Completer<Socket?> c) async {
    final t = await connectManual(host, port);
    if (t != null && !c.isCompleted) {
      c.complete(t);
    }
  }

  static Future<void> _acceptOnce(ServerSocket ln, Completer<Socket?> c) async {
    try {
      final sock = await ln.first.timeout(directTimeout);
      if (!c.isCompleted) {
        rbLog.info('direct inbound accept ${sock.remoteAddress}:${sock.remotePort}');
        c.complete(sock);
      }
    } catch (_) {}
  }

  static Future<void> _dialPeer(RbSignalSession sig, Completer<Socket?> c, RbConnectCancel? cancel) async {
    try {
      final eps = await sig.waitPeerEndpoints(timeout: directTimeout);
      cancel?.check();
      final dials = <Future<void>>[];
      for (final ep in eps.items) {
        if (ep.host.isEmpty || ep.port == 0) {
          continue;
        }
        dials.add(_dialOne(ep.host, ep.port, c));
      }
      if (dials.isNotEmpty) {
        await Future.wait(dials);
      }
    } catch (_) {}
  }

  static Future<void> _dialOne(String host, int port, Completer<Socket?> c) async {
    final t = await connectManual(host, port);
    if (t != null && !c.isCompleted) {
      c.complete(t);
    }
  }
}
