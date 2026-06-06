import 'dart:ffi';
import 'dart:typed_data';

import 'package:app/remotebrowse/rb_ffi_lib.dart';
import 'package:ffi/ffi.dart';

typedef _ViewerNewNative = Pointer<Void> Function(Int32 timeoutMs);
typedef _ViewerNewDart = Pointer<Void> Function(int timeoutMs);
typedef _ViewerPushNative =
    Int32 Function(Pointer<Void> h, Pointer<Uint8> data, IntPtr len);
typedef _ViewerPushDart =
    int Function(Pointer<Void> h, Pointer<Uint8> data, int len);
typedef _ViewerPollOutNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Uint8> buf,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _ViewerPollOutDart =
    int Function(
      Pointer<Void> h,
      Pointer<Uint8> buf,
      int cap,
      Pointer<IntPtr> outLen,
    );
typedef _ViewerStateNative = Int32 Function(Pointer<Void> h);
typedef _ViewerStateDart = int Function(Pointer<Void> h);
typedef _ViewerWriteNative =
    Int32 Function(Pointer<Void> h, Uint8 typ, Pointer<Uint8> data, IntPtr len);
typedef _ViewerWriteDart =
    int Function(Pointer<Void> h, int typ, Pointer<Uint8> data, int len);
typedef _ViewerReadNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Uint8> buf,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _ViewerReadDart =
    int Function(
      Pointer<Void> h,
      Pointer<Uint8> buf,
      int cap,
      Pointer<IntPtr> outLen,
    );
typedef _ViewerCloseNative = Void Function(Pointer<Void> h);
typedef _ViewerCloseDart = void Function(Pointer<Void> h);
typedef _AgentRunNative =
    Int32 Function(
      Pointer<Utf8> signal,
      Pointer<Utf8> room,
      Pointer<Utf8> root,
      Int32 timeoutMs,
    );
typedef _AgentRunDart =
    int Function(
      Pointer<Utf8> signal,
      Pointer<Utf8> room,
      Pointer<Utf8> root,
      int timeoutMs,
    );
