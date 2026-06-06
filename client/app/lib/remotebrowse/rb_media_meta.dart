import 'dart:convert';
import 'dart:io';

import 'package:app/remotebrowse/rb_media_http.dart';

/// Agent `/rb/v1/media/meta` 探测结果。
class RbMediaMeta {
  const RbMediaMeta({
    required this.durationMs,
    required this.directPlay,
    required this.videoCodec,
    required this.audioCodec,
  });

  final int durationMs;
  final bool directPlay;
  final String videoCodec;
  final String audioCodec;

  factory RbMediaMeta.fromJson(Map<String, dynamic> json) {
    return RbMediaMeta(
      durationMs: (json['duration_ms'] as num?)?.toInt() ?? 0,
      directPlay: json['direct_play'] == true,
      videoCodec: json['video_codec'] as String? ?? '',
      audioCodec: json['audio_codec'] as String? ?? '',
    );
  }
}

Uri? rbAgentMediaMetaUri({
  required String host,
  required int port,
  required String relPath,
  required String room,
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
  return Uri(
    scheme: 'http',
    host: rbMediaHttpHost(host),
    port: port,
    path: '/rb/v1/media/meta',
    queryParameters: q,
  );
}

Future<RbMediaMeta?> rbFetchAgentMediaMeta(Uri uri, {Duration timeout = const Duration(seconds: 12)}) async {
  final client = HttpClient();
  try {
    final req = await client.getUrl(uri).timeout(timeout);
    final res = await req.close().timeout(timeout);
    if (res.statusCode != HttpStatus.ok) {
      return null;
    }
    final body = await res.transform(utf8.decoder).join().timeout(timeout);
    final json = jsonDecode(body);
    if (json is! Map<String, dynamic>) {
      return null;
    }
    return RbMediaMeta.fromJson(json);
  } catch (_) {
    return null;
  } finally {
    client.close(force: true);
  }
}
