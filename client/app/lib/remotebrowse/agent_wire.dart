import 'dart:io';
import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:app/gen/pb/remotebrowse/browse.service.pb.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/wire_codec.dart';
import 'package:app/remotebrowse/wire_transport.dart';
import 'package:fixnum/fixnum.dart';
import 'package:image/image.dart' as img;

class RbAgentWire {
  static Future<void> serve(RbWireTransport transport, String root) async {
    try {
      while (true) {
        final frame = await transport.readFrame();
        await _handle(transport, root, frame.$1, frame.$2);
      }
    } catch (_) {
      return;
    }
  }

  static Future<void> _handle(RbWireTransport transport, String root, int typ, Uint8List payload) async {
    switch (typ) {
      case rbTypeFileIndex:
        await _replyList(transport, root, payload);
      case rbTypeThumbReq:
        await _replyThumb(transport, root, payload);
      default:
        break;
    }
  }

  static Future<void> _replyList(RbWireTransport transport, String root, Uint8List payload) async {
    final req = ListFilesRequest.fromBuffer(payload);
    final path = req.rootPath.isEmpty ? root : req.rootPath;
    final entries = await _listFiles(path);
    await transport.writeFrame(rbTypeFileIndex, ListFilesResponse(entries: entries).writeToBuffer());
  }

  static Future<void> _replyThumb(RbWireTransport transport, String root, Uint8List payload) async {
    final req = ThumbnailRequest.fromBuffer(payload);
    final file = _resolvePath(root, req.path);
    if (file == null || !await file.exists()) {
      await transport.writeFrame(rbTypeThumbData, ThumbnailResponse().writeToBuffer());
      return;
    }
    final maxEdge = req.maxEdge == 0 ? 256 : req.maxEdge;
    try {
      final bytes = await _thumbBytes(file, maxEdge);
      await transport.writeFrame(
        rbTypeThumbData,
        ThumbnailResponse(data: bytes, mime: 'image/jpeg', thumbHash: req.hash).writeToBuffer(),
      );
    } catch (_) {
      final data = await RemoteBrowseApi().fetchThumb(file.path, maxEdge: maxEdge, hash: req.hash);
      await transport.writeFrame(
        rbTypeThumbData,
        ThumbnailResponse(data: data, mime: 'image/jpeg', thumbHash: req.hash).writeToBuffer(),
      );
    }
  }

  static Future<List<FileEntry>> _listFiles(String root) async {
    final dir = Directory(root);
    if (!await dir.exists()) {
      return [];
    }
    final out = <FileEntry>[];
    await for (final e in dir.list(followLinks: false)) {
      if (e is! File) {
        continue;
      }
      final name = e.uri.pathSegments.last;
      final stat = await e.stat();
      out.add(FileEntry(
        id: '${dir.path}:$name',
        name: name,
        size: Int64(stat.size),
        mtimeUnixMs: Int64(stat.modified.millisecondsSinceEpoch),
        mime: _guessMime(name),
      ));
    }
    out.sort((a, b) => a.name.compareTo(b.name));
    return out;
  }

  static String _guessMime(String name) {
    final lower = name.toLowerCase();
    if (lower.endsWith('.jpg') || lower.endsWith('.jpeg')) {
      return 'image/jpeg';
    }
    if (lower.endsWith('.png')) {
      return 'image/png';
    }
    if (lower.endsWith('.webp')) {
      return 'image/webp';
    }
    return 'application/octet-stream';
  }

  static File? _resolvePath(String root, String path) {
    if (path.isEmpty) {
      return null;
    }
    final f = File(path);
    if (f.isAbsolute) {
      return f;
    }
    return File('$root/${path.replaceFirst(RegExp(r'^/'), '')}');
  }

  static Future<Uint8List> _thumbBytes(File file, int maxEdge) async {
    final raw = await file.readAsBytes();
    final decoded = img.decodeImage(raw);
    if (decoded == null) {
      throw StateError('decode failed');
    }
    final edge = maxEdge.clamp(32, 1024);
    final resized = decoded.width >= decoded.height
        ? img.copyResize(decoded, width: edge)
        : img.copyResize(decoded, height: edge);
    return Uint8List.fromList(img.encodeJpg(resized, quality: 85));
  }
}
