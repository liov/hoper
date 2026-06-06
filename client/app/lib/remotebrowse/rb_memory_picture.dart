import 'dart:typed_data';

import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:flutter/material.dart';

/// 网格/缩略图：仅显示 Flutter 可解码或 Agent 已转 WebP 的字节，不在客户端转码。
class RbMemoryPicture extends StatelessWidget {
  const RbMemoryPicture({
    super.key,
    required this.bytes,
    required this.fileName,
    this.mime = '',
    this.fit = BoxFit.cover,
    this.filterQuality = FilterQuality.low,
    this.gaplessPlayback = true,
  });

  final Uint8List bytes;
  final String fileName;
  final String mime;
  final BoxFit fit;
  final FilterQuality filterQuality;
  final bool gaplessPlayback;

  @override
  Widget build(BuildContext context) {
    if (bytes.isEmpty) {
      return const Center(child: Icon(Icons.broken_image_outlined, size: 28));
    }
    if (!rbPreviewBytesDisplayable(fileName: fileName, mime: mime, bytesHead: bytes)) {
      return const Center(child: Icon(Icons.broken_image_outlined, size: 28));
    }
    return Image.memory(bytes, fit: fit, gaplessPlayback: gaplessPlayback, filterQuality: filterQuality);
  }
}
