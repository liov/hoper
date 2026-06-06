import 'dart:async';
import 'dart:ffi';
import 'package:app/gen/pb/remotebrowse/signal.pb.dart';
import 'package:app/remotebrowse/connect_cancel.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_ice_ffi_io.dart';
import 'package:app/remotebrowse/rb_ice_grpc.dart';
import 'package:app/remotebrowse/signal_session.dart';

bool get rbIceFfiAvailable => RbIceFfi.available;

/// Viewer 侧 ICE（Rust webrtc-ice，controlling）。
class RbIceViewerDialer {
  static const iceTimeout = Duration(seconds: 12);

  static Future<RbIceGrpcBridge?> tryConnectGrpc(RbSignalSession sig, {RbConnectCancel? cancel}) =>
      _connectGrpc(sig, cancel: cancel);
}

@Deprecated('use RbIceViewerDialer')
typedef RbIceDialer = RbIceViewerDialer;

Future<RbIceGrpcBridge?> _connectGrpc(RbSignalSession sig, {RbConnectCancel? cancel}) async {
  if (!RbIceFfi.available) {
    rbLog.fine('ice ffi unavailable');
    return null;
  }
  rbLog.fine('ice grpc start timeout=${RbIceViewerDialer.iceTimeout.inMilliseconds}ms');
  final h = RbIceFfi.viewerNew(RbIceViewerDialer.iceTimeout.inMilliseconds);
  if (h == nullptr) {
    return null;
  }
  final sub = sig.bindIceInbound((data) => RbIceFfi.viewerPush(h, data));
  final deadline = DateTime.now().add(RbIceViewerDialer.iceTimeout);
  try {
    while (DateTime.now().isBefore(deadline)) {
      cancel?.check();
      while (true) {
        final out = RbIceFfi.viewerPollOut(h);
        if (out == null) {
          break;
        }
        if (sig.isClosed || cancel?.cancelled == true) {
          break;
        }
        try {
          await sig.send(SignalEnvelope.fromBuffer(out));
        } catch (_) {
          break;
        }
      }
      final st = RbIceFfi.viewerState(h);
      if (st == 1) {
        rbLog.info('ice ready, open grpc');
        return RbIceGrpcBridge.tryOpen(h);
      }
      if (st < 0) {
        rbLog.fine('ice failed state=$st');
        RbIceFfi.viewerClose(h);
        return null;
      }
      await Future<void>.delayed(const Duration(milliseconds: 20));
    }
    rbLog.fine('ice timeout');
    RbIceFfi.viewerClose(h);
    return null;
  } finally {
    await sub.cancel();
  }
}
