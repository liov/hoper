import 'dart:async';
import 'dart:io';
import 'dart:math' as math;
import 'dart:typed_data';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_file_media.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/rb_preview_store.dart';

/// 本机 loopback HTTP，按 Range 从远端 [RbGrpcSession.readFileRange] 拉块，供 video_player 边下边播。
class RbVideoStreamServer {
  RbVideoStreamServer._(
    this._server,
    this.playUri,
    this._wire,
    this._relPath,
    this._totalSize,
    this._mime,
    this._fileBaseOffset,
    this._entry,
    this._previewStore,
  );

  final HttpServer _server;
  final Uri playUri;
  final RbGrpcSession _wire;
  final String _relPath;
  int _totalSize;
  final String _mime;
  /// 读文件时叠加到 HTTP Range 偏移（Motion Photo 内嵌 MP4 起始）。
  final int _fileBaseOffset;
  final RbFileEntry? _entry;
  final RbPreviewContentStore? _previewStore;
  RandomAccessFile? _cacheWriter;

  /// 单次向播放器返回的上限（内网大视频宜大块，减少 loopback+ICE 往返）。
  static const _httpResponseMax = 16 << 20;
  /// 单次 ICE ReadFile 块（与 Agent cap 一致）。
  static const _iceReadChunk = rbReadFileRangeMax;
  /// 与 Agent read_sem 对齐，允许多 Range 并发拉流。
  static const _maxInflight = 8;
  var _inflight = 0;
  var _closed = false;
  final _cacheWriteQueue = <Future<void>>[];

  Future<void> close() async {
    _closed = true;
    if (_cacheWriteQueue.isNotEmpty) {
      await Future.wait(List<Future<void>>.from(_cacheWriteQueue));
    }
    final w = _cacheWriter;
    _cacheWriter = null;
    final store = _previewStore;
    final ent = _entry;
    if (w != null && ent != null && store != null && _totalSize > 0) {
      await store.finalizeWriter(ent, _relPath, logicalTotal: _totalSize, raf: w, fileBaseOffset: _fileBaseOffset);
    } else {
      await w?.close();
    }
    await _server.close(force: true);
  }

  static Future<RbVideoStreamServer> start({
    required RbGrpcSession wire,
    required String relPath,
    required int totalSize,
    String? mime,
    String? fileName,
    int fileBaseOffset = 0,
    RbFileEntry? entry,
    RbPreviewContentStore? previewStore,
  }) async {
    final server = await HttpServer.bind(InternetAddress.loopbackIPv4, 0);
    final path = '/rb/${relPath.hashCode.abs()}_${DateTime.now().microsecondsSinceEpoch}';
    final m = mime ?? (fileName != null ? rbVideoMimeType(fileName) : 'video/mp4');
    late final RbVideoStreamServer inst;
    inst = RbVideoStreamServer._(
      server,
      Uri.parse('http://127.0.0.1:${server.port}$path'),
      wire,
      relPath,
      totalSize,
      m,
      fileBaseOffset,
      entry,
      previewStore,
    );
    server.listen((req) {
      if (req.uri.path != path) {
        req.response.statusCode = HttpStatus.notFound;
        unawaited(req.response.close());
        return;
      }
      unawaited(inst._serve(req));
    });
    await probeReady(inst.playUri);
    return inst;
  }

  /// 确认 loopback HTTP 已能连上，再交给 ExoPlayer，避免「端口已关」的 ConnectException。
  static Future<void> probeReady(Uri uri) async {
    final client = HttpClient();
    try {
      final req = await client.headUrl(uri).timeout(const Duration(seconds: 8));
      final resp = await req.close().timeout(const Duration(seconds: 8));
      if (resp.statusCode >= 500) {
        throw StateError('本地视频代理未就绪 (HTTP ${resp.statusCode})');
      }
      await resp.drain();
    } on SocketException catch (e) {
      throw StateError('无法连接本地视频代理: $e');
    } finally {
      client.close(force: true);
    }
  }

  Future<void> _acquireSlot() async {
    for (var i = 0; i < 400; i++) {
      if (_closed) {
        throw StateError('stream closed');
      }
      if (_inflight < _maxInflight) {
        _inflight++;
        return;
      }
      await Future<void>.delayed(const Duration(milliseconds: 25));
    }
    throw StateError('video stream busy');
  }

  void _releaseSlot() {
    if (_inflight > 0) {
      _inflight--;
    }
  }

