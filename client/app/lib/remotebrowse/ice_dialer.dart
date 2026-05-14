import 'dart:async';
import 'dart:ffi';
import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:app/remotebrowse/rb_ice_ffi.dart';
import 'package:app/remotebrowse/signal_session.dart';
import 'package:app/remotebrowse/wire_codec.dart';
import 'package:app/remotebrowse/wire_transport.dart';

class RbIceDialer {
  static const iceTimeout = Duration(seconds: 12);

  static Future<RbWireTransport?> tryConnect(RbSignalSession sig) async {
    if (!RbIceFfi.available) {
      return null;
    }
    final h = RbIceFfi.viewerNew(iceTimeout.inMilliseconds);
    if (h == nullptr) {
      return null;
    }
    final sub = sig.bindIceInbound((data) => RbIceFfi.viewerPush(h, data));
    final deadline = DateTime.now().add(iceTimeout);
    try {
      while (DateTime.now().isBefore(deadline)) {
        while (true) {
          final out = RbIceFfi.viewerPollOut(h);
          if (out == null) {
            break;
          }
          await sig.send(SignalEnvelope.fromBuffer(out));
        }
        final st = RbIceFfi.viewerState(h);
        if (st == 1) {
          return RbIceNativeTransport(h);
        }
        if (st < 0) {
          RbIceFfi.viewerClose(h);
          return null;
        }
        await Future<void>.delayed(const Duration(milliseconds: 20));
      }
      RbIceFfi.viewerClose(h);
      return null;
    } finally {
      await sub.cancel();
    }
  }
}

class RbIceNativeTransport implements RbWireTransport {
  RbIceNativeTransport(this._h);

  final Pointer<Void> _h;
  var _closed = false;

  @override
  Future<void> writeFrame(int typ, Uint8List payload) async {
    if (_closed) {
      throw StateError('ice closed');
    }
    if (RbIceFfi.viewerWrite(_h, typ, payload) != 0) {
      throw StateError('ice write failed');
    }
  }

  @override
  Future<(int, Uint8List)> readFrame() async {
    if (_closed) {
      throw StateError('ice closed');
    }
    final buf = RbIceFfi.viewerRead(_h);
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
    RbIceFfi.viewerClose(_h);
  }
}
