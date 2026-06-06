import 'dart:convert';

import 'package:crypto/crypto.dart';

/// 与 Agent `media_token_for_room` 一致。
String rbMediaTokenForRoom(String room) {
  final r = room.trim();
  if (r.isEmpty) {
    return '';
  }
  return sha256.convert(utf8.encode('rb-media:$r')).toString().substring(0, 32);
}

String rbMediaHttpHost(String host) {
  final h = host.trim();
  if (h.contains(':') && !h.startsWith('[')) {
    return '[$h]';
  }
  return h;
}

/// Agent HTTP Range 原画播放（内网优先；无媒体端口时 Viewer 走 loopback+ICE 代理）。
Uri? rbAgentMediaPlayUri({
  required String host,
  required int port,
  required String relPath,
  required String room,
  int fileBaseOffset = 0,
  int agentImagePreviewMode = 0,
}) {
  final path = relPath.trim();
  if (path.isEmpty || host.isEmpty || port <= 0) {
    return null;
  }
  final q = <String, String>{'path': path};
  final tok = rbMediaTokenForRoom(room);
  if (tok.isNotEmpty) {
    q['t'] = tok;
  }
  if (agentImagePreviewMode == 1) {
    q['preview'] = '1';
  } else if (agentImagePreviewMode == 2) {
    q['preview'] = '2';
  }
  if (fileBaseOffset > 0) {
    q['off'] = '$fileBaseOffset';
  }
  return Uri(
    scheme: 'http',
    host: rbMediaHttpHost(host),
    port: port,
    path: '/rb/v1/media',
    queryParameters: q,
  );
}
