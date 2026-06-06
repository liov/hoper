import 'dart:async';
import 'dart:io';

import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_media_http.dart';

/// 与 Agent `FRAGMENT_MS` 一致。
const rbTranscodeFragmentMs = 2000;

/// 转码播放优先连续 MPEG-TS 流；失败时回退 HLS 分片。
const rbPreferTranscodeStream = true;

/// HLS 转码视频编码（Agent 侧，不用 H.264）。
class RbTranscodeVcodec {
  static const hevc = 'hevc';
  static const av1 = 'av1';

  /// 按源文件编码选择转码输出（VP8/VP9 等走原画直传，不转 H.264）。
  static String pickForMeta(String? videoCodec) {
    final v = videoCodec?.toLowerCase() ?? '';
    if (v.contains('av1')) {
      return av1;
    }
    return hevc;
  }
}

/// 远程视频转码档位（Agent 按需分片 HLS）；[source] 为原文件 HTTP Range 直传。
class RbTranscodePreset {
  const RbTranscodePreset({required this.id, required this.label, this.agentOnly = false});

  final String id;
  final String label;
  /// 仅 Agent HTTP 可转码；P2P 回退代理时不可用。
  final bool agentOnly;

  static const source = RbTranscodePreset(id: 'source', label: '原画');
  static const p1080 = RbTranscodePreset(id: '1080', label: '1080p', agentOnly: true);
  static const p720 = RbTranscodePreset(id: '720', label: '720p', agentOnly: true);
  static const p480 = RbTranscodePreset(id: '480', label: '480p', agentOnly: true);
  static const p360 = RbTranscodePreset(id: '360', label: '360p', agentOnly: true);

  static const all = [source, p1080, p720, p480, p360];

  bool get isSource => id == source.id;
}

Uri? rbAgentTranscodeHlsUri({
  required String host,
  required int port,
  required String relPath,
  required String room,
  required String presetId,
  String vcodec = RbTranscodeVcodec.hevc,
  int fileBaseOffset = 0,
  int beginMs = 0,
}) {
  if (presetId == RbTranscodePreset.source.id || fileBaseOffset > 0) {
    return null;
  }
  final path = relPath.trim();
  if (path.isEmpty || host.isEmpty || port <= 0) {
    return null;
  }
  final q = <String, String>{'path': path, 'preset': presetId, 'vcodec': vcodec};
  if (beginMs > 0) {
    q['begin_ms'] = '$beginMs';
  }
  final tok = rbMediaTokenForRoom(room);
  if (tok.isNotEmpty) {
    q['t'] = tok;
  }
  return Uri(
    scheme: 'http',
    host: rbMediaHttpHost(host),
    port: port,
    path: '/rb/v1/transcode/index.m3u8',
    queryParameters: q,
  );
}

/// 连续转码流（无 HLS 分片边界）；[startMs] 为起播位置（2s 对齐）。
Uri? rbAgentTranscodeStreamUri({
  required String host,
  required int port,
  required String relPath,
  required String room,
  required String presetId,
  String vcodec = RbTranscodeVcodec.hevc,
  int fileBaseOffset = 0,
  int startMs = 0,
}) {
  if (presetId == RbTranscodePreset.source.id || fileBaseOffset > 0) {
    return null;
  }
  final path = relPath.trim();
  if (path.isEmpty || host.isEmpty || port <= 0) {
    return null;
  }
  final q = <String, String>{'path': path, 'preset': presetId, 'vcodec': vcodec};
  if (startMs > 0) {
    q['start_ms'] = '$startMs';
  }
  final tok = rbMediaTokenForRoom(room);
  if (tok.isNotEmpty) {
    q['t'] = tok;
  }
  return Uri(
    scheme: 'http',
    host: rbMediaHttpHost(host),
    port: port,
    path: '/rb/v1/transcode/stream.ts',
    queryParameters: q,
  );
}

