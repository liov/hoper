import 'package:app/remotebrowse/viewer_session.dart';
import 'package:app/remotebrowse/wire_session.dart';

typedef RbRelaySession = RbWireSession;

Future<RbRelaySession> openRbViewer(Uri signalWs, String room) => RbViewerSession.connect(signalWs, room);
