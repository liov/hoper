import 'dart:io';

import 'package:file_picker/file_picker.dart';

/// Agent 沙箱默认根：用户主目录。
class RbAgentPaths {
  static String defaultRoot() {
    if (Platform.isWindows) {
      final profile = _windowsUserProfileFromEnv();
      if (profile != null && profile.isNotEmpty) {
        return profile;
      }
      throw StateError('无法解析 Windows 用户目录（USERPROFILE / PROFILE / HOMEDRIVE+HOMEPATH）');
    }
    final home = Platform.environment['HOME']?.trim();
    if (home != null && home.isNotEmpty) {
      return home;
    }
    return Directory.current.path;
  }

  /// 仅读 Windows 环境变量：`USERPROFILE`、`PROFILE`、`HOMEDRIVE`+`HOMEPATH`。
  static String? _windowsUserProfileFromEnv() {
    for (final key in ['USERPROFILE', 'PROFILE']) {
      final v = Platform.environment[key]?.trim();
      if (v != null && v.isNotEmpty) {
        return _normalizeWindowsProfileEnv(v) ?? v;
      }
    }
    final drive = Platform.environment['HOMEDRIVE']?.trim() ?? '';
    final homePath = Platform.environment['HOMEPATH']?.trim() ?? '';
    if (drive.isEmpty || homePath.isEmpty) {
      return null;
    }
    final combined = '$drive$homePath';
    return _normalizeWindowsProfileEnv(combined) ?? combined;
  }

  /// MSYS 可能把 USERPROFILE 设为 `/c/Users/…` → `C:\Users\…`
  static String? _normalizeWindowsProfileEnv(String raw) {
    final t = raw.trim().replaceAll('\\', '/');
    if (t.length >= 4 && t.startsWith('/')) {
      final drive = t[1];
      if (_isDriveLetter(drive) && t[2] == '/') {
        return '${drive.toUpperCase()}:${t.substring(2)}'.replaceAll('/', '\\');
      }
    }
    return null;
  }

  static bool _isDriveLetter(String c) {
    if (c.length != 1) {
      return false;
    }
    final u = c.codeUnitAt(0);
    return (u >= 0x41 && u <= 0x5A) || (u >= 0x61 && u <= 0x7A);
  }

  static List<String> windowsDrives() {
    if (!Platform.isWindows) {
      return [];
    }
    final out = <String>[];
    for (var c = 65; c <= 90; c++) {
      final root = '${String.fromCharCode(c)}:\\';
      if (Directory(root).existsSync()) {
        out.add(root);
      }
    }
    return out;
  }

  static Future<String?> pickDirectory({String? initial}) async {
    return FilePicker.getDirectoryPath(initialDirectory: initial);
  }
}
