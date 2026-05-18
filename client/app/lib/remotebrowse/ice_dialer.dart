import 'dart:async';
import 'dart:ffi';
import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:app/remotebrowse/rb_ice_ffi.dart';
import 'package:app/remotebrowse/signal_session.dart';
import 'package:app/remotebrowse/wire_codec.dart';
import 'package:app/remotebrowse/wire_transport.dart';

/// Viewer 侧 ICE（Rust webrtc-ice，controlling）。
class RbIceViewerDialer {
  static const iceTimeout = Duration(seconds: 12);

  static Future<RbWireTransport?> tryConnect(RbSignalSession sig) => _connect(sig, RbIceFfi.viewerNew, RbIceFfi.viewerPush, RbIceFfi.viewerPollOut, RbIceFfi.viewerState, RbIceFfi.viewerClose);
}

@Deprecated('use RbIceViewerDialer')
typedef RbIceDialer = RbIceViewerDialer;

Future<RbWireTransport?> _connect(
  RbSignalSession sig,
  Pointer<Void> Function(int) newFn,
  void Function(Pointer<Void>, Uint8List) pushFn,
  Uint8List? Function(Pointer<Void>) pollFn,
  int Function(Pointer<Void>) stateFn,
  void Function(Pointer<Void>) closeFn,
) async {
  if (!RbIceFfi.available) {
    return null;
  }
  final h = newFn(RbIceViewerDialer.iceTimeout.inMilliseconds);
  if (h == nullptr) {
    return null;
  }
  final sub = sig.bindIceInbound((data) => pushFn(h, data));
  final deadline = DateTime.now().add(RbIceViewerDialer.iceTimeout);
  try {
    while (DateTime.now().isBefore(deadline)) {
      while (true) {
        final out = pollFn(h);
        if (out == null) {
          break;
        }
        await sig.send(SignalEnvelope.fromBuffer(out));
      }
      final st = stateFn(h);
      if (st == 1) {
        return RbIceNativeTransport(h, readFn: RbIceFfi.viewerRead, writeFn: RbIceFfi.viewerWrite, closeFn: closeFn);
      }
      if (st < 0) {
        closeFn(h);
        return null;
      }
      await Future<void>.delayed(const Duration(milliseconds: 20));
    }
    closeFn(h);
    return null;
  } finally {
    await sub.cancel();
  }
}

class RbIceNativeTransport implements RbWireTransport {
  RbIceNativeTransport(
    this._h, {
    required this.readFn,
    required this.writeFn,
    required this.closeFn,
  });

  final Pointer<Void> _h;
  final Uint8List? Function(Pointer<Void>) readFn;
  final int Function(Pointer<Void>, int, Uint8List) writeFn;
  final void Function(Pointer<Void>) closeFn;
  var _closed = false;

  @override
  Future<void> writeFrame(int typ, Uint8List payload) async {
    if (_closed) {
      throw StateError('ice closed');
    }
    if (writeFn(_h, typ, payload) != 0) {
      throw StateError('ice write failed');
    }
  }

  @override
  Future<(int, Uint8List)> readFrame() async {
    if (_closed) {
      throw StateError('ice closed');
    }
    final buf = readFn(_h);
    if (buf == null) {
      throw StateError('ice read failed');
    }
    final frame = rbDecodeWireFrame(buf);
    return (frame.$1, frame.$2);
  }

  @override
  Future<void> close() async {
    if (_closed) {
      return;
    }
    _closed = true;
    closeFn(_h);
  }
}
