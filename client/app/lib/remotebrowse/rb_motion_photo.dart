import 'package:app/remotebrowse/api.dart';
import 'package:path/path.dart' as p;

/// 动态照片配套视频的 wire 相对路径；无视频可播时返回 null。
String? rbMotionVideoRelPath(String stillRelPath, RbFileEntry entry) {
  if (!entry.isMotionPhoto) {
    return null;
  }
  if (entry.motionCompanion.isNotEmpty) {
    return p.posix.join(p.posix.dirname(stillRelPath), entry.motionCompanion);
  }
  if (entry.motionLength > 0) {
    return stillRelPath;
  }
  return null;
}

int rbMotionVideoBaseOffset(RbFileEntry entry) =>
    entry.motionCompanion.isEmpty ? entry.motionOffset : 0;

int rbMotionVideoTotalSize(RbFileEntry entry) {
  if (entry.motionLength > 0) {
    return entry.motionLength;
  }
  return 0;
}

String rbMotionVideoMime(RbFileEntry entry) {
  final name = entry.motionCompanion.isNotEmpty ? entry.motionCompanion : entry.name;
  if (name.toLowerCase().endsWith('.mov')) {
    return 'video/quicktime';
  }
  return 'video/mp4';
}
