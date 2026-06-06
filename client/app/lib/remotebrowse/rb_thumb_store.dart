import 'dart:convert';
import 'dart:io';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:crypto/crypto.dart';
import 'package:flutter/foundation.dart';
import 'package:path_provider/path_provider.dart';

/// Viewer 本地缩略图缓存：按 Agent host + 相对路径分桶，默认保留 7 天。
class RbThumbStore {
  RbThumbStore({required String hostKey, this._cacheTtl = const Duration(days: 7)})
      : _hostKey = hostKey.trim().isEmpty ? 'rb' : hostKey.trim();

  final String _hostKey;
  final Duration _cacheTtl;
  Directory? _dir;

  static String sanitizeHostKey(String hostKey) =>
      hostKey.replaceAll(RegExp(r'[^\w.\-]+'), '_').replaceAll(RegExp(r'_+'), '_');

  String _dataFileName(String relPath, int maxEdge) {
    final digest = sha256.convert(utf8.encode('$_hostKey|${rbNormRemotePath(relPath)}|$maxEdge'));
    return '${digest.toString().substring(0, 32)}.webp';
  }

  Future<Directory?> _cacheDir() async {
    if (kIsWeb) {
      return null;
    }
    if (_dir != null) {
      return _dir;
    }
    final base = await getTemporaryDirectory();
    _dir = Directory('${base.path}/rb_thumb_cache/${sanitizeHostKey(_hostKey)}');
    if (!_dir!.existsSync()) {
      await _dir!.create(recursive: true);
    }
    return _dir;
  }

  bool _metaMatchesEntry(Map<String, dynamic> meta, String relPath, RbFileEntry entry, int maxEdge) {
    if (meta['host'] != _hostKey) {
      return false;
    }
    if (rbNormRemotePath(meta['path'] as String? ?? '') != rbNormRemotePath(relPath)) {
      return false;
    }
    if (meta['mtime'] != entry.mtimeUnixMs || meta['size'] != entry.size) {
      return false;
    }
    if ((meta['maxEdge'] as num?)?.toInt() != maxEdge) {
      return false;
    }
    final h = entry.thumbHash;
    if (h.isNotEmpty && meta['thumbHash'] != h) {
      return false;
    }
    final at = (meta['cachedAtMs'] as num?)?.toInt() ?? 0;
    if (at <= 0 || DateTime.now().millisecondsSinceEpoch - at > _cacheTtl.inMilliseconds) {
      return false;
    }
    return true;
  }

  Future<Uint8List?> load(RbFileEntry entry, String relPath, int maxEdge) async {
    final dir = await _cacheDir();
    if (dir == null) {
      return null;
    }
    final dataFile = File('${dir.path}/${_dataFileName(relPath, maxEdge)}');
    final metaFile = File('${dataFile.path}.meta');
    if (!dataFile.existsSync() || !metaFile.existsSync()) {
      return null;
    }
    try {
      final meta = jsonDecode(await metaFile.readAsString()) as Map<String, dynamic>;
      if (!_metaMatchesEntry(meta, relPath, entry, maxEdge)) {
        await _deletePair(dataFile, metaFile);
        return null;
      }
      return await dataFile.readAsBytes();
    } catch (_) {
      await _deletePair(dataFile, metaFile);
      return null;
    }
  }

  Future<void> put(RbFileEntry entry, String relPath, int maxEdge, Uint8List bytes) async {
    if (bytes.isEmpty) {
      return;
    }
    final dir = await _cacheDir();
    if (dir == null) {
      return;
    }
    final dataFile = File('${dir.path}/${_dataFileName(relPath, maxEdge)}');
    final meta = {
      'host': _hostKey,
      'path': rbNormRemotePath(relPath),
      'mtime': entry.mtimeUnixMs,
      'size': entry.size,
      'thumbHash': entry.thumbHash,
      'maxEdge': maxEdge,
      'cachedAtMs': DateTime.now().millisecondsSinceEpoch,
    };
    await dataFile.writeAsBytes(bytes, flush: true);
    await File('${dataFile.path}.meta').writeAsString(jsonEncode(meta), flush: true);
  }

  Future<void> clear() async {
    final dir = await _cacheDir();
    if (dir == null || !dir.existsSync()) {
      return;
    }
    await for (final e in dir.list()) {
      await e.delete(recursive: true);
    }
  }

  Future<void> purgeExpired() async {
    final dir = await _cacheDir();
    if (dir == null || !dir.existsSync()) {
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
          final dataPath = e.path.substring(0, e.path.length - 5);
          await _deletePair(File(dataPath), File(e.path));
        }
      } catch (_) {}
    }
  }

  Future<void> removeForRelPath(String relPath) async {
    final dir = await _cacheDir();
    if (dir == null || !dir.existsSync()) {
      return;
    }
    final want = rbNormRemotePath(relPath);
    await for (final e in dir.list()) {
      if (!e.path.endsWith('.meta')) {
        continue;
      }
      try {
        final meta = jsonDecode(await File(e.path).readAsString()) as Map<String, dynamic>;
        if (meta['host'] != _hostKey) {
          continue;
        }
        if (rbNormRemotePath(meta['path'] as String? ?? '') == want) {
          final dataPath = e.path.substring(0, e.path.length - 5);
          await _deletePair(File(dataPath), File(e.path));
        }
      } catch (_) {}
    }
  }

  Future<void> removeForThumbHash(String thumbHash) async {
    if (thumbHash.isEmpty) {
      return;
    }
    final dir = await _cacheDir();
    if (dir == null || !dir.existsSync()) {
      return;
    }
    await for (final e in dir.list()) {
      if (!e.path.endsWith('.meta')) {
        continue;
      }
      try {
        final meta = jsonDecode(await File(e.path).readAsString()) as Map<String, dynamic>;
        if (meta['host'] == _hostKey && meta['thumbHash'] == thumbHash) {
          final dataPath = e.path.substring(0, e.path.length - 5);
          await _deletePair(File(dataPath), File(e.path));
        }
      } catch (_) {}
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
}