typedef _AgentRunningNative = Int32 Function();
typedef _AgentRunningDart = int Function();
typedef _IceGrpcOpenNative = Int32 Function(Pointer<Void> h);
typedef _IceGrpcOpenDart = int Function(Pointer<Void> h);
typedef _IceGrpcCloseNative = Void Function();
typedef _IceGrpcCloseDart = void Function();
typedef _IceGrpcListNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Utf8> root,
      Pointer<Uint8> out,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcListDart =
    int Function(
      Pointer<Void> h,
      Pointer<Utf8> root,
      Pointer<Uint8> out,
      int cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcListMediaNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Utf8> root,
      Pointer<Uint8> cursor,
      IntPtr cursorLen,
      Uint32 pageSize,
      Pointer<Uint8> out,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcListMediaDart =
    int Function(
      Pointer<Void> h,
      Pointer<Utf8> root,
      Pointer<Uint8> cursor,
      int cursorLen,
      int pageSize,
      Pointer<Uint8> out,
      int cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcThumbNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Utf8> pathsNl,
      Uint32 maxEdge,
      Pointer<Uint8> out,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcThumbDart =
    int Function(
      Pointer<Void> h,
      Pointer<Utf8> pathsNl,
      int maxEdge,
      Pointer<Uint8> out,
      int cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcReadNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Utf8> path,
      Uint32 maxBytes,
      Pointer<Uint8> out,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcReadDart =
    int Function(
      Pointer<Void> h,
      Pointer<Utf8> path,
      int maxBytes,
      Pointer<Uint8> out,
      int cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcReadRangeNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Utf8> path,
      Int64 offset,
      Uint32 length,
      Pointer<Uint8> out,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcReadRangeDart =
    int Function(
      Pointer<Void> h,
      Pointer<Utf8> path,
      int offset,
      int length,
      Pointer<Uint8> out,
      int cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcDelNative =
    Int32 Function(
      Pointer<Void> h,
      Pointer<Utf8> path,
      Pointer<Uint8> out,
      IntPtr cap,
      Pointer<IntPtr> outLen,
    );
typedef _IceGrpcDelDart =
    int Function(
      Pointer<Void> h,
      Pointer<Utf8> path,
      Pointer<Uint8> out,
      int cap,
      Pointer<IntPtr> outLen,
    );

class RbIceFfi {
  static DynamicLibrary? _lib;
  static late final _ViewerNewDart _viewerNew;
  static late final _ViewerPushDart _viewerPush;
  static late final _ViewerPollOutDart _viewerPollOut;
  static late final _ViewerStateDart _viewerState;
  static late final _ViewerWriteDart _viewerWrite;
  static late final _ViewerReadDart _viewerRead;
  static late final _ViewerCloseDart _viewerClose;
  static _AgentRunDart? _agentRun;
  static _AgentRunningDart? _agentRunning;
  static late final _IceGrpcOpenDart _grpcOpen;
  static late final _IceGrpcCloseDart _grpcClose;
  static late final _IceGrpcListDart _grpcList;
  static _IceGrpcListMediaDart? _grpcListMedia;
  static late final _IceGrpcThumbDart _grpcThumb;
  static late final _IceGrpcReadDart _grpcRead;
  static _IceGrpcReadRangeDart? _grpcReadRange;
  static late final _IceGrpcDelDart _grpcDel;
  static var _inited = false;
  static var _ready = false;

  static bool get available {
    _ensure();
    return _ready;
  }

  static bool get agentRunAvailable {
    _ensure();
    return _agentRun != null;
  }

  static void _ensure() {
    if (_inited) {
      return;
    }
    _inited = true;
    try {
      _lib = findRbStaticLibrary();
      _viewerNew = _lib!.lookupFunction<_ViewerNewNative, _ViewerNewDart>(
        'rb_ice_viewer_new',
      );
      _viewerPush = _lib!.lookupFunction<_ViewerPushNative, _ViewerPushDart>(
        'rb_ice_viewer_push',
      );
      _viewerPollOut = _lib!
          .lookupFunction<_ViewerPollOutNative, _ViewerPollOutDart>(
            'rb_ice_viewer_poll_out',
          );
      _viewerState = _lib!.lookupFunction<_ViewerStateNative, _ViewerStateDart>(
        'rb_ice_viewer_state',
      );
      _viewerWrite = _lib!.lookupFunction<_ViewerWriteNative, _ViewerWriteDart>(
        'rb_ice_viewer_write',
      );
      _viewerRead = _lib!.lookupFunction<_ViewerReadNative, _ViewerReadDart>(
        'rb_ice_viewer_read',
      );
      _viewerClose = _lib!.lookupFunction<_ViewerCloseNative, _ViewerCloseDart>(
        'rb_ice_viewer_close',
      );
      _ready = true;
      _grpcOpen = _lib!.lookupFunction<_IceGrpcOpenNative, _IceGrpcOpenDart>(
        'rb_ice_grpc_open',
      );
      _grpcClose = _lib!.lookupFunction<_IceGrpcCloseNative, _IceGrpcCloseDart>(
        'rb_ice_grpc_close',
      );
      _grpcList = _lib!.lookupFunction<_IceGrpcListNative, _IceGrpcListDart>(
        'rb_ice_grpc_list_files',
      );
      try {
        _grpcListMedia = _lib!.lookupFunction<_IceGrpcListMediaNative, _IceGrpcListMediaDart>(
          'rb_ice_grpc_list_media',
        );
      } catch (_) {
        _grpcListMedia = null;
      }
      _grpcThumb = _lib!.lookupFunction<_IceGrpcThumbNative, _IceGrpcThumbDart>(
        'rb_ice_grpc_thumb_batch',
      );
      _grpcRead = _lib!.lookupFunction<_IceGrpcReadNative, _IceGrpcReadDart>(
        'rb_ice_grpc_read_file',
      );
      try {
        _grpcReadRange = _lib!.lookupFunction<_IceGrpcReadRangeNative, _IceGrpcReadRangeDart>(
          'rb_ice_grpc_read_file_range',
        );
      } catch (_) {
        _grpcReadRange = null;
      }
      _grpcDel = _lib!.lookupFunction<_IceGrpcDelNative, _IceGrpcDelDart>(
        'rb_ice_grpc_delete_file',
      );
      try {
        _agentRun = _lib!.lookupFunction<_AgentRunNative, _AgentRunDart>(
          'rb_agent_run',
        );
        _agentRunning = _lib!
            .lookupFunction<_AgentRunningNative, _AgentRunningDart>(
              'rb_agent_running',
            );
      } catch (_) {
        _agentRun = null;
        _agentRunning = null;
      }
    } catch (_) {
      _lib = null;
      _ready = false;
    }
  }

  static Pointer<Void> viewerNew(int timeoutMs) {
    if (!available) {
      return nullptr;
    }
    return _viewerNew(timeoutMs);
  }

  static void viewerPush(Pointer<Void> h, Uint8List data) {
    if (!available || h == nullptr) {
      return;
    }
    final p = calloc<Uint8>(data.length);
    try {
      p.asTypedList(data.length).setAll(0, data);
      _viewerPush(h, p, data.length);
    } finally {
      calloc.free(p);
    }
  }

  static Uint8List? viewerPollOut(Pointer<Void> h) {
    if (!available || h == nullptr) {
      return null;
    }
    final buf = calloc<Uint8>(65536);
    final outLen = calloc<IntPtr>();
    try {
      final rc = _viewerPollOut(h, buf, 65536, outLen);
      if (rc != 0) {
        return null;
      }
      final n = outLen.value;
      if (n <= 0) {
        return null;
      }
      return Uint8List.fromList(buf.asTypedList(n));
    } finally {
      calloc.free(outLen);
      calloc.free(buf);
    }
  }

  static int viewerState(Pointer<Void> h) =>
      available && h != nullptr ? _viewerState(h) : -1;

  static int viewerWrite(Pointer<Void> h, int typ, Uint8List payload) {
    if (!available || h == nullptr) {
      return -1;
    }
    final p = calloc<Uint8>(payload.length);
    try {
      p.asTypedList(payload.length).setAll(0, payload);
      return _viewerWrite(h, typ, p, payload.length);
    } finally {
      calloc.free(p);
    }
  }

  static Uint8List? viewerRead(Pointer<Void> h) {
    if (!available || h == nullptr) {
      return null;
    }
    final buf = calloc<Uint8>(1 << 20);
    final outLen = calloc<IntPtr>();
    try {
      final rc = _viewerRead(h, buf, 1 << 20, outLen);
      if (rc != 0) {
        return null;
      }
      final n = outLen.value;
      if (n <= 0) {
        return null;
      }
      return Uint8List.fromList(buf.asTypedList(n));
    } finally {
      calloc.free(outLen);
      calloc.free(buf);
    }
  }

  static void viewerClose(Pointer<Void> h) {
    if (!available || h == nullptr) {
      return;
    }
    _viewerClose(h);
  }

  /// 0=已启动；-1=参数；-2=已在运行
  static int agentRun(
    String signalWs,
    String room,
    String root,
    int timeoutMs,
  ) {
    if (!agentRunAvailable) {
      return -1;
    }
    final s = signalWs.toNativeUtf8();
    final r = room.toNativeUtf8();
    final t = root.toNativeUtf8();
    try {
      return _agentRun!(s, r, t, timeoutMs);
    } finally {
      calloc.free(s);
      calloc.free(r);
      calloc.free(t);
    }
  }

  static int agentRunning() => _agentRunning?.call() ?? 0;

  static int grpcOpen(Pointer<Void> h) =>
      available && h != nullptr ? _grpcOpen(h) : -1;

  static void grpcClose() {
    if (available) {
      _grpcClose();
    }
  }

  static int grpcListFiles(
    Pointer<Void> h,
    Pointer<Utf8> root,
    Pointer<Uint8> out,
    int cap,
    Pointer<IntPtr> outLen,
  ) => _grpcList(h, root, out, cap, outLen);

  static int grpcListMedia(
    Pointer<Void> h,
    Pointer<Utf8> root,
    Pointer<Uint8> cursor,
    int cursorLen,
    int pageSize,
    Pointer<Uint8> out,
    int cap,
    Pointer<IntPtr> outLen,
  ) {
    final fn = _grpcListMedia;
    if (fn == null) {
      return -1;
    }
    return fn(h, root, cursor, cursorLen, pageSize, out, cap, outLen);
  }

  static int grpcThumbBatch(
    Pointer<Void> h,
    Pointer<Utf8> pathsNl,
    int maxEdge,
    Pointer<Uint8> out,
    int cap,
    Pointer<IntPtr> outLen,
  ) => _grpcThumb(h, pathsNl, maxEdge, out, cap, outLen);

  static int grpcReadFile(
    Pointer<Void> h,
    Pointer<Utf8> path,
    int maxBytes,
    Pointer<Uint8> out,
    int cap,
    Pointer<IntPtr> outLen,
  ) => _grpcRead(h, path, maxBytes, out, cap, outLen);

  static int grpcReadFileRange(
    Pointer<Void> h,
    Pointer<Utf8> path,
    int offset,
    int length,
    Pointer<Uint8> out,
    int cap,
    Pointer<IntPtr> outLen,
  ) {
    final fn = _grpcReadRange;
    if (fn == null) {
      return grpcReadFile(h, path, length > 0 ? length : cap, out, cap, outLen);
    }
    return fn(h, path, offset, length, out, cap, outLen);
  }

  static int grpcDeleteFile(
    Pointer<Void> h,
    Pointer<Utf8> path,
    Pointer<Uint8> out,
    int cap,
    Pointer<IntPtr> outLen,
  ) => _grpcDel(h, path, out, cap, outLen);
}
