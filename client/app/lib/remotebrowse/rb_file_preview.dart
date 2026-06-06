import 'dart:typed_data';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_file_media.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:path/path.dart' as p;

enum RbPreviewKind { image, video, motionPhoto, text, markdown, html, pdf, unsupported }

String rbPreviewKindLabel(RbPreviewKind kind) => switch (kind) {
      RbPreviewKind.image => '图片',
      RbPreviewKind.video => '视频',
      RbPreviewKind.motionPhoto => '动态照片',
      RbPreviewKind.text => '文本',
      RbPreviewKind.markdown => 'Markdown',
      RbPreviewKind.html => 'HTML',
      RbPreviewKind.pdf => 'PDF',
      RbPreviewKind.unsupported => '不支持预览',
    };

const rbPreviewTextMaxBytes = 4 << 20;
const rbPreviewPdfMaxBytes = 24 << 20;
const rbPreviewImageMaxBytes = 16 << 20;

/// Agent 预览：源 bpp > 此值时返回 0.3bpp WebP；否则直传原文件（Flutter 可解码时）。
const rbPreviewImageBppThreshold = 0.3;

/// 「原图」：请求 Agent 全分辨率高清 WebP（目标 bpp 2.0），非磁盘原文件。
const rbReadOriginalPathSuffix = '#rb-original';

String rbReadOriginalRequestPath(String relPath) => '${rbNormRemotePath(relPath)}$rbReadOriginalPathSuffix';

/// Flutter 无法解码、预览须 Agent 转 0.3bpp WebP（与 Agent [preview_transcode_for_client] 一致）。
bool rbImageNeedsAgentPreviewWebp(String fileName) {
  switch (p.extension(fileName).toLowerCase()) {
    case '.avif':
    case '.jxl':
    case '.heic':
    case '.heif':
    case '.hif':
    case '.tif':
    case '.tiff':
      return true;
    default:
      return false;
  }
}

/// Agent ReadFile 返回的是转码预览（非磁盘原图）时为 true。
bool rbImageAgentPreviewTranscoded({
  required String mime,
  required String fileName,
  required int totalSize,
  required int entrySize,
}) {
  final ext = p.extension(fileName).toLowerCase();
  if (mime == 'image/webp' && ext != '.webp') {
    return true;
  }
  if (rbImageNeedsAgentPreviewWebp(fileName) && mime == 'image/webp') {
    return true;
  }
  if (entrySize > 0 && totalSize > 0 && totalSize != entrySize) {
    if (ext == '.webp' || mime == 'image/webp') {
      return true;
    }
    if (totalSize < entrySize) {
      return true;
    }
  }
  return false;
}

bool rbBytesLookLikeWebp(Uint8List bytes) =>
    bytes.length >= 12 &&
    bytes[0] == 0x52 &&
    bytes[1] == 0x49 &&
    bytes[2] == 0x46 &&
    bytes[3] == 0x46 &&
    bytes[8] == 0x57 &&
    bytes[9] == 0x45 &&
    bytes[10] == 0x42 &&
    bytes[11] == 0x50;

/// 当前预览字节是否为 Agent 原图（非转码 WebP 预览）。
bool rbImagePreviewIsOriginal({
  required String fileName,
  required String mime,
  required int dataBytes,
  required int entrySize,
  Uint8List? bytesHead,
}) {
  if (rbImageAgentPreviewTranscoded(mime: mime, fileName: fileName, totalSize: dataBytes, entrySize: entrySize)) {
    return false;
  }
  final head = bytesHead;
  if (head != null && head.isNotEmpty && mime.isEmpty && rbBytesLookLikeWebp(head)) {
    if (p.extension(fileName).toLowerCase() != '.webp') {
      return false;
    }
  }
  return true;
}

/// Flutter [Image.memory] 可直接绘制的格式；AVIF/HEIC/JXL 等须由 Agent ReadFile/缩略图转 WebP。
bool rbFlutterLikelyDecodesImageBytes(String fileName, String mime) {
  if (mime.isNotEmpty) {
    switch (mime) {
      case 'image/jpeg':
      case 'image/png':
      case 'image/gif':
      case 'image/webp':
      case 'image/bmp':
        return true;
      default:
        break;
    }
  }
  switch (p.extension(fileName).toLowerCase()) {
    case '.jpg':
    case '.jpeg':
    case '.png':
    case '.gif':
    case '.webp':
    case '.bmp':
      return true;
    default:
      return false;
  }
}

