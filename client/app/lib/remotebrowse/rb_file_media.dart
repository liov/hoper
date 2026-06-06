import 'package:path/path.dart' as p;

/// `.ts` 多为 TypeScript；纯数字或 segment/chunk 等前缀视为 MPEG-TS 分片。
bool rbFileIsLikelyMpegTsSegment(String fileName) {
  if (p.extension(fileName).toLowerCase() != '.ts') {
    return false;
  }
  final stem = p.basenameWithoutExtension(fileName);
  if (stem.isEmpty) {
    return false;
  }
  if (RegExp(r'^\d+$').hasMatch(stem)) {
    return true;
  }
  final lower = stem.toLowerCase();
  return lower.startsWith('segment') ||
      lower.startsWith('chunk') ||
      lower.startsWith('stream') ||
      lower.startsWith('media_') ||
      lower.startsWith('video_');
}

/// 远程浏览支持预览/转码/缩略图的视频扩展名（不含普通 `.ts` 源码）。
bool rbFileIsVideo(String fileName) {
  final ext = p.extension(fileName).toLowerCase();
  switch (ext) {
    case '.mp4':
    case '.mov':
    case '.avi':
    case '.mkv':
    case '.webm':
    case '.m4v':
    case '.flv':
    case '.rmvb':
    case '.3gp':
    case '.m2ts':
    case '.mts':
      return true;
    case '.ts':
      return rbFileIsLikelyMpegTsSegment(fileName);
    default:
      return false;
  }
}

String rbFormatFileSize(int bytes) {
  if (bytes < 1024) {
    return '$bytes B';
  }
  if (bytes < 1024 * 1024) {
    return '${(bytes / 1024).toStringAsFixed(1)} KB';
  }
  if (bytes < 1024 * 1024 * 1024) {
    return '${(bytes / (1024 * 1024)).toStringAsFixed(1)} MB';
  }
  return '${(bytes / (1024 * 1024 * 1024)).toStringAsFixed(2)} GB';
}

/// 将毫秒格式化为 `m:ss` 或 `h:mm:ss`；无效返回空串。
String rbFormatDurationMs(int ms) {
  if (ms <= 0) {
    return '';
  }
  final totalSec = ms ~/ 1000;
  final h = totalSec ~/ 3600;
  final m = (totalSec % 3600) ~/ 60;
  final s = totalSec % 60;
  if (h > 0) {
    return '$h:${m.toString().padLeft(2, '0')}:${s.toString().padLeft(2, '0')}';
  }
  return '$m:${s.toString().padLeft(2, '0')}';
}

String rbVideoMimeType(String fileName) {
  switch (p.extension(fileName).toLowerCase()) {
    case '.webm':
      return 'video/webm';
    case '.mov':
    case '.m4v':
      return 'video/quicktime';
    case '.mkv':
      return 'video/x-matroska';
    case '.avi':
      return 'video/x-msvideo';
    case '.ts':
    case '.m2ts':
    case '.mts':
      return 'video/mp2t';
    default:
      return 'video/mp4';
  }
}
