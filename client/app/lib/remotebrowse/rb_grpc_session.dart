import 'dart:async';
import 'dart:io';
import 'dart:math' as math;
import 'dart:typed_data';

import 'package:fixnum/fixnum.dart';
import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:app/gen/pb/remotebrowse/browse.service.pb.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/direct_dialer.dart';
import 'package:app/remotebrowse/link_kind.dart';
import 'package:app/remotebrowse/rb_file_icon.dart';
import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:app/remotebrowse/rb_h2_transport.dart';
import 'package:app/remotebrowse/rb_ice_grpc.dart';
import 'package:app/remotebrowse/rb_media_http.dart';
import 'package:app/remotebrowse/rb_media_meta.dart';
import 'package:app/remotebrowse/rb_transcode.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';

class RbReadFileResult {
  RbReadFileResult({required this.bytes, required this.mime, required this.totalSize});

  final Uint8List bytes;
  final String mime;
  final int totalSize;
}

class RbReadRangeResult {
  RbReadRangeResult({required this.bytes, required this.mime, required this.totalSize});

  final Uint8List bytes;
  final String mime;
  final int totalSize;
}

/// 单次读块上限（与 Agent `read_file_once` 一致）。
const rbReadFileRangeMax = 2 << 20;

/// ICE ReadFile 全量拉取时的并发 Range 数（与 Agent read_sem 对齐）。
const _readFullParallel = 6;

/// 预览切文件时取消进行中的读请求，勿当作断线。
const rbReadSuperseded = 'RbReadSuperseded';

typedef RbReadProgressCallback = void Function(int received, int total);

/// P2P：直连/中继 HTTP/2 + protobuf；ICE 经 Rust FFI。
class RbGrpcSession {
  RbGrpcSession._(
    this.linkKind, {
    this._h2,
    this._ice,
    this.room = '',
    this.agentMediaHost,
    this.agentMediaPort,
  });

  final RbLinkKind linkKind;
  final String room;
  final String? agentMediaHost;
  final int? agentMediaPort;
  final RbBoundH2Socket? _h2;
  final RbIceGrpcBridge? _ice;
  var _closed = false;
  var _readGeneration = 0;

  bool get isOpen => !_closed && !(_h2?.broken ?? false);

  factory RbGrpcSession.tcp(
    RbBoundH2Socket socket,
    RbLinkKind kind, {
    String room = '',
    String? agentMediaHost,
    int? agentMediaPort,
  }) {
    return RbGrpcSession._(
      kind,
      h2: socket,
      room: room,
      agentMediaHost: agentMediaHost,
      agentMediaPort: agentMediaPort,
    );
  }

  factory RbGrpcSession.ice(
    RbIceGrpcBridge bridge, {
    String room = '',
    String? agentMediaHost,
    int? agentMediaPort,
  }) {
    return RbGrpcSession._(
      RbLinkKind.ice,
      ice: bridge,
      room: room,
      agentMediaHost: agentMediaHost,
      agentMediaPort: agentMediaPort,
    );
  }

  Uri? agentTranscodeStreamUri(
    String relPath, {
    required String presetId,
    String vcodec = RbTranscodeVcodec.hevc,
    int fileBaseOffset = 0,
    int startMs = 0,
  }) {
    if (linkKind == RbLinkKind.directTcp) {
      final sock = _h2?.socket;
      if (sock == null) {
        return null;
      }
      return rbAgentTranscodeStreamUri(
        host: sock.remoteAddress.address,
        port: sock.remotePort,
        relPath: relPath,
        room: room,
        presetId: presetId,
        vcodec: vcodec,
        fileBaseOffset: fileBaseOffset,
        startMs: startMs,
      );
    }
    final host = agentMediaHost?.trim();
    if (host == null || host.isEmpty) {
      return null;
    }
    return rbAgentTranscodeStreamUri(
      host: host,
      port: agentMediaPort ?? RbDirectDialer.defaultPort,
      relPath: relPath,
      room: room,
      presetId: presetId,
      vcodec: vcodec,
      fileBaseOffset: fileBaseOffset,
      startMs: startMs,
    );
  }

