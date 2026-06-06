import 'dart:convert';
import 'dart:io';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:crypto/crypto.dart';
import 'package:flutter/foundation.dart';
import 'package:path_provider/path_provider.dart';

/// 预览正文缓存（整文件或视频/PDF 流式落盘），默认保留 1 天。
class RbPreviewContentStore {
  RbPreviewContentStore({required String hostKey, this._cacheTtl = const Duration(days: 1)})
      : _hostKey = hostKey.trim().isEmpty ? 'rb' : hostKey.trim();

  final String _hostKey;
  final Duration _cacheTtl;
  Directory? _dir;

  static String sanitizeHostKey(String hostKey) =>
      hostKey.replaceAll(RegExp(r'[^\w.\-]+'), '_').replaceAll(RegExp(r'_+'), '_');

  String _baseName(String relPath, RbFileEntry entry, {int fileBaseOffset = 0, int? logicalTotal}) {
    final total = logicalTotal ?? entry.size;
    final digest = sha256.convert(utf8.encode(
      '$_hostKey|${rbNormRemotePath(relPath)}|${entry.mtimeUnixMs}|${entry.size}|$fileBaseOffset|$total',
    ));
    return digest.toString().substring(0, 32);
  }

  Future<Directory?> _cacheDir() async {
    if (kIsWeb) {
      return null;
    }
    if (_dir != null) {
      return _dir;
    }
    final base = await getTemporaryDirectory();
    _dir = Directory('${base.path}/rb_preview_cache/${sanitizeHostKey(_hostKey)}');
    if (!_dir!.existsSync()) {
      await _dir!.create(recursive: true);
    }
    return _dir;
  }

  File _dataFile(Directory dir, String base) => File('${dir.path}/$base.bin');
  File _metaFile(Directory dir, String base) => File('${dir.path}/$base.meta');

  bool _metaValid(Map<String, dynamic> meta, String relPath, RbFileEntry entry, {int fileBaseOffset = 0}) {
    if (meta['host'] != _hostKey) {
      return false;
    }
    if (rbNormRemotePath(meta['path'] as String? ?? '') != rbNormRemotePath(relPath)) {
      return false;
    }
    if (meta['mtime'] != entry.mtimeUnixMs || meta['size'] != entry.size) {
      return false;
    }
    if ((meta['fileBaseOffset'] as num?)?.toInt() != fileBaseOffset) {
      return false;
    }
    final at = (meta['cachedAtMs'] as num?)?.toInt() ?? 0;
    if (at <= 0 || DateTime.now().millisecondsSinceEpoch - at > _cacheTtl.inMilliseconds) {
      return false;
    }
    return true;
  }

  Future<Map<String, dynamic>?> _readMeta(Directory dir, String base, String relPath, RbFileEntry entry, {int fileBaseOffset = 0}) async {
    final metaFile = _metaFile(dir, base);
    if (!metaFile.existsSync()) {
      return null;
    }
    try {
      final meta = jsonDecode(await metaFile.readAsString()) as Map<String, dynamic>;
      if (!_metaValid(meta, relPath, entry, fileBaseOffset: fileBaseOffset)) {
        await _deletePair(_dataFile(dir, base), metaFile);
        return null;
      }
      return meta;
    } catch (_) {
      return null;
    }
  }

  Future<void> _deletePair(File data, File meta) async {
    if (data.existsSync()) {
      await data.delete();
    }
    if (meta.existsSync()) {
      await meta.delete();
    }
  }

  Future<Uint8List?> load(
    RbFileEntry entry,
    String relPath, {
    int maxBytes = 0,
    int fileBaseOffset = 0,
    int? logicalTotal,
  }) async {
    final file = await mediaFileIfComplete(entry, relPath, fileBaseOffset: fileBaseOffset, logicalTotal: logicalTotal);
    if (file == null) {
      return null;
    }
    final len = await file.length();
    if (maxBytes > 0 && len > maxBytes) {
      return null;
    }
    return file.readAsBytes();
  }

  Future<File?> mediaFileIfComplete(
    RbFileEntry entry,
    String relPath, {
    int fileBaseOffset = 0,
    int? logicalTotal,
  }) async {
    final dir = await _cacheDir();
    if (dir == null) {
      return null;
    }
    final base = _baseName(relPath, entry, fileBaseOffset: fileBaseOffset, logicalTotal: logicalTotal);
    final data = _dataFile(dir, base);
    final meta = await _readMeta(dir, base, relPath, entry, fileBaseOffset: fileBaseOffset);
    if (meta == null || meta['complete'] != true) {
      return null;
    }
    final expect = (meta['logicalTotal'] as num?)?.toInt() ?? 0;
    if (expect <= 0 || !data.existsSync()) {
      return null;
    }
    final len = await data.length();
    if (len != expect) {
      return null;
    }
    return data;
  }

  Future<void> put(RbFileEntry entry, String relPath, Uint8List bytes, {int fileBaseOffset = 0}) async {
    if (bytes.isEmpty) {
      return;
    }
    final dir = await _cacheDir();
    if (dir == null) {
      return;
    }
    final base = _baseName(relPath, entry, fileBaseOffset: fileBaseOffset, logicalTotal: bytes.length);
    final data = _dataFile(dir, base);
    final meta = {
      'host': _hostKey,
      'path': rbNormRemotePath(relPath),
      'mtime': entry.mtimeUnixMs,
      'size': entry.size,
      'fileBaseOffset': fileBaseOffset,
      'logicalTotal': bytes.length,
      'complete': true,
      'cachedAtMs': DateTime.now().millisecondsSinceEpoch,
    };
    await data.writeAsBytes(bytes, flush: true);
    await _metaFile(dir, base).writeAsString(jsonEncode(meta), flush: true);
  }

