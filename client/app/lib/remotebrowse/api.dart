import 'dart:convert';
import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:dio/dio.dart';

class RbListFilesResult {
  RbListFilesResult({required this.entries, required this.resolvedRootPath, this.hiddenSkipped = 0});

  final List<RbFileEntry> entries;
  final String resolvedRootPath;
  final int hiddenSkipped;
}

class RbListMediaResult {
  RbListMediaResult({
    required this.items,
    required this.nextCursor,
    required this.done,
    required this.resolvedRootPath,
    this.scannedDirs = 0,
  });

  final List<RbMediaItem> items;
  final MediaListCursor nextCursor;
  final bool done;
  final String resolvedRootPath;
  final int scannedDirs;
}

bool rbMediaCursorEmpty(MediaListCursor c) => c.stack.isEmpty && c.pending.isEmpty;

class RbMediaItem {
  const RbMediaItem({required this.relPath, required this.entry});

  final String relPath;
  final RbFileEntry entry;
}

/// 与 proto `FileEntry.flags` 一致：bit0 = 目录；bit1 = Live / Motion Photo
const rbFileFlagDirectory = 1;
const rbFileFlagMotionPhoto = 2;

class RbFileEntry {
  RbFileEntry({
    required this.id,
    required this.name,
    required this.size,
    this.mtimeUnixMs = 0,
    this.thumbHash = '',
    this.isDirectory = false,
    this.durationMs = 0,
    this.isMotionPhoto = false,
    this.motionOffset = 0,
    this.motionLength = 0,
    this.motionCompanion = '',
  });

  final String id;
  final String name;
  final int size;
  final int mtimeUnixMs;
  final String thumbHash;
  final bool isDirectory;
  final int durationMs;
  final bool isMotionPhoto;
  final int motionOffset;
  final int motionLength;
  final String motionCompanion;

  /// 列表展示：修改时间从新到旧，同时间按文件名。
  static int compareByMtimeDesc(RbFileEntry a, RbFileEntry b) {
    final byTime = b.mtimeUnixMs.compareTo(a.mtimeUnixMs);
    if (byTime != 0) {
      return byTime;
    }
    return a.name.compareTo(b.name);
  }

  factory RbFileEntry.fromJson(Map<String, dynamic> json) {
    final flags = (json['flags'] as num?)?.toInt() ?? 0;
    return RbFileEntry(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      size: (json['size'] as num?)?.toInt() ?? 0,
      mtimeUnixMs: (json['mtimeUnixMs'] as num?)?.toInt() ?? (json['mtime_unix_ms'] as num?)?.toInt() ?? 0,
      thumbHash: json['thumbHash'] as String? ?? json['thumb_hash'] as String? ?? '',
      isDirectory: (flags & rbFileFlagDirectory) != 0,
      durationMs: (json['durationMs'] as num?)?.toInt() ?? (json['duration_ms'] as num?)?.toInt() ?? 0,
      isMotionPhoto: (flags & rbFileFlagMotionPhoto) != 0,
      motionOffset: (json['motionOffset'] as num?)?.toInt() ?? (json['motion_offset'] as num?)?.toInt() ?? 0,
      motionLength: (json['motionLength'] as num?)?.toInt() ?? (json['motion_length'] as num?)?.toInt() ?? 0,
      motionCompanion: json['motionCompanion'] as String? ?? json['motion_companion'] as String? ?? '',
    );
  }
}

void rbSortEntriesByMtimeDesc(List<RbFileEntry> entries) {
  entries.sort(RbFileEntry.compareByMtimeDesc);
}

class RemoteBrowseApi {
  RemoteBrowseApi({String? baseUrl, Dio? client})
      : _baseUrl = baseUrl ?? rbDebugBaseUrl,
        _dio = client ?? Dio(BaseOptions(baseUrl: baseUrl ?? rbDebugBaseUrl, connectTimeout: const Duration(seconds: 8)));

  static const rbDebugBaseUrl = 'https://api.hoper.xyz';

  final String _baseUrl;
  final Dio _dio;

  Future<Map<String, dynamic>> health() async {
    final res = await _dio.get<Map<String, dynamic>>('/rb/health', options: Options(responseType: ResponseType.json));
    return res.data ?? {};
  }