  Uri? agentTranscodeHlsUri(
    String relPath, {
    required String presetId,
    String vcodec = RbTranscodeVcodec.hevc,
    int fileBaseOffset = 0,
    int beginMs = 0,
  }) {
    if (linkKind == RbLinkKind.directTcp) {
      final sock = _h2?.socket;
      if (sock == null) {
        return null;
      }
      return rbAgentTranscodeHlsUri(
        host: sock.remoteAddress.address,
        port: sock.remotePort,
        relPath: relPath,
        room: room,
        presetId: presetId,
        vcodec: vcodec,
        fileBaseOffset: fileBaseOffset,
        beginMs: beginMs,
      );
    }
    final host = agentMediaHost?.trim();
    if (host == null || host.isEmpty) {
      return null;
    }
    return rbAgentTranscodeHlsUri(
      host: host,
      port: agentMediaPort ?? RbDirectDialer.defaultPort,
      relPath: relPath,
      room: room,
      presetId: presetId,
      vcodec: vcodec,
      fileBaseOffset: fileBaseOffset,
      beginMs: beginMs,
    );
  }

  Uri? agentMediaPlayUri(String relPath, {int fileBaseOffset = 0, int agentImagePreviewMode = 0}) {
    if (linkKind == RbLinkKind.directTcp) {
      final sock = _h2?.socket;
      if (sock == null) {
        return null;
      }
      return rbAgentMediaPlayUri(
        host: sock.remoteAddress.address,
        port: sock.remotePort,
        relPath: relPath,
        room: room,
        fileBaseOffset: fileBaseOffset,
        agentImagePreviewMode: agentImagePreviewMode,
      );
    }
    final host = agentMediaHost?.trim();
    if (host == null || host.isEmpty) {
      return null;
    }
    return rbAgentMediaPlayUri(
      host: host,
      port: agentMediaPort ?? RbDirectDialer.defaultPort,
      relPath: relPath,
      room: room,
      fileBaseOffset: fileBaseOffset,
      agentImagePreviewMode: agentImagePreviewMode,
    );
  }

  Uri? agentMediaMetaUri(String relPath) {
    if (linkKind == RbLinkKind.directTcp) {
      final sock = _h2?.socket;
      if (sock == null) {
        return null;
      }
      return rbAgentMediaMetaUri(
        host: sock.remoteAddress.address,
        port: sock.remotePort,
        relPath: relPath,
        room: room,
      );
    }
    final host = agentMediaHost?.trim();
    if (host == null || host.isEmpty) {
      return null;
    }
    return rbAgentMediaMetaUri(
      host: host,
      port: agentMediaPort ?? RbDirectDialer.defaultPort,
      relPath: relPath,
      room: room,
    );
  }

  Future<RbMediaMeta?> fetchAgentMediaMeta(String relPath) {
    final uri = agentMediaMetaUri(relPath);
    if (uri == null) {
      return Future.value(null);
    }
    return rbFetchAgentMediaMeta(uri);
  }

  String thumbCacheHostKey({String fallback = 'rb'}) {
    if (linkKind == RbLinkKind.directTcp) {
      final sock = _h2?.socket;
      if (sock != null) {
        return '${sock.remoteAddress.address}:${sock.remotePort}';
      }
    }
    final h = agentMediaHost?.trim();
    if (h != null && h.isNotEmpty) {
      return '$h:${agentMediaPort ?? RbDirectDialer.defaultPort}';
    }
    return fallback;
  }

  Future<void> close() async {
    if (_closed) {
      return;
    }
    _closed = true;
    _readGeneration++;
    _ice?.close();
    unawaited(_h2?.shutdown() ?? Future<void>.value());
  }

  void _ensureOpen() {
    if (_closed) {
      throw StateError('连接已关闭');
    }
  }

  void cancelPendingReads() {
    _readGeneration++;
  }

  Future<void> cancelPendingReadsAndWait() async {
    _readGeneration++;
  }

  void _checkReadGeneration(int gen) {
    if (_closed) {
      throw StateError('连接已关闭');
    }
    if (gen != _readGeneration) {
      throw StateError(rbReadSuperseded);
    }
  }