  Future<RandomAccessFile?> openWriter(
    RbFileEntry entry,
    String relPath, {
    required int logicalTotal,
    int fileBaseOffset = 0,
  }) async {
    if (logicalTotal <= 0) {
      return null;
    }
    final dir = await _cacheDir();
    if (dir == null) {
      return null;
    }
    final base = _baseName(relPath, entry, fileBaseOffset: fileBaseOffset, logicalTotal: logicalTotal);
    if (await mediaFileIfComplete(entry, relPath, fileBaseOffset: fileBaseOffset, logicalTotal: logicalTotal) != null) {
      return null;
    }
    final data = _dataFile(dir, base);
    final meta = {
      'host': _hostKey,
      'path': rbNormRemotePath(relPath),
      'mtime': entry.mtimeUnixMs,
      'size': entry.size,
      'fileBaseOffset': fileBaseOffset,
      'logicalTotal': logicalTotal,
      'complete': false,
      'cachedAtMs': DateTime.now().millisecondsSinceEpoch,
    };
    if (data.existsSync()) {
      await data.delete();
    }
    await data.create(recursive: true);
    await _metaFile(dir, base).writeAsString(jsonEncode(meta), flush: true);
    return data.open(mode: FileMode.write);
  }

  Future<void> writeRangeAt(RandomAccessFile raf, int logicalOffset, Uint8List bytes) async {
    if (bytes.isEmpty) {
      return;
    }
    await raf.setPosition(logicalOffset);
    await raf.writeFrom(bytes);
  }

  Future<void> finalizeWriter(
    RbFileEntry entry,
    String relPath, {
    required int logicalTotal,
    required RandomAccessFile raf,
    int fileBaseOffset = 0,
  }) async {
    await raf.flush();
    final len = await raf.length();
    await raf.close();
    if (len < logicalTotal) {
      return;
    }
    final dir = await _cacheDir();
    if (dir == null) {
      return;
    }
    final base = _baseName(relPath, entry, fileBaseOffset: fileBaseOffset, logicalTotal: logicalTotal);
    final metaFile = _metaFile(dir, base);
    try {
      final meta = jsonDecode(await metaFile.readAsString()) as Map<String, dynamic>;
      meta['complete'] = true;
      meta['cachedAtMs'] = DateTime.now().millisecondsSinceEpoch;
      await metaFile.writeAsString(jsonEncode(meta), flush: true);
    } catch (_) {}
  }

  Future<Uint8List?> readRange(
    RbFileEntry entry,
    String relPath, {
    required int logicalOffset,
    required int length,
    int fileBaseOffset = 0,
    int? logicalTotal,
  }) async {
    final file = await mediaFileIfComplete(entry, relPath, fileBaseOffset: fileBaseOffset, logicalTotal: logicalTotal);
    if (file == null) {
      return null;
    }
    final total = await file.length();
    if (logicalOffset < 0 || logicalOffset >= total) {
      return null;
    }
    final raf = await file.open();
    try {
      await raf.setPosition(logicalOffset);
      final want = (logicalOffset + length).clamp(0, total) - logicalOffset;
      final buf = Uint8List(want);
      final n = await raf.readInto(buf);
      if (n <= 0) {
        return null;
      }
      return n == buf.length ? buf : Uint8List.sublistView(buf, 0, n);
    } finally {
      await raf.close();
    }
  }

  Future<void> evict(String relPath) async {
    final dir = await _cacheDir();
    if (dir == null) {
      return;
    }
    final want = rbNormRemotePath(relPath);
    await for (final e in dir.list()) {
      if (!e.path.endsWith('.meta')) {
        continue;
      }
      try {
        final meta = jsonDecode(await File(e.path).readAsString()) as Map<String, dynamic>;
        if (meta['host'] == _hostKey && rbNormRemotePath(meta['path'] as String? ?? '') == want) {
          final name = File(e.path).uri.pathSegments.last;
          final base = name.endsWith('.meta') ? name.substring(0, name.length - 5) : name;
          await _deletePair(_dataFile(dir, base), File(e.path));
        }
      } catch (_) {}
    }
  }

  Future<void> purgeExpired() async {
    final dir = await _cacheDir();
    if (dir == null) {
      return;
    }
    final now = DateTime.now().millisecondsSinceEpoch;
    await for (final e in dir.list()) {
      if (!e.path.endsWith('.meta')) {
        continue;
      }
      try {
        final meta = jsonDecode(await File(e.path).readAsString()) as Map<String, dynamic>;
        final at = (meta['cachedAtMs'] as num?)?.toInt() ?? 0;
        if (at <= 0 || now - at > _cacheTtl.inMilliseconds) {
          final name = File(e.path).uri.pathSegments.last;
          final base = name.endsWith('.meta') ? name.substring(0, name.length - 5) : name;
          await _deletePair(_dataFile(dir, base), File(e.path));
        }
      } catch (_) {}
    }
  }
}
