import 'dart:async';
import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_platform.dart';
import 'package:app/remotebrowse/signal_ws.dart';
import 'package:app/remotebrowse/signal_ws_queue.dart';

class RbSignalSession {
  RbSignalSession._(this._ws);

  final RbSignalWsConnection _ws;
  StreamSubscription? _frameSub;
  final _registerWaiters = <Completer<RegisterResp>>[];
  final _ice = StreamController<SignalEnvelope>.broadcast();
  final _peer = StreamController<PeerEndpoints>.broadcast();
  final _relay = StreamController<RelayToken>.broadcast();
  RelayToken? _pendingRelay;
  PeerEndpoints? _pendingPeer;
  var _closed = false;
  Completer<void>? _closeDone;
  Future<void> _sendMux = Future<void>.value();
  StreamSubscription<Uint8List>? _iceInboundSub;

  static const connectTimeout = Duration(seconds: 20);

  static Future<RbSignalSession> connect(Uri url) async {
    if (url.hasPort && (url.port < 1 || url.port > 65535)) {
      throw ArgumentError('Invalid port ${url.port}');
    }
    rbLog.info('signal ws connect $url');
    RbSignalWsConnection ws;
    try {
      ws = await RbSignalWsConnection.connect(url);
    } catch (e) {
      throw StateError('信令 WebSocket 连接失败: $e');
    }
    final sess = RbSignalSession._(ws);
    sess._frameSub = ws.frames.listen(
      sess._onFrame,
      onError: sess._onFrameError,
      onDone: sess._onFrameDone,
      cancelOnError: true,
    );
    rbLog.fine('signal ws connected');
    return sess;
  }

  bool get isClosed => _closed || _ws.isClosed;

  void _onFrame(Uint8List raw) {
    if (_closed) {
      return;
    }
    final SignalEnvelope env;
    try {
      env = SignalEnvelope.fromBuffer(raw);
    } catch (e, st) {
      rbLog.warning('signal decode failed', e, st);
      return;
    }
    if (env.hasRegisterAck() || env.hasError()) {
      if (env.hasError()) {
        rbLog.warning('signal register error: ${env.error}');
      } else {
        rbLog.info('signal register ok peer=${env.registerAck.peerId} room=${env.registerAck.roomId}');
      }
      for (final w in List<Completer<RegisterResp>>.from(_registerWaiters)) {
        if (env.hasError()) {
          w.completeError(StateError(env.error));
        } else {
          w.complete(env.registerAck);
        }
      }
      _registerWaiters.clear();
      return;
    }
    if (env.hasIceParameters() || env.hasIceCandidate() || env.hasIceComplete()) {
      if (!_ice.isClosed) {
        _ice.add(env);
      }
      return;
    }
    if (env.hasPeerEndpoints()) {
      final e = env.peerEndpoints;
      rbLog.fine('signal recv peer_endpoints n=${e.items.length}');
      if (e.items.isNotEmpty) {
        _pendingPeer = e;
      }
      if (!_peer.isClosed) {
        _peer.add(e);
      }
      return;
    }
    if (env.hasRelayToken()) {
      final t = env.relayToken;
      rbLog.info('signal recv relay_token ${t.relayHost}:${t.relayPort} session=${t.sessionId}');
      _pendingRelay = t;
      if (!_relay.isClosed) {
        _relay.add(t);
      }
    }
  }

  void _onFrameError(Object err) {
    if (_closed) {
      return;
    }
    rbLog.warning('signal frame stream error', err);
    _shutdownWaiters(StateError('signal closed'));
  }

  void _onFrameDone() {
    if (_closed) {
      return;
    }
    rbLog.fine('signal frame stream done');
    _shutdownWaiters(StateError('signal closed'));
  }

  void _shutdownWaiters(Object err) {
    _markClosed();
    _failRegisterWaiters(err);
    _closeSideStreams();
  }

  void _failRegisterWaiters(Object err) {
    for (final w in List<Completer<RegisterResp>>.from(_registerWaiters)) {
      if (!w.isCompleted) {
        w.completeError(err);
      }
    }
    _registerWaiters.clear();
  }

  void _markClosed() => _closed = true;

  void _closeSideStreams() {
    _safeCloseCtrl(_ice);
    _safeCloseCtrl(_peer);
    _safeCloseCtrl(_relay);
  }

  void _safeCloseCtrl<T>(StreamController<T> c) {
    if (c.isClosed) {
      return;
    }
    unawaited(c.close().catchError((Object e) {
      if (!rbWsSinkBoundError(e)) {
        rbLog.fine('signal side stream close: $e');
      }
    }));
  }

  Future<void> _awaitCloseDone() async {
    final done = _closeDone;
    if (done != null) {
      await done.future.timeout(const Duration(seconds: 3), onTimeout: () {});
    }
  }

  Future<void> send(SignalEnvelope env) async {
    if (_closed || _ws.isClosed) {
      return;
    }
    final run = _sendMux.then((_) async {
      if (_closed || _ws.isClosed) {
        return;
      }
      await _ws.send(env.writeToBuffer());
    });
    _sendMux = run.then((_) {}, onError: (_) {});
    await run.catchError((_) {});
  }

