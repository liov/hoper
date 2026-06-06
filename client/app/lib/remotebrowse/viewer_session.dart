import 'dart:async';
import 'dart:io';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/connect_cancel.dart';
import 'package:app/remotebrowse/direct_dialer.dart';
import 'package:app/remotebrowse/ice_dialer.dart';
import 'package:app/remotebrowse/link_kind.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_connect_endpoint.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/rb_h2_transport.dart';
import 'package:app/remotebrowse/rb_ice_grpc.dart';
import 'package:app/remotebrowse/rb_user_message.dart';
import 'package:app/remotebrowse/signal_session.dart';
import 'package:app/remotebrowse/wire_codec.dart';

class RbViewerSession {
  static const iceTimeout = Duration(seconds: 12);
  static const connectTimeout = Duration(seconds: 45);

  static Future<RbGrpcSession> connect(
    Uri signalWs,
    String room, {
    String? directHost,
    int? directPort,
    RbConnectCancel? cancel,
  }) async {
    rbLog.info(
      'viewer connect signal=$signalWs room=$room direct=$directHost:$directPort',
    );
    await probeSignalHealth(signalWs);
    rbLog.fine('signal health ok ${signalHealthUri(signalWs)}');
    final sig = await RbSignalSession.connect(signalWs);
    cancel?.bind(sig);
    cancel?.check();
    try {
      await sig.registerViewer(room).timeout(connectTimeout, onTimeout: () {
        throw StateError('连接超时');
      });
      cancel?.check();
      rbLog.fine('viewer try direct tcp');
      var direct = await RbDirectDialer.tryConnect(
        sig,
        manualHost: directHost,
        manualPort: directPort,
        cancel: cancel,
      );
      cancel?.check();
      if (direct == null && directHost == null) {
        rbLog.fine('viewer direct retry (wait agent signal)');
        await Future<void>.delayed(const Duration(milliseconds: 800));
        cancel?.check();
        direct = await RbDirectDialer.tryConnect(
          sig,
          manualHost: directHost,
          manualPort: directPort,
          cancel: cancel,
        );
        cancel?.check();
      }
      final agentMedia = await _pickAgentMediaEndpoint(sig, cancel);
      if (direct != null) {
        rbLog.info('viewer link direct tcp http2');
        await _dropSignal(sig, cancel);
        return RbGrpcSession.tcp(RbBoundH2Socket(direct), RbLinkKind.directTcp, room: room);
      }
      if (rbIceFfiAvailable) {
        rbLog.fine('viewer try ice');
        final ice = await _pickIceGrpc(sig, cancel: cancel);
        await sig.detachIceInbound();
        cancel?.check();
        if (ice != null) {
          rbLog.info('viewer link ice http2');
          await _dropSignal(sig, cancel);
          return RbGrpcSession.ice(
            ice,
            room: room,
            agentMediaHost: agentMedia.$1,
            agentMediaPort: agentMedia.$2,
          );
        }
      } else {
        rbLog.fine('viewer skip ice (需 staticLibs/<平台>/librb.a，见 build_flutter_lib.sh)');
      }
      rbLog.fine('viewer try relay');
      final tok = await sig.waitRelayToken();
      cancel?.check();
      rbLog.info('viewer link relay ${tok.relayHost}:${tok.relayPort}');
      await _dropSignal(sig, cancel);
      final sock = await RbRelaySocket.connect(tok.relayHost, tok.relayPort, tok.sessionId, rbRoleViewer);
      cancel?.check();
      return RbGrpcSession.tcp(
        RbBoundH2Socket(sock),
        RbLinkKind.relayTcp,
        room: room,
        agentMediaHost: agentMedia.$1,
        agentMediaPort: agentMedia.$2,
      );
    } on RbConnectCancelled {
      unawaited(_closeSignalQuiet(sig));
      rethrow;
    } catch (e, st) {
      rbLog.warning('viewer connect failed', e, st);
      unawaited(_closeSignalQuiet(sig));
      throw StateError(rbUserMessage(e));
    }
  }

  static const directConnectTimeout = Duration(seconds: 20);

  static Future<RbGrpcSession> connectDirect(String host, int port, {String room = ''}) async {
    final h = rbSocketConnectHost(host);
    rbLog.info('viewer direct dial $h:$port');
    final sock = await Socket.connect(h, port, timeout: directConnectTimeout);
    rbLog.info('viewer direct connected $h:$port');
    return RbGrpcSession.tcp(RbBoundH2Socket(sock), RbLinkKind.directTcp, room: room);
  }

  /// Agent 广播的直连端点（中继/ICE 时旁路 TCP 拉 `/rb/v1/media`）。
  static Future<(String?, int?)> _pickAgentMediaEndpoint(RbSignalSession sig, RbConnectCancel? cancel) async {
    try {
      final eps = await sig.waitPeerEndpoints(timeout: const Duration(seconds: 18));
      cancel?.check();
      for (final ep in eps.items) {
        final host = ep.host.trim();
        final port = ep.port.toInt();
        if (host.isNotEmpty && port > 0) {
          rbLog.fine('agent media endpoint $host:$port');
          return (host, port);
        }
      }
    } catch (_) {}
    return (null, null);
  }

  static Future<void> _closeSignalQuiet(RbSignalSession sig) async {
    try {
      await sig.close().timeout(const Duration(seconds: 2));
    } catch (_) {}
  }

  static Future<void> _dropSignal(RbSignalSession sig, RbConnectCancel? cancel) async {
    cancel?.detachSignal();
    await _closeSignalQuiet(sig);
  }

  static Future<RbIceGrpcBridge?> _pickIceGrpc(RbSignalSession sig, {RbConnectCancel? cancel}) async {
    final c = Completer<RbIceGrpcBridge?>();
    unawaited(
      RbIceViewerDialer.tryConnectGrpc(sig, cancel: cancel).then((v) {
        if (!c.isCompleted) {
          c.complete(v);
        }
      }).catchError((Object e) {
        if (e is RbConnectCancelled && !c.isCompleted) {
          c.completeError(e);
        }
      }),
    );
    final t = Timer(iceTimeout, () {
      if (!c.isCompleted) {
        c.complete(null);
      }
    });
    try {
      final v = await c.future;
      cancel?.check();
      return v;
    } on RbConnectCancelled {
      if (!c.isCompleted) {
        c.complete(null);
      }
      rethrow;
    } finally {
      t.cancel();
    }
  }
}

/// 中继：RBRL 握手后裸 TCP，供 HTTP/2 使用。
class RbRelaySocket {
  static Future<Socket> connect(String host, int port, String sessionId, int role) async {
    rbLog.info('relay dial $host:$port session=$sessionId role=$role');
    final sock = await Socket.connect(host, port);
    sock.add(rbRelayJoinBytes(sessionId, role));
    await sock.flush();
    rbLog.info('relay joined $host:$port');
    return sock;
  }
}
