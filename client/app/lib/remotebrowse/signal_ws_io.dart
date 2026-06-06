import 'dart:async';
import 'dart:io';
import 'dart:typed_data';

import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/signal_ws_queue.dart';

/// 信令二进制 WebSocket（dart:io）；send/close 串行，避免 sink 竞态。
class RbSignalWsConnection {
  RbSignalWsConnection._(this._ws);

  final WebSocket _ws;
  final _frames = StreamController<Uint8List>.broadcast();
  final _ops = RbWsOpQueue();
  StreamSubscription? _sub;
  var _closed = false;
  var _userClosed = false;
  var _framesEnded = false;

  static Future<RbSignalWsConnection> connect(Uri url) async {
    final ws = await WebSocket.connect(
      url.toString(),
      compression: CompressionOptions.compressionOff,
    );
    ws.pingInterval = const Duration(seconds: 25);
    final conn = RbSignalWsConnection._(ws);
    conn._sub = ws.listen(
      conn._onMessage,
      onError: conn._onError,
      onDone: conn._onDone,
      cancelOnError: true,
    );
    return conn;
  }

  Stream<Uint8List> get frames => _frames.stream;

  bool get isClosed => _closed;

  static Uint8List? _bytes(dynamic data) {
    if (data is Uint8List) {
      return data;
    }
    if (data is List<int>) {
      return Uint8List.fromList(data);
    }
    return null;
  }

  void _onMessage(dynamic data) {
    if (_closed || _frames.isClosed) {
      return;
    }
    final raw = _bytes(data);
    if (raw == null) {
      rbLog.warning('signal ws frame ignored type=${data.runtimeType}');
      return;
    }
    _frames.add(raw);
  }

  void _onError(Object err, StackTrace st) {
    rbLog.warning('signal ws error', err, st);
    unawaited(_ops.run(() => _endFrames(err, st)));
  }

  void _onDone() {
    if (_closed) {
      return;
    }
    rbLog.fine('signal ws done');
    unawaited(_ops.run(() async {
      _closed = true;
      if (!_userClosed) {
        await _endFrames(StateError('signal closed'), StackTrace.current);
      }
    }));
  }

  Future<void> _endFrames(Object err, StackTrace st) async {
    if (_framesEnded || _frames.isClosed) {
      return;
    }
    _framesEnded = true;
    _closed = true;
    _frames.addError(err, st);
    try {
      await _frames.close();
    } catch (e) {
      if (!rbWsSinkBoundError(e)) {
        rbLog.fine('signal frames close: $e');
      }
    }
  }

  Future<void> send(Uint8List data) => _ops.run(() => _sendRaw(data));

  Future<void> _sendRaw(Uint8List data) async {
    if (_closed) {
      return;
    }
    try {
      _ws.add(data);
    } catch (e, st) {
      if (rbWsSinkBoundError(e)) {
        rbLog.fine('signal ws send skipped (sink bound)');
      } else {
        rbLog.fine('signal ws send skipped: $e', e, st);
      }
      _closed = true;
    }
  }

  Future<void> close() => _ops.run(_closeRaw);

  Future<void> _closeRaw() async {
    if (_closed && _sub == null) {
      return;
    }
    _userClosed = true;
    _closed = true;
    final sub = _sub;
    _sub = null;
    try {
      await sub?.cancel();
    } catch (e, st) {
      rbLog.fine('signal ws sub cancel: $e', e, st);
    }
    try {
      await _ws.close(WebSocketStatus.goingAway);
    } catch (e, st) {
      if (!rbWsSinkBoundError(e)) {
        rbLog.fine('signal ws close: $e', e, st);
      }
    }
    if (!_framesEnded && !_frames.isClosed) {
      _framesEnded = true;
      try {
        await _frames.close();
      } catch (e) {
        if (!rbWsSinkBoundError(e)) {
          rbLog.fine('signal frames close: $e');
        }
      }
    }
  }
}
