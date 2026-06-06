import 'package:app/remotebrowse/connect_cancel.dart';
import 'package:app/remotebrowse/rb_ice_grpc_stub.dart';
import 'package:app/remotebrowse/signal_session.dart';

bool get rbIceFfiAvailable => false;

class RbIceViewerDialer {
  static Future<RbIceGrpcBridge?> tryConnectGrpc(RbSignalSession sig, {RbConnectCancel? cancel}) async => null;
}

@Deprecated('use RbIceViewerDialer')
typedef RbIceDialer = RbIceViewerDialer;
