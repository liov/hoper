import 'dart:ffi';

import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:app/gen/pb/remotebrowse/browse.service.pb.dart';

class RbIceGrpcBridge {
  RbIceGrpcBridge(this._h);
  final Pointer<Void> _h;
  void close() {}
  static RbIceGrpcBridge? tryOpen(Pointer<Void> h) => null;
  Future<ListFilesResponse> listFiles(String root) async => throw UnsupportedError('ice ffi');
  Future<ListMediaResponse> listMedia(String root, {MediaListCursor? cursor, int pageSize = 64}) async =>
      throw UnsupportedError('ice ffi');
  Future<ThumbnailBatchResponse> thumbBatch(List<String> paths, int maxEdge) async =>
      throw UnsupportedError('ice ffi');
  Future<ReadFileResponse> readFile(String path) async => throw UnsupportedError('ice ffi');
  Future<ReadFileResponse> readFileRange(String path, {required int offset, required int length}) async =>
      throw UnsupportedError('ice ffi');
  Future<DeleteFileResponse> deleteFile(String path) async => throw UnsupportedError('ice ffi');
}
