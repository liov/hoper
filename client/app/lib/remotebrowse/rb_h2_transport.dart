import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'dart:typed_data';

import 'package:app/remotebrowse/rb_zstd_codec.dart';
import 'package:http2/transport.dart';

/// 远程浏览 HTTP/2 API：`Content-Type: application/protobuf`，载荷可选 zstd。
abstract final class RbHttpPaths {
  static const health = '/rb/health';
  static const list = '/rb/v1/list';
  static const listMedia = '/rb/v1/list/media';
  static const thumb = '/rb/v1/thumb';
  static const thumbBatch = '/rb/v1/thumb/batch';
  static const read = '/rb/v1/read';
  static const delete = '/rb/v1/delete';
}

const rbProtoContentType = 'application/protobuf';

/// 已连接 TCP 上的单路 HTTP/2 会话（P2P 直连/中继）。
class RbBoundH2Socket {
  RbBoundH2Socket(this._sock) {
    unawaited(_sock.done.whenComplete(_markDead));
  }

  final Socket _sock;
  Socket get socket => _sock;
  ClientTransportConnection? _conn;
  Future<ClientTransportConnection>? _connecting;
  var _dead = false;

  bool get broken => _dead;

  void _markDead() {
    _dead = true;
    _conn = null;
    _connecting = null;
  }

  Future<ClientTransportConnection> _connection() async {
    if (_dead) {
      throw const SocketException('P2P HTTP/2 连接已关闭，请重新连接远程相册');
    }
    final open = _conn;
    if (open != null && open.isOpen) {
      return open;
    }
    if (_conn != null) {
      for (var i = 0; i < 8; i++) {
        await Future<void>.delayed(const Duration(milliseconds: 25));
        final c = _conn;
        if (c != null && c.isOpen) {
          return c;
        }
      }
      _markDead();
      throw const SocketException('P2P HTTP/2 连接已关闭，请重新连接远程相册');
    }
    final inflight = _connecting;
    if (inflight != null) {
      return inflight;
    }
    final job = _connectOnce();
    _connecting = job;
    try {
      return await job;
    } finally {
      if (identical(_connecting, job)) {
        _connecting = null;
      }
    }
  }

  Future<ClientTransportConnection> _connectOnce() async {
    final conn = ClientTransportConnection.viaStreams(_sock, _sock);
    _conn = conn;
    return conn;
  }

  Future<Uint8List> post(String path, Uint8List body) async {
    final conn = await _connection();
    final headers = <Header>[
      Header.ascii(':method', 'POST'),
      Header.ascii(':path', path),
      Header.ascii(':scheme', 'http'),
      Header.ascii('content-type', rbProtoContentType),
      Header.ascii('accept', rbProtoContentType),
      Header.ascii('accept-encoding', RbZstdCodec.contentEncoding),
    ];
    var payload = body;
    if (body.isNotEmpty) {
      final packed = await RbZstdCodec.maybeCompressRequest(body);
      payload = packed.bytes;
      if (packed.zstd) {
        headers.add(Header.ascii('content-encoding', RbZstdCodec.contentEncoding));
      }
    }
    final stream = conn.makeRequest(headers, endStream: payload.isEmpty);
    if (payload.isNotEmpty) {
      stream.sendData(payload, endStream: true);
    }
    return _readResponse(stream, requestPath: path);
  }

  Future<Uint8List> get(String path) async {
    final conn = await _connection();
    final stream = conn.makeRequest([
      Header.ascii(':method', 'GET'),
      Header.ascii(':path', path),
      Header.ascii(':scheme', 'http'),
      Header.ascii('accept', rbProtoContentType),
      Header.ascii('accept-encoding', RbZstdCodec.contentEncoding),
    ], endStream: true);
    return _readResponse(stream, requestPath: path);
  }

  Future<Uint8List> _readResponse(ClientTransportStream stream, {required String requestPath}) async {
    final buf = BytesBuilder(copy: false);
    var status = 0;
    String? responseEncoding;
    await for (final msg in stream.incomingMessages) {
      if (msg is HeadersStreamMessage) {
        for (final h in msg.headers) {
          final name = String.fromCharCodes(h.name).toLowerCase();
          if (name == ':status') {
            status = int.tryParse(String.fromCharCodes(h.value)) ?? 0;
          } else if (name == 'content-encoding') {
            responseEncoding = String.fromCharCodes(h.value);
          }
        }
      } else if (msg is DataStreamMessage) {
        if (msg.bytes.isNotEmpty) {
          buf.add(msg.bytes);
        }
      }
    }
    if (status >= 400) {
      final detail = utf8.decode(buf.toBytes(), allowMalformed: true).trim();
      final msg = detail.isEmpty ? 'HTTP $status' : 'HTTP $status: $detail';
      throw HttpException(msg, uri: Uri(path: requestPath));
    }
    final raw = buf.toBytes();
    return RbZstdCodec.decodeBody(raw, contentEncoding: responseEncoding);
  }

  Future<void> shutdown() async {
    _markDead();
    final c = _conn;
    _conn = null;
    if (c != null && c.isOpen) {
      try {
        await c.finish().timeout(const Duration(seconds: 3));
      } catch (_) {}
    }
    try {
      await _sock.close();
    } catch (_) {}
  }
}
