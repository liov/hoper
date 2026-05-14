import 'dart:typed_data';

import 'package:app/gen/pb/remotebrowse/browse.pb.dart';
import 'package:app/gen/pb/remotebrowse/browse.service.pb.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/wire_codec.dart';
import 'package:app/remotebrowse/wire_transport.dart';

class RbWireSession {
  RbWireSession._(this._transport);

  final RbWireTransport _transport;

  factory RbWireSession(RbWireTransport transport) => RbWireSession._(transport);

  Future<void> close() => _transport.close();

  Future<List<RbFileEntry>> listFiles(String root) async {
    final req = ListFilesRequest(rootPath: root);
    await _transport.writeFrame(rbTypeFileIndex, req.writeToBuffer());
    final frame = await _transport.readFrame();
    if (frame.$1 != rbTypeFileIndex) {
      throw StateError('unexpected wire type ${frame.$1}');
    }
    final resp = ListFilesResponse.fromBuffer(frame.$2);
    return resp.entries.map((e) => RbFileEntry(id: e.id, name: e.name, size: e.size.toInt(), thumbHash: e.thumbHash)).toList();
  }

  Future<Uint8List> fetchThumb(String path, {int maxEdge = 256}) async {
    final req = ThumbnailRequest(path: path, maxEdge: maxEdge);
    await _transport.writeFrame(rbTypeThumbReq, req.writeToBuffer());
    final frame = await _transport.readFrame();
    if (frame.$1 != rbTypeThumbData) {
      throw StateError('unexpected wire type ${frame.$1}');
    }
    return Uint8List.fromList(ThumbnailResponse.fromBuffer(frame.$2).data);
  }
}