  Future<Uint8List> _post(String path, Uint8List body) async {
    final h2 = _h2;
    if (h2 == null) {
      throw StateError('ICE 通道请走 FFI');
    }
    return h2.post(path, body);
  }

  Future<RbListFilesResult> listFiles(String root) => _listFilesImpl(root);

  Future<RbListMediaResult> listMedia(String root, {MediaListCursor? cursor, int pageSize = 64}) =>
      _listMediaImpl(root, cursor: cursor ?? MediaListCursor(), pageSize: pageSize);

  RbFileEntry _mapFileEntry(FileEntry e) => RbFileEntry(
        id: e.id,
        name: e.name,
        size: e.size.toInt(),
        mtimeUnixMs: e.mtimeUnixMs.toInt(),
        thumbHash: e.thumbHash,
        isDirectory: (e.flags & rbFileFlagDirectory) != 0,
        durationMs: e.durationMs.toInt(),
        isMotionPhoto: (e.flags & rbFileFlagMotionPhoto) != 0,
        motionOffset: e.motionOffset.toInt(),
        motionLength: e.motionLength.toInt(),
        motionCompanion: e.motionCompanion,
      );

  Future<RbListFilesResult> _listFilesImpl(String root) async {
    _ensureOpen();
    final path = root.trim().isEmpty ? '.' : root.trim().replaceAll('\\', '/');
    final ListFilesResponse resp;
    if (_h2 != null) {
      final bytes = await _post(
        RbHttpPaths.list,
        ListFilesRequest(rootPath: path).writeToBuffer(),
      );
      resp = ListFilesResponse.fromBuffer(bytes);
    } else {
      resp = await _ice!.listFiles(path);
    }
    final entries = resp.entries.map(_mapFileEntry).toList();
    final resolved = rbNormResolvedRootPath(resp.resolvedRootPath);
    return RbListFilesResult(
      entries: entries,
      resolvedRootPath: resolved,
      hiddenSkipped: resp.hiddenSkipped,
    );
  }

  Future<RbListMediaResult> _listMediaImpl(String root, {required MediaListCursor cursor, required int pageSize}) async {
    _ensureOpen();
    final path = root.trim().isEmpty ? '.' : root.trim().replaceAll('\\', '/');
    final ListMediaResponse resp;
    if (_h2 != null) {
      final req = ListMediaRequest(rootPath: path, pageSize: pageSize);
      if (!rbMediaCursorEmpty(cursor)) {
        req.cursor = cursor;
      }
      final bytes = await _post(RbHttpPaths.listMedia, req.writeToBuffer());
      resp = ListMediaResponse.fromBuffer(bytes);
    } else {
      resp = await _ice!.listMedia(path, cursor: cursor, pageSize: pageSize);
    }
    final items = resp.items
        .where((m) => m.hasEntry())
        .map((m) => RbMediaItem(relPath: rbNormRemotePath(m.relPath), entry: _mapFileEntry(m.entry)))
        .toList();
    final resolved = rbNormResolvedRootPath(resp.resolvedRootPath);
    return RbListMediaResult(
      items: items,
      nextCursor: resp.nextCursor,
      done: resp.done,
      resolvedRootPath: resolved,
      scannedDirs: resp.scannedDirs,
    );
  }

  Future<Map<String, Uint8List>> fetchThumbsBatch(List<String> paths, {int maxEdge = 256}) =>
      _fetchThumbsBatchImpl(paths, maxEdge: maxEdge);