Uri rbTranscodeFragmentUri(Uri indexUri, int startMs, {int prefetchAhead = 0}) {
  final q = Map<String, String>.from(indexUri.queryParameters);
  q['start_ms'] = '$startMs';
  if (prefetchAhead > 0) {
    q['prefetch'] = '$prefetchAhead';
  }
  return indexUri.replace(path: '/rb/v1/transcode/fragment.ts', queryParameters: q);
}

/// 切换码率前预热：拉 m3u8，并按需拉 [fragmentCount] 个分片（从 [warmStartMs] 对齐的 2s 边界起）。
Future<bool> rbWarmTranscodePlaylist(
  Uri indexUri, {
  Duration timeout = const Duration(seconds: 90),
  int fragmentCount = 8,
  int retries = 2,
  int warmStartMs = 0,
}) async {
  if (fragmentCount < 0) {
    fragmentCount = 0;
  }
  for (var attempt = 0; attempt <= retries; attempt++) {
    if (await _warmTranscodeOnce(indexUri, timeout: timeout, fragmentCount: fragmentCount, warmStartMs: warmStartMs)) {
      return true;
    }
    if (attempt < retries) {
      await Future<void>.delayed(Duration(milliseconds: 400 * (attempt + 1)));
    }
  }
  return false;
}

Future<bool> _warmTranscodeOnce(Uri indexUri, {required Duration timeout, required int fragmentCount, int warmStartMs = 0}) async {
  final client = HttpClient();
  final baseMs = warmStartMs > 0 ? (warmStartMs ~/ rbTranscodeFragmentMs) * rbTranscodeFragmentMs : 0;
  try {
    var req = await client.getUrl(indexUri).timeout(timeout);
    var res = await req.close().timeout(timeout);
    if (res.statusCode != HttpStatus.ok) {
      rbLog.warning('transcode m3u8 http ${res.statusCode} $indexUri');
      return false;
    }
    await res.drain().timeout(timeout);
    if (fragmentCount <= 0) {
      return true;
    }
    final frag = rbTranscodeFragmentUri(indexUri, baseMs, prefetchAhead: fragmentCount);
    req = await client.getUrl(frag).timeout(timeout);
    res = await req.close().timeout(timeout);
    if (res.statusCode != HttpStatus.ok) {
      rbLog.warning('transcode frag http ${res.statusCode} prefetch=$fragmentCount $frag');
      return false;
    }
    await res.drain().timeout(timeout);
    return true;
  } catch (e, st) {
    rbLog.warning('transcode warm fail: $e', e, st);
    return false;
  } finally {
    client.close(force: true);
  }
}

Uri rbTranscodeUriCacheBust(Uri indexUri) {
  final q = Map<String, String>.from(indexUri.queryParameters);
  q['_sw'] = '${DateTime.now().millisecondsSinceEpoch}';
  return indexUri.replace(queryParameters: q);
}

/// 播放中后台续热后续分片（不阻塞 UI）。
/// 预热连续流：拉取前 [maxBytes] 字节，让 Agent 先开始转码。
Future<bool> rbWarmTranscodeStream(Uri streamUri, {Duration timeout = const Duration(seconds: 60), int maxBytes = 384 << 10}) async {
  final client = HttpClient();
  try {
    final req = await client.getUrl(streamUri).timeout(timeout);
    final res = await req.close().timeout(timeout);
    if (res.statusCode != HttpStatus.ok) {
      rbLog.warning('transcode stream warm http ${res.statusCode} $streamUri');
      return false;
    }
    var got = 0;
    await for (final chunk in res.timeout(timeout)) {
      got += chunk.length;
      if (got >= maxBytes) {
        break;
      }
    }
    return got > 1024;
  } catch (e, st) {
    rbLog.warning('transcode stream warm fail: $e', e, st);
    return false;
  } finally {
    client.close(force: true);
  }
}

void rbWarmTranscodeAhead(Uri indexUri, int fromStartMs, {int fragmentCount = 10}) {
  if (fragmentCount <= 0) {
    return;
  }
  unawaited(_warmTranscodeOnce(indexUri, timeout: const Duration(seconds: 180), fragmentCount: fragmentCount, warmStartMs: fromStartMs));
}
