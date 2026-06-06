import 'package:flutter/material.dart';

/// 远端 Agent 操作系统（与 [rbPlatformName] / Rust `std::env::consts::OS` 对齐）。
enum RbRemoteOs { windows, macos, linux, android, ios, unknown }

RbRemoteOs rbParseRemoteOs(String? raw) {
  final s = raw?.trim().toLowerCase() ?? '';
  if (s.isEmpty) {
    return RbRemoteOs.unknown;
  }
  if (s.contains('win')) {
    return RbRemoteOs.windows;
  }
  if (s.contains('mac') || s == 'darwin' || s == 'osx') {
    return RbRemoteOs.macos;
  }
  if (s.contains('linux') || s.contains('android')) {
    return s.contains('android') ? RbRemoteOs.android : RbRemoteOs.linux;
  }
  if (s.contains('ios') || s == 'iphone') {
    return RbRemoteOs.ios;
  }
  return RbRemoteOs.unknown;
}

/// Windows canonicalize / verbatim 前缀（`\\?\` 或 `\`→`/` 后的 `//?/`）。
String rbStripVerbatimPrefix(String path) {
  var s = path.trim();
  if (s.startsWith(r'\\?\')) {
    s = s.substring(4);
    if (s.startsWith(r'UNC\')) {
      s = r'\\${s.substring(4)}';
    }
  } else if (s.startsWith('//?/')) {
    s = s.substring(4);
  } else if (s.startsWith('/?/')) {
    s = s.substring(3);
  }
  return s;
}

/// 远端目录 + 文件名（统一 `/`，兼容 `D:/...` 绝对路径）。
String rbJoinRemotePath(String dir, String name) {
  final d = rbNormRemotePath(dir);
  final n = name.trim().replaceAll('\\', '/');
  if (n.isEmpty) {
    return d;
  }
  if (rbIsWindowsRemotePath(n)) {
    return n;
  }
  if (d.isEmpty || d == '/') {
    return n;
  }
  return '$d/$n';
}

/// Unix 远端根 `/`（空串不是根）。
bool rbIsUnixRemoteRoot(String path) => rbNormRemotePath(path) == '/';

/// Agent `resolved_root_path`：空保持空，避免误变成 `/` 导致 Windows 相对路径丢盘符。
String rbNormResolvedRootPath(String raw) {
  final t = raw.trim();
  if (t.isEmpty) {
    return '';
  }
  return rbNormRemotePath(t);
}

/// 路径规范化：统一 `/`、去 verbatim 前缀与尾斜杠。wire 路径以 Agent `resolved_root_path`（如 `D:/…`）为准，不做终端形态转换。
String rbNormRemotePath(String path) {
  var t = rbStripVerbatimPrefix(path).trim().replaceAll('\\', '/');
  if (t.isEmpty) {
    return '';
  }
  if (t.length > 1 && t.endsWith('/')) {
    t = t.substring(0, t.length - 1);
  }
  return t;
}

bool rbIsWindowsRemotePath(String path) => RegExp(r'^[A-Za-z]:').hasMatch(rbNormRemotePath(path));

/// Viewer 连接配置起始路径：Windows 远端仅保留 `C:/…`，丢弃 Unix 形态旧值。
String rbSanitizeViewerListPath(String path, {String? agentPlatform}) {
  final t = path.trim();
  if (t.isEmpty || t == '.') {
    return '';
  }
  final plat = agentPlatform?.trim() ?? '';
  final winAgent = rbParseRemoteOs(plat) == RbRemoteOs.windows || rbIsWindowsRemotePath(t);
  if (winAgent) {
    return rbIsWindowsRemotePath(t) ? rbNormRemotePath(t) : '';
  }
  return rbNormRemotePath(t);
}

bool rbRemotePathsEqual(String a, String b) {
  final pa = rbNormRemotePath(a);
  final pb = rbNormRemotePath(b);
  if (pa == pb) {
    return true;
  }
  if (rbIsWindowsRemotePath(pa) || rbIsWindowsRemotePath(pb)) {
    return pa.toLowerCase() == pb.toLowerCase();
  }
  return false;
}

/// 远端路径是否在浏览沙箱 [root] 下（Windows 盘符路径大小写不敏感）。
bool rbIsUnderRemoteRoot(String path, String root) {
  final p = rbNormRemotePath(path);
  final r = rbNormRemotePath(root);
  if (rbRemotePathsEqual(p, r)) {
    return true;
  }
  final win = rbIsWindowsRemotePath(r) || rbIsWindowsRemotePath(p);
  if (win) {
    return p.toLowerCase().startsWith('${r.toLowerCase()}/');
  }
  return p.startsWith('$r/');
}

/// 绝对路径按层级拆分（列表/预览路径栏共用）。
List<({String label, String path})> rbRemotePathSegments(String absPath) {
  final text = rbNormRemotePath(absPath);
  if (text == '/') {
    return [(label: '/', path: '/')];
  }
  final isWinDrive = text.length >= 2 && text[1] == ':';
  final parts = text.split('/').where((s) => s.isNotEmpty).toList();
  if (parts.isEmpty) {
    return [(label: text, path: text)];
  }
  final out = <({String label, String path})>[];
  if (isWinDrive) {
    var acc = parts[0];
    out.add((label: parts[0], path: parts[0]));
    for (var i = 1; i < parts.length; i++) {
      acc = '$acc/${parts[i]}';
      out.add((label: parts[i], path: acc));
    }
    return out;
  }
  var acc = '';
  for (final part in parts) {
    acc = acc.isEmpty ? '/$part' : '$acc/$part';
    out.add((label: part, path: acc));
  }
  return out;
}

/// 绝对路径的父目录；已在根时返回 null。
String? rbRemotePathParent(String absPath) {
  final cur = rbNormRemotePath(absPath);
  if (cur == '/' || cur.isEmpty) {
    return null;
  }
  if (RegExp(r'^[A-Za-z]:$').hasMatch(cur) || RegExp(r'^[A-Za-z]:/$').hasMatch(cur)) {
    return null;
  }
  final i = cur.lastIndexOf('/');
  if (i < 0) {
    return null;
  }
  if (i == 0) {
    return '/';
  }
  if (RegExp(r'^[A-Za-z]:').hasMatch(cur) && i <= 3) {
    return cur.substring(0, 2);
  }
  return cur.substring(0, i);
}

String? rbInferOsFromPath(String path) {
  final p = rbNormRemotePath(path);
  if (p.isEmpty) {
    return null;
  }
  if (RegExp(r'^[A-Za-z]:/').hasMatch(p)) {
    return 'windows';
  }
  if (p.startsWith('/Users/') || p.startsWith('/private/var/')) {
    return 'macos';
  }
  // `/home/` 可能是 MSYS 输入形态，Windows Agent 解析后应为 `D:/…/home/…`，勿据此判 linux。
  if (p.startsWith('/home/') || p.startsWith('/root/')) {
    return null;
  }
  return null;
}

class RbRemoteOsAvatar extends StatelessWidget {
  const RbRemoteOsAvatar({super.key, required this.os, this.radius = 22});

  final RbRemoteOs os;
  final double radius;

  @override
  Widget build(BuildContext context) {
    final (bg, fg, icon) = switch (os) {
      RbRemoteOs.windows => (const Color(0xFFE3F2FD), const Color(0xFF0078D4), Icons.desktop_windows_rounded),
      RbRemoteOs.macos => (const Color(0xFFF5F5F5), const Color(0xFF424242), Icons.laptop_mac_rounded),
      RbRemoteOs.linux => (const Color(0xFFFFF3E0), const Color(0xFFE95420), Icons.terminal_rounded),
      RbRemoteOs.android => (const Color(0xFFE8F5E9), const Color(0xFF3DDC84), Icons.phone_android_rounded),
      RbRemoteOs.ios => (const Color(0xFFF3E5F5), const Color(0xFF7B1FA2), Icons.phone_iphone_rounded),
      RbRemoteOs.unknown => (
          Theme.of(context).colorScheme.primaryContainer,
          Theme.of(context).colorScheme.onPrimaryContainer,
          Icons.computer_rounded,
        ),
    };
    return CircleAvatar(
      radius: radius,
      backgroundColor: bg,
      child: Icon(icon, size: radius, color: fg),
    );
  }
}