  Future<RegisterResp> registerAgent(String room) async {
    await _awaitCloseDone();
    if (_closed) {
      throw StateError('signal closed');
    }
    rbLog.info('signal register agent room=$room');
    final c = Completer<RegisterResp>();
    _registerWaiters.add(c);
    await send(SignalEnvelope(
      register: RegisterReq(
        roomCode: room,
        role: 'agent',
        caps: DeviceCapabilities(hasIpv6: rbHasIpv6(), platform: rbPlatformName()),
      ),
    ));
    if (_closed && !c.isCompleted) {
      throw StateError('signal closed');
    }
    return _waitRegister(c, 'agent');
  }

  Future<RegisterResp> registerViewer(String room) async {
    await _awaitCloseDone();
    if (_closed) {
      throw StateError('signal closed');
    }
    rbLog.info('signal register viewer room=$room');
    final c = Completer<RegisterResp>();
    _registerWaiters.add(c);
    await send(SignalEnvelope(
      register: RegisterReq(
        roomCode: room,
        role: 'viewer',
        caps: DeviceCapabilities(hasIpv6: rbHasIpv6(), platform: rbPlatformName()),
      ),
    ));
    if (_closed && !c.isCompleted) {
      throw StateError('signal closed');
    }
    return _waitRegister(c, 'viewer');
  }

  Future<RegisterResp> _waitRegister(Completer<RegisterResp> c, String role) async {
    try {
      return await c.future.timeout(connectTimeout, onTimeout: () {
        throw StateError(_registerTimeoutMsg(role));
      });
    } catch (e) {
      throw StateError(_safeErr(e));
    }
  }

  String _registerTimeoutMsg(String role) {
    if (_closed) {
      return '信令已断开：daemon 可能重启或网络中断，请确认信令地址与 Agent 一致后重试';
    }
    return '信令注册超时（${connectTimeout.inSeconds}s 内未收到 daemon 确认）：'
        '请核对信令地址是否与 Agent 的 RB_SIGNAL_URL 为同一台机器、rb-daemon 是否在跑、房间码是否一致';
  }

  static String _safeErr(Object e) {
    final s = e.toString();
    if (s.contains('StreamSink is bound to a stream')) {
      return 'signal closed';
    }
    if (e is StateError && e.message.isNotEmpty) {
      return e.message;
    }
    return s.replaceFirst('Bad state: ', '');
  }

  Future<void> sendIceParameters(String ufrag, String pwd) async {
    await send(SignalEnvelope(iceParameters: IceParameters(ufrag: ufrag, pwd: pwd)));
  }

  Future<void> sendIceCandidate(IceCandidateInit cand) async {
    await send(SignalEnvelope(iceCandidate: cand));
  }

  Future<void> sendIceComplete() async {
    await send(SignalEnvelope(iceComplete: true));
  }

  Future<SignalEnvelope> recvIce() async {
    try {
      return await _ice.stream.first;
    } catch (e) {
      throw StateError(_safeErr(e));
    }
  }

  StreamSubscription<Uint8List> bindIceInbound(void Function(Uint8List data) push) {
    unawaited(_iceInboundSub?.cancel());
    _iceInboundSub = _ice.stream.map((e) => Uint8List.fromList(e.writeToBuffer())).listen(
      push,
      onError: (_) {},
    );
    return _iceInboundSub!;
  }

  /// ICE 选路结束或放弃后调用，避免继续往信令 WS 发 ICE 帧导致 close/send 竞态。
  Future<void> detachIceInbound() async {
    final sub = _iceInboundSub;
    _iceInboundSub = null;
    await sub?.cancel();
  }

  Future<void> sendPeerEndpoints(PeerEndpoints eps) async {
    rbLog.fine('signal send peer_endpoints n=${eps.items.length}');
    await send(SignalEnvelope(peerEndpoints: eps));
  }

  Future<PeerEndpoints> waitPeerEndpoints({Duration timeout = const Duration(seconds: 30)}) async {
    final cached = _pendingPeer;
    if (cached != null && cached.items.isNotEmpty) {
      _pendingPeer = null;
      return cached;
    }
    try {
      return await _peer.stream.timeout(timeout).firstWhere((e) => e.items.isNotEmpty);
    } catch (_) {
      throw StateError('signal closed');
    }
  }

  Future<RelayToken> waitRelayToken({Duration timeout = const Duration(seconds: 120)}) async {
    final cached = _pendingRelay;
    if (cached != null) {
      _pendingRelay = null;
      return cached;
    }
    try {
      return await _relay.stream.timeout(timeout).first;
    } catch (_) {
      throw StateError('signal closed');
    }
  }

  Future<void> close() async {
    if (_closeDone != null) {
      return _closeDone!.future.timeout(const Duration(seconds: 3), onTimeout: () {});
    }
    final done = Completer<void>();
    _closeDone = done;
    _markClosed();
    _failRegisterWaiters(StateError('signal closed'));
    try {
      await detachIceInbound();
      await _frameSub?.cancel();
      _frameSub = null;
      _closeSideStreams();
      await _sendMux.catchError((_) {});
      await _ws.close();
      done.complete();
    } catch (e, st) {
      rbLog.fine('signal close: $e', e, st);
      if (!done.isCompleted) {
        done.complete();
      }
    }
  }
}
