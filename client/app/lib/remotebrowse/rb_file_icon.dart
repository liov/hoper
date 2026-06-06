import 'package:app/remotebrowse/rb_file_media.dart';
import 'package:flutter/material.dart';
import 'package:path/path.dart' as p;

/// 宫格中仅图片/视频不叠文件名；目录与其它类型显示在缩略图上。
bool rbFileIsPhotoOrVideo(String fileName) => rbFileSupportsAgentThumb(fileName);

/// Agent `ensure_thumbnail` 仅支持图片与部分视频；`.ts` 默认按源码排除（见 [rbFileIsVideo]）。
bool rbFileSupportsAgentThumb(String fileName) {
  final ext = p.extension(fileName).toLowerCase();
  if (ext.isEmpty) {
    return false;
  }
  const images = {'.jpg', '.jpeg', '.jfif', '.png', '.gif', '.webp', '.avif', '.jxl', '.bmp', '.heic', '.heif', '.hif', '.tif', '.tiff'};
  if (images.contains(ext)) {
    return true;
  }
  return rbFileIsVideo(fileName);
}

/// 无缩略图时按扩展名选择 Material 图标（非媒体文件占位）。
IconData rbFileTypeIcon(String fileName) {
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
      return Icons.image_outlined;
    case '.mp4':
    case '.mov':
    case '.avi':
    case '.mkv':
    case '.webm':
    case '.m4v':
    case '.flv':
    case '.m2ts':
    case '.mts':
      return Icons.videocam_outlined;
    case '.mp3':
    case '.wav':
    case '.flac':
    case '.aac':
    case '.m4a':
    case '.ogg':
    case '.wma':
      return Icons.audiotrack_outlined;
    case '.zip':
    case '.rar':
    case '.7z':
    case '.tar':
    case '.gz':
    case '.bz2':
    case '.xz':
      return Icons.folder_zip_outlined;
    case '.pdf':
      return Icons.picture_as_pdf_outlined;
    case '.doc':
    case '.docx':
    case '.rtf':
      return Icons.description_outlined;
    case '.xls':
    case '.xlsx':
    case '.csv':
      return Icons.table_chart_outlined;
    case '.ppt':
    case '.pptx':
      return Icons.slideshow_outlined;
    case '.txt':
    case '.md':
    case '.json':
    case '.xml':
    case '.yaml':
    case '.yml':
    case '.toml':
    case '.ini':
    case '.log':
    case '.conf':
      return Icons.article_outlined;
    case '.dart':
    case '.go':
    case '.rs':
    case '.java':
    case '.js':
    case '.jsx':
    case '.tsx':
    case '.py':
    case '.cpp':
    case '.c':
    case '.h':
    case '.swift':
    case '.kt':
    case '.rb':
    case '.php':
    case '.sql':
    case '.sh':
    case '.proto':
      return Icons.code_outlined;
    case '.ts':
      return rbFileIsLikelyMpegTsSegment(fileName) ? Icons.videocam_outlined : Icons.code_outlined;
    case '.apk':
      return Icons.android_outlined;
    case '.exe':
    case '.dmg':
    case '.deb':
    case '.rpm':
    case '.msi':
      return Icons.apps_outlined;
    default:
      return Icons.insert_drive_file_outlined;
  }
}

Widget rbFileTypeIconWidget(String fileName, {double size = 36, Color? color}) {
  return Icon(rbFileTypeIcon(fileName), size: size, color: color);
}
