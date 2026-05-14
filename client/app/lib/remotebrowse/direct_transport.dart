import 'dart:async';
import 'dart:io';
import 'dart:typed_data';

import 'package:app/remotebrowse/wire_codec.dart';
import 'package:app/remotebrowse/wire_transport.dart';

class RbDirectTransport implements RbWireTransport {
  RbDirectTransport._(this._sock, this._buf);

  final Socket _sock;
  final _SockBuf _buf;

  static Future<RbDirectTransport> connect(String host, int port) async {
    final sock = await Socket.connect(host, port);
    return RbDirectTransport.fromSocket(sock);
  }

  factory RbDirectTransport.fromSocket(Socket sock) => RbDirectTransport._(sock, _SockBuf(sock));

  @override
  Future<void> writeFrame(int typ, Uint8List payload) async {
    _sock.add(rbEncodeWireFrame(typ, payload));
    await _sock.flush();
  }

  @override
  Future<(int, Uint8List)> readFrame() async {
    final hdr = await _buf.take(6);
    final n = ByteData.sublistView(hdr).getUint32(2, Endian.big);
    final body = await _buf.take(n);
    return rbDecodeWireFrame(Uint8List.fromList([...hdr, ...body]));
  }

  @override
  Future<void> close() async {
    await _buf.dispose();
    await _sock.close();
  }
}

class _SockBuf {
  _SockBuf(Socket s) {
    _sub = s.listen((chunk) {
      _pending.addAll(chunk);
      for (final w in _waiters) {
        if (_pending.length >= w.$1) {
          w.$2.complete();
        }
      }
    });
  }

  final _pending = <int>[];
  final _waiters = <(int, Completer<void>)>[];
  late final StreamSubscription<Uint8List> _sub;

  Future<void> dispose() async {
    await _sub.cancel();
  }

  Future<Uint8List> take(int n) async {
    while (_pending.length < n) {
      final c = Completer<void>();
      _waiters.add((n, c));
      await c.future;
      _waiters.removeWhere((w) => w.$2 == c);
    }
    final out = Uint8List.fromList(_pending.sublist(0, n));
    _pending.removeRange(0, n);
    return out;
  }
}