  Future<Map<String, Uint8List>> _fetchThumbsBatchImpl(List<String> paths, {int maxEdge = 256}) async {
    _ensureOpen();
    if (paths.isEmpty) {
      return {};
    }
    final reqPaths = paths.map(rbNormRemotePath).toList();
    final ThumbnailBatchResponse resp;
    if (_h2 != null) {
      final bytes = await _post(
        RbHttpPaths.thumbBatch,
        ThumbnailBatchRequest(paths: reqPaths, maxEdge: maxEdge).writeToBuffer(),
      );
      resp = ThumbnailBatchResponse.fromBuffer(bytes);
    } else {
      resp = await _ice!.thumbBatch(reqPaths, maxEdge);
    }
    final out = <String, Uint8List>{};
    void put(String key, Uint8List bytes) {
      if (key.isEmpty) {
        return;
      }
      out[key] = bytes;
      final n = rbNormRemotePath(key);
      if (n != key) {
        out[n] = bytes;
      }
    }

    for (var i = 0; i < resp.items.length; i++) {
      final item = resp.items[i];
      if (item.data.isEmpty) {
        if (item.error.isNotEmpty && i < reqPaths.length) {
          rbLog.warning('thumb fail path=${reqPaths[i]} err=${item.error}');
        }
        continue;
      }
      final bytes = Uint8List.fromList(item.data);
      if (i < paths.length) {
        put(paths[i], bytes);
      }
      if (i < reqPaths.length) {
        put(reqPaths[i], bytes);
      }
      if (item.path.isNotEmpty) {
        put(item.path, bytes);
      }
    }
    return out;
  }

  Future<RbReadFileResult> readFileFull(
    String path, {
    int maxBytes = rbPreviewPdfMaxBytes,
    RbReadProgressCallback? onProgress,
  }) async {
    _ensureOpen();
    final rel = rbNormRemotePath(path);
    if (rbIsStillImagePreviewPath(rel)) {
      return _readImagePreviewFull(rel, maxBytes: maxBytes, onProgress: onProgress);
    }
    if (_shouldReadFullViaAgentMedia(rel)) {
      final viaHttp = await _tryReadFullViaAgentMedia(rel, maxBytes: maxBytes, onProgress: onProgress);
      if (viaHttp != null) {
        return viaHttp;
      }
    }
    return _readFileFullImpl(rel, maxBytes: maxBytes, onProgress: onProgress);
  }

  /// 文本/PDF 等优先 Agent HTTP；图片见 [_readImagePreviewFull]。
  bool _shouldReadFullViaAgentMedia(String relPath) {
    final name = relPath.split('/').last;
    return !rbFileSupportsAgentThumb(name);
  }

  /// 内网/旁路：`preview=1` 由 Agent 按 bpp/格式返回原图或 0.3bpp WebP；无媒体端口走 ICE ReadFile。
  Future<RbReadFileResult> _readImagePreviewFull(
    String rel, {
    required int maxBytes,
    RbReadProgressCallback? onProgress,
  }) async {
    final prev = await _tryReadFullViaAgentMedia(
      rel,
      maxBytes: maxBytes,
      onProgress: onProgress,
      agentImagePreviewMode: 1,
    );
    if (prev != null) {
      rbLog.info('image preview agent http path=$rel bytes=${prev.bytes.length} mime=${prev.mime}');
      return prev;
    }
    rbLog.info('image preview ice read path=$rel agentMedia=${agentMediaHost ?? "none"}');
    return _readFileFullImpl(rel, maxBytes: maxBytes, onProgress: onProgress);
  }

  /// 「原图」：高清 WebP（`preview=2` / `#rb-original`），非磁盘原文件。
  Future<RbReadFileResult> readFileFullOriginal(
    String path, {
    int maxBytes = rbPreviewPdfMaxBytes,
    RbReadProgressCallback? onProgress,
  }) async {
    _ensureOpen();
    final rel = rbNormRemotePath(path);
    final hq = await _tryReadFullViaAgentMedia(
      rel,
      maxBytes: maxBytes,
      onProgress: onProgress,
      agentImagePreviewMode: 2,
    );
    if (hq != null) {
      rbLog.info('image hq webp agent http path=$rel bytes=${hq.bytes.length}');
      return hq;
    }
    return _readFileFullImpl(rbReadOriginalRequestPath(rel), maxBytes: maxBytes, onProgress: onProgress);
  }