  Future<List<RbFileEntry>> listFiles(String path) async {
    final res = await _dio.get<String>('/rb/v1/list', queryParameters: {'path': path}, options: Options(responseType: ResponseType.plain));
    final body = res.data;
    if (body == null || body.isEmpty) {
      return [];
    }
    final map = jsonDecode(body) as Map<String, dynamic>;
    final entries = map['entries'] as List<dynamic>? ?? [];
    return entries.map((e) => RbFileEntry.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<Uint8List> fetchThumb(String path, {int maxEdge = 256, String? hash}) async {
    final res = await _dio.get<List<int>>(
      '/rb/v1/thumb',
      queryParameters: {'path': path, 'max_edge': maxEdge, if (hash != null && hash.isNotEmpty) 'hash': hash},
      options: Options(responseType: ResponseType.bytes, headers: {'accept': 'application/json, application/x-protobuf, */*'}),
    );
    final data = res.data;
    if (data == null) {
      return Uint8List(0);
    }
    return Uint8List.fromList(data);
  }

  String signalWsUrl() => signalWsUrlFrom(_baseUrl);
}

String defaultSignalWs() => signalWsUrlFrom(RemoteBrowseApi.rbDebugBaseUrl);

String signalWsUrlFrom(String base) => parseSignalWsUri(base).toString();

/// 规范化信令地址：补 ws/wss、补 `/rb/signal`、支持 `host:port` 简写。
Uri parseSignalWsUri(String raw) {
  var t = raw.trim();
  if (t.isEmpty) {
    throw ArgumentError('信令地址不能为空');
  }
  if (!t.contains('://')) {
    t = 'ws://$t';
  }
  Uri uri;
  try {
    uri = Uri.parse(t);
  } on FormatException {
    throw ArgumentError('信令地址格式无效: $raw');
  }
  if (uri.host.isEmpty) {
    throw ArgumentError('信令地址缺少主机: $raw');
  }
  var scheme = uri.scheme;
  if (scheme == 'http') {
    scheme = 'ws';
  } else if (scheme == 'https') {
    scheme = 'wss';
  } else if (scheme != 'ws' && scheme != 'wss') {
    throw ArgumentError('信令须为 ws/wss 或 http(s): $raw');
  }
  var path = uri.path;
  if (!path.endsWith('/rb/signal')) {
    path = '${path.replaceAll(RegExp(r'/+$'), '')}/rb/signal';
  }
  uri = Uri(scheme: scheme, host: uri.host, port: uri.hasPort ? uri.port : null, path: path);
  if (uri.hasPort && (uri.port < 1 || uri.port > 65535)) {
    throw ArgumentError('信令端口无效(${uri.port})，请检查是否为 8079 等常用端口');
  }
  return uri;
}

/// 由信令 WebSocket 地址推导 daemon 健康检查 URL（`GET /rb/health`）。
Uri signalHealthUri(Uri signalWs) {
  final scheme = signalWs.scheme == 'wss' ? 'https' : 'http';
  final port = signalWs.hasPort ? signalWs.port : (scheme == 'https' ? 443 : 80);
  return Uri(scheme: scheme, host: signalWs.host, port: port, path: '/rb/health');
}

/// 连接前探测信令服务是否可达（非 WebSocket，仅 HTTP health）。
Future<void> probeSignalHealth(Uri signalWs, {Duration timeout = const Duration(seconds: 5)}) async {
  await _getSignalHealth(signalWs, timeout: timeout);
}

/// 查询房间内在线 Agent 的平台（daemon `GET /rb/health?room=`）。
Future<String?> probeAgentPlatform(Uri signalWs, String room, {Duration timeout = const Duration(seconds: 4)}) async {
  if (room.trim().isEmpty) {
    return null;
  }
  final data = await _getSignalHealth(signalWs, room: room.trim(), timeout: timeout);
  final plat = data['agentPlatform'] as String? ?? data['agent_platform'] as String?;
  final s = plat?.trim() ?? '';
  return s.isEmpty ? null : s;
}

Future<Map<String, dynamic>> _getSignalHealth(Uri signalWs, {String? room, Duration timeout = const Duration(seconds: 5)}) async {
  var uri = signalHealthUri(signalWs);
  if (room != null && room.isNotEmpty) {
    uri = uri.replace(queryParameters: {'room': room});
  }
  final dio = Dio(BaseOptions(connectTimeout: timeout, receiveTimeout: timeout));
  try {
    final res = await dio.getUri<Map<String, dynamic>>(uri);
    if (res.statusCode != 200) {
      throw StateError('HTTP ${res.statusCode}');
    }
    return res.data ?? {};
  } on DioException catch (e) {
    throw StateError('无法访问 $uri：${e.message ?? e.type}');
  }
}
