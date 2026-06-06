import 'dart:ffi';
import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:app/gen/pb/remotebrowse/browse.service.pb.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_ice_ffi_io.dart';
import 'package:ffi/ffi.dart';

class RbIceGrpcBridge {
  RbIceGrpcBridge(this._h);

  final Pointer<Void> _h;

  void close() {
    if (RbIceFfi.available) {
      RbIceFfi.grpcClose();
    }
    RbIceFfi.viewerClose(_h);
  }

  static RbIceGrpcBridge? tryOpen(Pointer<Void> h) {
    if (!RbIceFfi.available || h == nullptr) {
      return null;
    }
    if (RbIceFfi.grpcOpen(h) != 0) {
      RbIceFfi.viewerClose(h);
      return null;
    }
    return RbIceGrpcBridge(h);
  }

  Future<ListFilesResponse> listFiles(String root) async {
    final bytes = await _callUtf8(root, RbIceFfi.grpcListFiles);
    return ListFilesResponse.fromBuffer(bytes);
  }

  Future<ListMediaResponse> listMedia(String root, {MediaListCursor? cursor, int pageSize = 64}) async {
    final bytes = await _callListMedia(root, cursor ?? MediaListCursor(), pageSize);
    return ListMediaResponse.fromBuffer(bytes);
  }

  Future<ThumbnailBatchResponse> thumbBatch(List<String> paths, int maxEdge) async {
    final bytes = await _callPaths(paths.join('\n'), maxEdge, RbIceFfi.grpcThumbBatch);
    return ThumbnailBatchResponse.fromBuffer(bytes);
  }

  Future<ReadFileResponse> readFile(String path) async {
    final bytes = await _callRead(path, RbIceFfi.grpcReadFile);
    return ReadFileResponse.fromBuffer(bytes);
  }

  Future<ReadFileResponse> readFileRange(String path, {required int offset, required int length}) async {
    final bytes = await _callReadRange(path, offset, length);
    return ReadFileResponse.fromBuffer(bytes);
  }

  Future<DeleteFileResponse> deleteFile(String path) async {
    final bytes = await _callUtf8(path, RbIceFfi.grpcDeleteFile);
    return DeleteFileResponse.fromBuffer(bytes);
  }

  Future<Uint8List> _callUtf8(
    String arg,
    int Function(Pointer<Void>, Pointer<Utf8>, Pointer<Uint8>, int, Pointer<IntPtr>) fn,
  ) async {
    final p = arg.toNativeUtf8();
    final buf = calloc<Uint8>(8 << 20);
    final outLen = calloc<IntPtr>();
    try {
      final rc = fn(_h, p, buf, 8 << 20, outLen);
      if (rc != 0) {
        throw StateError('ice grpc failed ($rc)');
      }
      return Uint8List.fromList(buf.asTypedList(outLen.value));
    } finally {
      calloc.free(p);
      calloc.free(outLen);
      calloc.free(buf);
    }
  }

  Future<Uint8List> _callListMedia(String root, MediaListCursor cursor, int pageSize) async {
    final pRoot = root.toNativeUtf8();
    final cursorBytes = rbMediaCursorEmpty(cursor) ? Uint8List(0) : cursor.writeToBuffer();
    final pCursor = calloc<Uint8>(cursorBytes.length);
    final buf = calloc<Uint8>(8 << 20);
    final outLen = calloc<IntPtr>();
    try {
      if (cursorBytes.isNotEmpty) {
        pCursor.asTypedList(cursorBytes.length).setAll(0, cursorBytes);
      }
      final rc = RbIceFfi.grpcListMedia(_h, pRoot, pCursor, cursorBytes.length, pageSize, buf, 8 << 20, outLen);
      if (rc != 0) {
        throw StateError('ice grpc list media failed ($rc)');
      }
      return Uint8List.fromList(buf.asTypedList(outLen.value));
    } finally {
      calloc.free(pRoot);
      calloc.free(pCursor);
      calloc.free(outLen);
      calloc.free(buf);
    }
  }

  Future<Uint8List> _callPaths(
    String pathsNl,
    int maxEdge,
    int Function(Pointer<Void>, Pointer<Utf8>, int, Pointer<Uint8>, int, Pointer<IntPtr>) fn,
  ) async {
    final p = pathsNl.toNativeUtf8();
    final buf = calloc<Uint8>(8 << 20);
    final outLen = calloc<IntPtr>();
    try {
      final rc = fn(_h, p, maxEdge, buf, 8 << 20, outLen);
      if (rc != 0) {
        throw StateError('ice grpc thumb failed ($rc)');
      }
      return Uint8List.fromList(buf.asTypedList(outLen.value));
    } finally {
      calloc.free(p);
      calloc.free(outLen);
      calloc.free(buf);
    }
  }

  Future<Uint8List> _callReadRange(String path, int offset, int length) async {
    final p = path.toNativeUtf8();
    final buf = calloc<Uint8>(8 << 20);
    final outLen = calloc<IntPtr>();
    try {
      final rc = RbIceFfi.grpcReadFileRange(_h, p, offset, length, buf, 8 << 20, outLen);
      if (rc != 0) {
        throw StateError('ice grpc read range failed ($rc)');
      }
      return Uint8List.fromList(buf.asTypedList(outLen.value));
    } finally {
      calloc.free(p);
      calloc.free(outLen);
      calloc.free(buf);
    }
  }

  Future<Uint8List> _callRead(
    String path,
    int Function(Pointer<Void>, Pointer<Utf8>, int, Pointer<Uint8>, int, Pointer<IntPtr>) fn,
  ) async {
    final p = path.toNativeUtf8();
    final buf = calloc<Uint8>(8 << 20);
    final outLen = calloc<IntPtr>();
    try {
      final rc = fn(_h, p, 32 << 20, buf, 8 << 20, outLen);
      if (rc != 0) {
        throw StateError('ice grpc read failed ($rc)');
      }
      return Uint8List.fromList(buf.asTypedList(outLen.value));
    } finally {
      calloc.free(p);
      calloc.free(outLen);
      calloc.free(buf);
    }
  }
}
