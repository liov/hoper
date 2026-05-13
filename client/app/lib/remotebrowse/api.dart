import 'dart:convert';
import 'dart:typed_data';

import 'package:dio/dio.dart';

const rbBaseUrl = 'https://api.hoper.xyz';

class RbFileEntry {
  RbFileEntry({required this.id, required this.name, required this.size, this.thumbHash = ''});

  final String id;
  final String name;
  final int size;
  final String thumbHash;

  factory RbFileEntry.fromJson(Map<String, dynamic> json) {
    return RbFileEntry(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      size: (json['size'] as num?)?.toInt() ?? 0,
      thumbHash: json['thumbHash'] as String? ?? json['thumb_hash'] as String? ?? '',
    );
  }
}

class RemoteBrowseApi {
  RemoteBrowseApi({Dio? client}) : _dio = client ?? Dio(BaseOptions(baseUrl: rbBaseUrl, connectTimeout: const Duration(seconds: 8)));

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

  String signalWsUrl() {
    final uri = Uri.parse(rbBaseUrl);
    final scheme = uri.scheme == 'https' ? 'wss' : 'ws';
    return Uri(scheme: scheme, host: uri.host, port: uri.hasPort ? uri.port : null, path: '/rb/signal').toString();
  }
}
