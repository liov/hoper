import 'dart:async';
import 'dart:typed_data';
import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:flutter/foundation.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class RbSignalSession {
  RbSignalSession._(this._channel);

  final WebSocketChannel _channel;
  final _registerWaiters = <Completer<RegisterResp>>[];
  final _ice = StreamController<SignalEnvelope>.broadcast();
  final _peer = StreamController<PeerEndpoints>.broadcast();
  final _relay = StreamController<RelayToken>.broadcast();
  final _errs = StreamController<Object>.broadcast();
  var _closed = false;

  static Future<RbSignalSession> connect(Uri url) async {
    final ch = WebSocketChannel.connect(url);
    final sess = RbSignalSession._(ch);
    ch.stream.listen(sess._onData, onError: sess._onError, onDone: sess._onDone);
    return sess;
  }

  void _onData(Object? data) {
    if (data is! List<int>) {
      return;
    }
    final env = SignalEnvelope.fromBuffer(data);
    if (env.hasRegisterAck() || env.hasError()) {
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
      _ice.add(env);
      return;
    }
    if (env.hasPeerEndpoints()) {
      _peer.add(env.peerEndpoints);
      return;
    }
    if (env.hasRelayToken()) {
      _relay.add(env.relayToken);
    }
  }

  void _onError(Object err) {
    _errs.add(err);
    _ice.addError(err);
    _relay.addError(err);
  }

  void _onDone() {
    if (_closed) {
      return;
    }
    _closed = true;
    final err = StateError('signal closed');
    _errs.add(err);
    _ice.addError(err);
    _relay.addError(err);
  }

  Future<void> send(SignalEnvelope env) async {
    _channel.sink.add(env.writeToBuffer());
  }

  Future<RegisterResp> registerViewer(String room) async {
    final c = Completer<RegisterResp>();
    _registerWaiters.add(c);
    await send(SignalEnvelope(
      register: RegisterReq(
        roomCode: room,
        role: 'viewer',
        caps: DeviceCapabilities(hasIpv6: !kIsWeb, platform: kIsWeb ? 'web' : 'flutter'),
      ),
    ));
    return c.future;
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
    return _ice.stream.first;
  }

  StreamSubscription<Uint8List> bindIceInbound(void Function(Uint8List data) push) {
    return _ice.stream.map((e) => Uint8List.fromList(e.writeToBuffer())).listen(push);
  }

  Future<void> sendPeerEndpoints(PeerEndpoints eps) async {
    await send(SignalEnvelope(peerEndpoints: eps));
  }

  Future<PeerEndpoints> waitPeerEndpoints({Duration timeout = const Duration(seconds: 30)}) async {
    return _peer.stream.timeout(timeout).firstWhere((e) => e.items.isNotEmpty);
  }

  Future<RelayToken> waitRelayToken({Duration timeout = const Duration(seconds: 120)}) async {
    return _relay.stream.timeout(timeout).first;
  }

  Future<void> close() async {
    _closed = true;
    await _channel.sink.close();
    await _ice.close();
    await _peer.close();
    await _relay.close();
    await _errs.close();
  }

}