  /// Agent 旁路 HTTP Range：单连接流式读原文件，避免 ICE protobuf 多次 2MB 往返。
  Future<RbReadFileResult?> _tryReadFullViaAgentMedia(
    String relPath, {
    required int maxBytes,
    int fileBaseOffset = 0,
    int agentImagePreviewMode = 0,
    RbReadProgressCallback? onProgress,
  }) async {
    final uri = agentMediaPlayUri(relPath, fileBaseOffset: fileBaseOffset, agentImagePreviewMode: agentImagePreviewMode);
    if (uri == null) {
      return null;
    }
    final client = HttpClient();
    client.connectionTimeout = const Duration(seconds: 12);
    try {
      final req = await client.getUrl(uri);
      if (maxBytes > 0) {
        req.headers.set(HttpHeaders.rangeHeader, 'bytes=0-${maxBytes - 1}');
      }
      final resp = await req.close().timeout(const Duration(minutes: 8));
      if (resp.statusCode != HttpStatus.ok && resp.statusCode != HttpStatus.partialContent) {
        rbLog.fine('agent media read http ${resp.statusCode}');
        return null;
      }
      final mime = resp.headers.contentType?.mimeType ?? 'application/octet-stream';
      final total = _contentRangeTotal(resp) ??
          int.tryParse(resp.headers.value(HttpHeaders.contentLengthHeader) ?? '') ??
          0;
      final buf = BytesBuilder(copy: false);
      void report() {
        final t = total > 0 ? total : maxBytes;
        onProgress?.call(buf.length, t);
      }
      report();
      await for (final chunk in resp) {
        if (buf.length >= maxBytes) {
          break;
        }
        final room = maxBytes - buf.length;
        if (chunk.length <= room) {
          buf.add(chunk);
        } else {
          buf.add(Uint8List.fromList(chunk.sublist(0, room)));
        }
        report();
      }
      final bytes = buf.toBytes();
      if (bytes.isEmpty) {
        return null;
      }
      return RbReadFileResult(bytes: bytes, mime: mime, totalSize: total > 0 ? total : bytes.length);
    } catch (e, st) {
      rbLog.fine('agent media read fail', e, st);
      return null;
    } finally {
      client.close(force: true);
    }
  }

  int? _contentRangeTotal(HttpClientResponse resp) {
    final raw = resp.headers.value(HttpHeaders.contentRangeHeader);
    if (raw == null || !raw.contains('/')) {
      return null;
    }
    final total = raw.split('/').last.trim();
    return int.tryParse(total);
  }

  Future<RbReadFileResult> _readFileFullImpl(
    String path, {
    int maxBytes = rbPreviewPdfMaxBytes,
    RbReadProgressCallback? onProgress,
  }) async {
    _ensureOpen();
    final gen = _readGeneration;
    final first = await _readFileProto(path, offset: 0, length: rbReadFileRangeMax, maxBytes: rbReadFileRangeMax);
    _checkReadGeneration(gen);
    if (first.error.isNotEmpty) {
      throw StateError(first.error);
    }
    var mime = first.mime.isNotEmpty ? first.mime : 'application/octet-stream';
    final total = first.totalSize.toInt();
    final buf = BytesBuilder(copy: false);
    void report() {
      final t = total > 0 ? math.min(total, maxBytes) : maxBytes;
      onProgress?.call(buf.length, t);
    }
    if (first.data.isNotEmpty) {
      final n = math.min(first.data.length, maxBytes);
      buf.add(first.data.length == n ? first.data : first.data.sublist(0, n));
    }
    report();
    if (first.eof || buf.length >= maxBytes || total <= 0 || buf.length >= total) {
      return RbReadFileResult(bytes: buf.toBytes(), mime: mime, totalSize: total > 0 ? total : buf.length);
    }
    var offset = buf.length;
    while (buf.length < maxBytes && offset < total) {
      _checkReadGeneration(gen);
      final futures = <Future<ReadFileResponse>>[];
      var batchEnd = offset;
      for (var i = 0; i < _readFullParallel && batchEnd < total && buf.length < maxBytes; i++) {
        final want = math.min(rbReadFileRangeMax, math.min(maxBytes - buf.length, total - batchEnd));
        if (want <= 0) {
          break;
        }
        futures.add(_readFileProto(path, offset: batchEnd, length: want, maxBytes: want));
        batchEnd += want;
      }
      if (futures.isEmpty) {
        break;
      }
      final chunks = await Future.wait(futures);
      _checkReadGeneration(gen);
      var done = false;
      for (final chunk in chunks) {
        if (chunk.error.isNotEmpty) {
          throw StateError(chunk.error);
        }
        if (chunk.mime.isNotEmpty) {
          mime = chunk.mime;
        }
        if (chunk.data.isEmpty) {
          done = true;
          continue;
        }
        final room = maxBytes - buf.length;
        if (room <= 0) {
          done = true;
          break;
        }
        if (chunk.data.length <= room) {
          buf.add(chunk.data);
        } else {
          buf.add(chunk.data.sublist(0, room));
          done = true;
        }
        if (chunk.eof) {
          done = true;
        }
      }
      offset = buf.length;
      report();
      if (done) {
        break;
      }
    }
    return RbReadFileResult(bytes: buf.toBytes(), mime: mime, totalSize: total > 0 ? total : buf.length);
  }

