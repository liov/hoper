import 'dart:typed_data';

import 'package:zstandard/zstandard.dart';

/// 与 rb `proto_zstd` 对齐：HTTP protobuf 请求/响应 zstd。
abstract final class RbZstdCodec {
  static const contentEncoding = 'zstd';
  static const minCompressBytes = 64;
  static const compressionLevel = 3;

  static final Zstandard _z = Zstandard();

  static Future<({Uint8List bytes, bool zstd})> maybeCompressRequest(Uint8List plain) async {
    if (plain.length < minCompressBytes) {
      return (bytes: plain, zstd: false);
    }
    final out = await _z.compress(plain, compressionLevel);
    if (out == null || out.length >= plain.length) {
      return (bytes: plain, zstd: false);
    }
    return (bytes: out, zstd: true);
  }

  static Future<Uint8List> decodeBody(Uint8List body, {String? contentEncoding}) async {
    final enc = contentEncoding?.trim().toLowerCase();
    if (enc != RbZstdCodec.contentEncoding) {
      return body;
    }
    final out = await _z.decompress(body);
    if (out == null) {
      throw StateError('zstd 解压失败');
    }
    return out;
  }
}
