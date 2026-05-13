import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:app/gen/pb/remotebrowse/remotebrowse/browse.service.pb.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/signal.dart';

const _wireVer = 1;
const _typeFileIndex = 2;
const _typeThumbReq = 3;
const _typeThumbData = 4;
const _roleViewer = 0;

class RbRelaySession {
  RbRelaySession._(this._sock, this._buf);

  final Socket _sock;
  final _SockBuf _buf;

  static Future<RbRelaySession> openViewer(Uri signalWs, String room) async {
    final sig = await RbSignalClient.connect(signalWs);
    await sig.registerViewer(room);
    final tok = await sig.waitRelayToken();
    await sig.close();
    final sock = await Socket.connect(tok.relayHost, tok.relayPort);
    final buf = _SockBuf(sock);
    await _writeJoin(sock, tok.sessionId, _roleViewer);
    return RbRelaySession._(sock, buf);
  }

  Future<void> close() async {
    await _buf.dispose();
    await _sock.close();
  }

  Future<List<RbFileEntry>> listFiles(String root) async {
    final req = ListFilesRequest(rootPath: root);
    await _writeWire(_sock, _buf, _typeFileIndex, req.writeToBuffer());
    final frame = await _readWire(_buf);
    if (frame.$1 != _typeFileIndex) {
      throw StateError('unexpected wire type ${frame.$1}');
    }
    final resp = ListFilesResponse.fromBuffer(frame.$2);
    return resp.entries.map((e) => RbFileEntry(id: e.id, name: e.name, size: e.size.toInt(), thumbHash: e.thumbHash)).toList();
  }

  Future<Uint8List> fetchThumb(String path, {int maxEdge = 256}) async {
    final req = ThumbnailRequest(path: path, maxEdge: maxEdge);
    await _writeWire(_sock, _buf, _typeThumbReq, req.writeToBuffer());
    final frame = await _readWire(_buf);
    if (frame.$1 != _typeThumbData) {
      throw StateError('unexpected wire type ${frame.$1}');
    }
    return Uint8List.fromList(ThumbnailResponse.fromBuffer(frame.$2).data);
  }

  static Future<void> _writeJoin(Socket s, String sessionId, int role) async {
    final id = _uuidBytes(sessionId);
    final buf = BytesBuilder();
    buf.add(utf8.encode('RBRL'));
    buf.add([1, ...id, role]);
    s.add(buf.toBytes());
    await s.flush();
  }

  static List<int> _uuidBytes(String sessionId) {
    final hex = sessionId.replaceAll('-', '');
    return List.generate(16, (i) => int.parse(hex.substring(i * 2, i * 2 + 2), radix: 16));
  }

  static Future<void> _writeWire(Socket s, _SockBuf buf, int typ, Uint8List payload) async {
    final hdr = ByteData(6);
    hdr.setUint8(0, _wireVer);
    hdr.setUint8(1, typ);
    hdr.setUint32(2, payload.length, Endian.big);
    final frame = BytesBuilder();
    frame.add(hdr.buffer.asUint8List());
    frame.add(payload);
    await _relayWrite(s, frame.toBytes());
  }

  static Future<(int, Uint8List)> _readWire(_SockBuf buf) async {
    final raw = await _relayRead(buf);
    final hdr = ByteData.sublistView(raw, 0, 6);
    if (hdr.getUint8(0) != _wireVer) {
      throw StateError('bad wire version');
    }
    final n = hdr.getUint32(2, Endian.big);
    return (hdr.getUint8(1), Uint8List.sublistView(raw, 6, 6 + n));
  }

  static Future<void> _relayWrite(Socket s, Uint8List payload) async {
    final sz = ByteData(4)..setUint32(0, payload.length, Endian.big);
    s.add(sz.buffer.asUint8List());
    s.add(payload);
    await s.flush();
  }

  static Future<Uint8List> _relayRead(_SockBuf buf) async {
    final sz = await buf.take(4);
    final n = ByteData.sublistView(sz).getUint32(0, Endian.big);
    return buf.take(n);
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
