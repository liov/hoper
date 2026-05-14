import 'dart:ffi';
import 'dart:typed_data';

import 'package:applib/util/ffi.dart';
import 'package:ffi/ffi.dart';

typedef _ViewerNewNative = Pointer<Void> Function(Int32 timeoutMs);
typedef _ViewerNewDart = Pointer<Void> Function(int timeoutMs);
typedef _ViewerPushNative = Int32 Function(Pointer<Void> h, Pointer<Uint8> data, IntPtr len);
typedef _ViewerPushDart = int Function(Pointer<Void> h, Pointer<Uint8> data, int len);
typedef _ViewerPollOutNative = Int32 Function(Pointer<Void> h, Pointer<Uint8> buf, IntPtr cap, Pointer<IntPtr> outLen);
typedef _ViewerPollOutDart = int Function(Pointer<Void> h, Pointer<Uint8> buf, int cap, Pointer<IntPtr> outLen);
typedef _ViewerStateNative = Int32 Function(Pointer<Void> h);
typedef _ViewerStateDart = int Function(Pointer<Void> h);
typedef _ViewerWriteNative = Int32 Function(Pointer<Void> h, Uint8 typ, Pointer<Uint8> data, IntPtr len);
typedef _ViewerWriteDart = int Function(Pointer<Void> h, int typ, Pointer<Uint8> data, int len);
typedef _ViewerReadNative = Int32 Function(Pointer<Void> h, Pointer<Uint8> buf, IntPtr cap, Pointer<IntPtr> outLen);
typedef _ViewerReadDart = int Function(Pointer<Void> h, Pointer<Uint8> buf, int cap, Pointer<IntPtr> outLen);
typedef _ViewerCloseNative = Void Function(Pointer<Void> h);
typedef _ViewerCloseDart = void Function(Pointer<Void> h);

class RbIceFfi {
  static DynamicLibrary? _lib;
  static late final _ViewerNewDart _viewerNew;
  static late final _ViewerPushDart _viewerPush;
  static late final _ViewerPollOutDart _viewerPollOut;
  static late final _ViewerStateDart _viewerState;
  static late final _ViewerWriteDart _viewerWrite;
  static late final _ViewerReadDart _viewerRead;
  static late final _ViewerCloseDart _viewerClose;
  static var _inited = false;
  static var _ready = false;

  static bool get available {
    _ensure();
    return _ready;
  }

  static void _ensure() {
    if (_inited) {
      return;
    }
    _inited = true;
    try {
      _lib = findDynamicLibrary('rbclient', 'libraries');
      _viewerNew = _lib!.lookupFunction<_ViewerNewNative, _ViewerNewDart>('rb_ice_viewer_new');
      _viewerPush = _lib!.lookupFunction<_ViewerPushNative, _ViewerPushDart>('rb_ice_viewer_push');
      _viewerPollOut = _lib!.lookupFunction<_ViewerPollOutNative, _ViewerPollOutDart>('rb_ice_viewer_poll_out');
      _viewerState = _lib!.lookupFunction<_ViewerStateNative, _ViewerStateDart>('rb_ice_viewer_state');
      _viewerWrite = _lib!.lookupFunction<_ViewerWriteNative, _ViewerWriteDart>('rb_ice_viewer_write');
      _viewerRead = _lib!.lookupFunction<_ViewerReadNative, _ViewerReadDart>('rb_ice_viewer_read');
      _viewerClose = _lib!.lookupFunction<_ViewerCloseNative, _ViewerCloseDart>('rb_ice_viewer_close');
      _ready = true;
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

  static int viewerState(Pointer<Void> h) => available && h != nullptr ? _viewerState(h) : -1;

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
}