  Future<RbReadRangeResult> readFileRange(String path, {required int offset, required int length}) =>
      _readFileRangeImpl(path, offset: offset, length: length);

  Future<RbReadRangeResult> _readFileRangeImpl(String path, {required int offset, required int length}) async {
    _ensureOpen();
    final gen = _readGeneration;
    if (length <= 0 || length > rbReadFileRangeMax) {
      throw ArgumentError.value(length, 'length', '1..$rbReadFileRangeMax');
    }
    final chunk = await _readFileProto(path, offset: offset, length: length, maxBytes: length);
    _checkReadGeneration(gen);
    if (chunk.error.isNotEmpty) {
      throw StateError(chunk.error);
    }
    final raw = Uint8List.fromList(chunk.data);
    final bytes = raw.length > length ? Uint8List.sublistView(raw, 0, length) : raw;
    return RbReadRangeResult(
      bytes: bytes,
      mime: chunk.mime.isEmpty ? 'application/octet-stream' : chunk.mime,
      totalSize: chunk.totalSize.toInt(),
    );
  }

  Future<ReadFileResponse> _readFileProto(
    String path, {
    required int offset,
    required int length,
    required int maxBytes,
  }) async {
    final wirePath = rbNormRemotePath(path);
    if (wirePath.isEmpty) {
      throw StateError('文件路径无效');
    }
    final ice = _ice;
    if (ice != null) {
      if (offset == 0 && length == 0) {
        return ice.readFile(wirePath);
      }
      return ice.readFileRange(wirePath, offset: offset, length: length);
    }
    final reqLen = length == 0 ? rbReadFileRangeMax : length;
    final bytes = await _post(
      RbHttpPaths.read,
      ReadFileRequest(path: wirePath, offset: Int64(offset), length: reqLen).writeToBuffer(),
    );
    final resp = ReadFileResponse.fromBuffer(bytes);
    if (resp.data.length > maxBytes) {
      return ReadFileResponse(
        path: resp.path,
        totalSize: resp.totalSize,
        mime: resp.mime,
        data: resp.data.sublist(0, maxBytes),
        eof: resp.eof,
        error: resp.error,
      );
    }
    return resp;
  }

  Future<void> deleteFile(String path) => _deleteFileImpl(path);

  Future<void> _deleteFileImpl(String path) async {
    _ensureOpen();
    final wirePath = rbNormRemotePath(path);
    if (wirePath.isEmpty) {
      throw StateError('文件路径无效');
    }
    DeleteFileResponse resp;
    if (_h2 != null) {
      final bytes = await _post(
        RbHttpPaths.delete,
        DeleteFileRequest(path: wirePath).writeToBuffer(),
      ).timeout(const Duration(seconds: 120), onTimeout: () {
        throw StateError('删除超时：请重新编译并重启 Agent');
      });
      resp = DeleteFileResponse.fromBuffer(bytes);
    } else {
      resp = await _ice!.deleteFile(wirePath).timeout(const Duration(seconds: 120), onTimeout: () {
        throw StateError('删除超时：请重新编译并重启 Agent');
      });
    }
    if (resp.error.isNotEmpty) {
      throw StateError(resp.error);
    }
  }
}
