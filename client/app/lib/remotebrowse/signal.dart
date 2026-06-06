import 'dart:async';

import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class RbSignalClient {
  RbSignalClient(this._channel);

  final WebSocketChannel _channel;
  final _inbox = StreamController<SignalEnvelope>.broadcast();

  static Future<RbSignalClient> connect(Uri url) async {
    final ch = WebSocketChannel.connect(url);
    final cli = RbSignalClient(ch);
    cli._channel.stream.listen(
      (data) {
        if (data is! List<int>) {
          return;
        }
        cli._inbox.add(SignalEnvelope.fromBuffer(data));
      },
      onError: cli._inbox.addError,
      onDone: () => cli._inbox.close(),
    );
    return cli;
  }

  Future<void> close() async {
    await _channel.sink.close();
    await _inbox.close();
  }

  Future<void> send(SignalEnvelope env) async {
    _channel.sink.add(env.writeToBuffer());
  }

  Future<RegisterResp> registerViewer(String room) async {
    await send(SignalEnvelope(register: RegisterReq(roomCode: room, role: 'viewer')));
    return _waitAck();
  }

  Future<RelayToken> waitRelayToken({Duration timeout = const Duration(seconds: 120)}) async {
    final tok = await _inbox.stream.timeout(timeout).firstWhere((e) => e.hasRelayToken() || e.hasError());
    if (tok.hasError()) {
      throw StateError(tok.error);
    }
    return tok.relayToken;
  }

  Future<RegisterResp> _waitAck() async {
    final ack = await _inbox.stream.firstWhere((e) => e.hasRegisterAck() || e.hasError());
    if (ack.hasError()) {
      throw StateError(ack.error);
    }
    return ack.registerAck;
  }
}