  Future<void> _serve(HttpRequest req) async {
    if (_closed) {
      req.response.statusCode = HttpStatus.serviceUnavailable;
      await req.response.close();
      return;
    }
    if (req.method != 'GET' && req.method != 'HEAD') {
      req.response.statusCode = HttpStatus.methodNotAllowed;
      await req.response.close();
      return;
    }
    final isHead = req.method == 'HEAD';
    if (!isHead) {
      try {
        await _acquireSlot();
      } catch (_) {
        try {
          req.response.statusCode = HttpStatus.serviceUnavailable;
          await req.response.close();
        } catch (_) {}
        return;
      }
    }
    try {
      if (_closed || !_wire.isOpen) {
        req.response.statusCode = HttpStatus.serviceUnavailable;
        await req.response.close();
        return;
      }
      var total = _totalSize;
      if (total <= 0) {
        final probe = await _wire.readFileRange(_relPath, offset: 0, length: 1);
        total = probe.totalSize;
        if (total > 0) {
          _totalSize = total;
        }
      }
      if (total <= 0) {
        req.response.statusCode = HttpStatus.notFound;
        await req.response.close();
        return;
      }
      final parsed = _parseRange(req.headers.value(HttpHeaders.rangeHeader), total);
      var start = parsed.start;
      var end = parsed.end;
      var partial = parsed.partial;
      var len = end - start + 1;
      final cap = _httpResponseMax;
      if (len > cap) {
        end = start + cap - 1;
        len = cap;
        partial = true;
      }
      if (req.method == 'HEAD') {
        _setResponseHeaders(req.response, total: total, start: start, end: end, len: len, partial: partial);
        await req.response.close();
        return;
      }
      _setResponseHeaders(req.response, total: total, start: start, end: end, len: len, partial: partial);
      var pos = start;
      var sent = 0;
      while (pos <= end && sent < len) {
        if (_closed || !_wire.isOpen) {
          break;
        }
        final chunkLen = math.min(_iceReadChunk, end - pos + 1);
        final chunk = await _readChunk(pos, chunkLen);
        if (chunk.isEmpty) {
          break;
        }
        req.response.add(chunk);
        sent += chunk.length;
        pos += chunk.length;
        if (chunk.length < chunkLen) {
          break;
        }
      }
      await req.response.close();
    } catch (e) {
      try {
        if (!_closed) {
          req.response.statusCode = HttpStatus.internalServerError;
          await req.response.close();
        }
      } catch (_) {}
    } finally {
      if (!isHead) {
        _releaseSlot();
      }
    }
  }

  Future<Uint8List> _readChunk(int logicalOffset, int chunkLen) async {
    final entry = _entry;
    final store = _previewStore;
    if (entry != null && store != null) {
      final hit = await store.readRange(
        entry,
        _relPath,
        logicalOffset: logicalOffset,
        length: chunkLen,
        fileBaseOffset: _fileBaseOffset,
        logicalTotal: _totalSize,
      );
      if (hit != null && hit.isNotEmpty) {
        return hit;
      }
    }
    final r = await _wire.readFileRange(_relPath, offset: logicalOffset + _fileBaseOffset, length: chunkLen);
    final bytes = r.bytes;
    if (bytes.isEmpty) {
      return Uint8List(0);
    }
    if (entry != null && store != null) {
      _cacheWriter ??= await store.openWriter(entry, _relPath, logicalTotal: _totalSize, fileBaseOffset: _fileBaseOffset);
      final w = _cacheWriter;
      if (w != null) {
        _enqueueCacheWrite(store.writeRangeAt(w, logicalOffset, bytes));
      }
    }
    return bytes;
  }

  void _enqueueCacheWrite(Future<void> op) {
    _cacheWriteQueue.add(op);
    unawaited(op.whenComplete(() {
      _cacheWriteQueue.remove(op);
    }));
  }

  void _setResponseHeaders(
    HttpResponse resp, {
    required int total,
    required int start,
    required int end,
    required int len,
    required bool partial,
  }) {
    resp.headers.set(HttpHeaders.acceptRangesHeader, 'bytes');
    resp.headers.contentType = ContentType.parse(_mime);
    if (partial) {
      resp.statusCode = HttpStatus.partialContent;
      resp.headers.set(HttpHeaders.contentRangeHeader, 'bytes $start-$end/$total');
    } else {
      resp.statusCode = HttpStatus.ok;
    }
    resp.contentLength = len;
  }
}

class _RangeSlice {
  _RangeSlice({required this.start, required this.end, required this.partial});

  final int start;
  final int end;
  final bool partial;
}

_RangeSlice _parseRange(String? header, int total) {
  if (total <= 0) {
    return _RangeSlice(start: 0, end: 0, partial: false);
  }
  if (header == null || !header.startsWith('bytes=')) {
    return _RangeSlice(start: 0, end: total - 1, partial: false);
  }
  final spec = header.substring(6).split(',').first.trim();
  final dash = spec.indexOf('-');
  if (dash < 0) {
    return _RangeSlice(start: 0, end: total - 1, partial: false);
  }
  final a = spec.substring(0, dash);
  final b = spec.substring(dash + 1);
  if (a.isEmpty) {
    final suffix = int.tryParse(b) ?? 0;
    final start = (total - suffix).clamp(0, total - 1);
    return _RangeSlice(start: start, end: total - 1, partial: true);
  }
  final start = int.parse(a).clamp(0, total - 1);
  final end = b.isEmpty ? total - 1 : int.parse(b).clamp(start, total - 1);
  return _RangeSlice(start: start, end: end, partial: true);
}