/// 当前字节是否可在 Flutter 中绘制（含 Agent 转好的 WebP）；不在客户端做 HEIC/JXL/AVIF 解码。
bool rbPreviewBytesDisplayable({required String fileName, required String mime, Uint8List? bytesHead}) {
  if (rbImageNeedsAgentPreviewWebp(fileName)) {
    final head = bytesHead;
    return head != null && head.isNotEmpty && (mime == 'image/webp' || rbBytesLookLikeWebp(head));
  }
  final head = bytesHead;
  if (head != null && head.isNotEmpty && rbBytesLookLikeWebp(head)) {
    return true;
  }
  return rbFlutterLikelyDecodesImageBytes(fileName, mime);
}

/// 旧版误把 512 缩略图写入预览缓存时跳过（原图预览应明显更大）。
bool rbIsLikelyStaleImagePreviewCache(RbFileEntry entry, int cachedBytes) {
  if (entry.size <= 0 || cachedBytes <= 0) {
    return false;
  }
  if (cachedBytes >= entry.size) {
    return false;
  }
  if (cachedBytes >= 256 << 10) {
    return false;
  }
  return cachedBytes < entry.size ~/ 8;
}

bool rbIsPreviewable(String fileName) => rbPreviewKind(fileName) != RbPreviewKind.unsupported;

/// 全屏图片预览（非视频/动态照片）。
bool rbIsStillImagePreviewPath(String relPath) => rbPreviewKind(p.basename(relPath)) == RbPreviewKind.image;

bool rbIsPreviewableEntry(RbFileEntry entry) {
  if (entry.isDirectory) {
    return false;
  }
  if (entry.isMotionPhoto) {
    return true;
  }
  return rbIsPreviewable(entry.name);
}

RbPreviewKind rbPreviewKindForEntry(RbFileEntry entry) {
  if (entry.isMotionPhoto) {
    return RbPreviewKind.motionPhoto;
  }
  return rbPreviewKind(entry.name);
}

RbPreviewKind rbPreviewKind(String fileName) {
  switch (p.extension(fileName).toLowerCase()) {
    case '.jpg':
    case '.jpeg':
    case '.png':
    case '.gif':
    case '.webp':
    case '.avif':
    case '.jxl':
    case '.bmp':
    case '.heic':
    case '.heif':
    case '.tif':
    case '.tiff':
      return RbPreviewKind.image;
    case '.md':
    case '.markdown':
      return RbPreviewKind.markdown;
    case '.html':
    case '.htm':
      return RbPreviewKind.html;
    case '.txt':
    case '.json':
    case '.xml':
    case '.yaml':
    case '.yml':
    case '.log':
    case '.ini':
    case '.conf':
    case '.dart':
    case '.go':
    case '.rs':
    case '.js':
    case '.py':
    case '.sh':
    case '.sql':
    case '.proto':
    case '.css':
    case '.csv':
      return RbPreviewKind.text;
    case '.pdf':
      return RbPreviewKind.pdf;
    case '.ts':
      if (rbFileIsLikelyMpegTsSegment(fileName)) {
        return RbPreviewKind.video;
      }
      return RbPreviewKind.text;
    default:
      if (rbFileIsVideo(fileName)) {
        return RbPreviewKind.video;
      }
      return RbPreviewKind.unsupported;
  }
}

int rbPreviewMaxBytes(RbPreviewKind kind) => switch (kind) {
      RbPreviewKind.image => rbPreviewImageMaxBytes,
      RbPreviewKind.video => rbReadFileRangeMaxBytes,
      RbPreviewKind.pdf => rbPreviewPdfMaxBytes,
      _ => rbPreviewTextMaxBytes,
    };

/// 与 Agent ReadFile 单次 length 上限一致；视频预览为 Range 流式，不整文件下载。
const rbReadFileRangeMaxBytes = 2 << 20;
const rbPreviewVideoMaxBytes = rbReadFileRangeMaxBytes;
