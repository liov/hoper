import 'dart:async';
import 'dart:html' as html;
import 'dart:typed_data';

import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/signal_ws_queue.dart';

/// Web 信令 WS（dart:html），不用 web_socket_channel，避免 sink 竞态。
class RbSignalWsConnection {
  RbSignalWsConnection._(this._ws);

  final html.WebSocket _ws;
  final _frames = StreamController<Uint8List>.broadcast();
  final _ops = RbWsOpQueue();
  StreamSubscription<html.MessageEvent>? _msgSub;
  StreamSubscription<html.Event>? _errSub;
  StreamSubscription<html.CloseEvent>? _closeSub;
  var _closed = false;
  var _userClosed = false;
  var _framesEnded = false;

  static Future<RbSignalWsConnection> connect(Uri url) async {
    final ws = html.WebSocket(url.toString());
    ws.binaryType = 'arraybuffer';
    final done = Completer<void>();
    void onOpen(html.Event _) {
      if (!done.isCompleted) {
        done.complete();
      }
    }

    void onFail(html.Event _) {
      if (!done.isCompleted) {
        done.completeError(StateError('signal ws open failed'));
      }
    }

    ws.onOpen.listen(onOpen);
    ws.onError.listen(onFail);
    await done.future.timeout(const Duration(seconds: 20), onTimeout: () {
      throw StateError('signal ws open timeout');
    });
    final conn = RbSignalWsConnection._(ws);
    conn._msgSub = ws.onMessage.listen(conn._onMessage);
    conn._errSub = ws.onError.listen(conn._onError);
    conn._closeSub = ws.onClose.listen(conn._onClose);
    return conn;
  }

  Stream<Uint8List> get frames => _frames.stream;

  bool get isClosed => _closed;

  static Uint8List? _bytes(Object? data) {
    if (data is ByteBuffer) {
      return data.asUint8List();
    }
    if (data is Uint8List) {
      return data;
    }
    if (data is List<int>) {
      return Uint8List.fromList(data);
    }
    return null;
  }

  void _onMessage(html.MessageEvent e) {
    if (_closed || _frames.isClosed) {
      return;
    }
    final raw = _bytes(e.data);
    if (raw == null) {
      rbLog.warning('signal ws frame ignored type=${e.data.runtimeType}');
      return;
    }
    _frames.add(raw);
  }

  void _onError(html.Event _) {
    rbLog.warning('signal ws error');
    unawaited(_ops.run(() => _endFrames(StateError('signal closed'), StackTrace.current)));
  }

  void _onClose(html.CloseEvent _) {
    if (_closed) {
      return;
    }
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
    if (_closed || _ws.readyState != html.WebSocket.OPEN) {
      return;
    }
    try {
      _ws.send(data);
    } catch (e, st) {
      rbLog.fine('signal ws send skipped: $e', e, st);
      _closed = true;
    }
  }

  Future<void> close() => _ops.run(_closeRaw);

  Future<void> _closeRaw() async {
    if (_closed && _msgSub == null) {
      return;
    }
    _userClosed = true;
    _closed = true;
    await _msgSub?.cancel();
    await _errSub?.cancel();
    await _closeSub?.cancel();
    _msgSub = null;
    _errSub = null;
    _closeSub = null;
    try {
      _ws.close();
    } catch (e, st) {
      rbLog.fine('signal ws close: $e', e, st);
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
