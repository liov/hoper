import 'dart:async';

import 'package:app/remotebrowse/signal_session.dart';

/// Viewer 建链过程可取消；[cancel] 会关闭信令 WS。
class RbConnectCancel {
  RbSignalSession? sig;
  bool cancelled = false;

  void bind(RbSignalSession s) => sig = s;

  /// 建链成功关闭信令：仅解绑，不置 [cancelled]（避免后续 [check] 误判为已取消）。
  void detachSignal() => sig = null;

  void cancel() {
    unawaited(closeSignal());
  }

  /// 关闭信令 WS；可 await，避免与下一次 connect 的 close/send 竞态。
  Future<void> closeSignal() async {
    if (cancelled) {
      return;
    }
    cancelled = true;
    final s = sig;
    sig = null;
    if (s != null) {
      await s.close();
    }
  }

  void check() {
    if (cancelled) {
      throw const RbConnectCancelled();
    }
  }
}

class RbConnectCancelled implements Exception {
  const RbConnectCancelled();
}
