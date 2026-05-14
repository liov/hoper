import 'dart:convert';
import 'dart:typed_data';

const rbWireVer = 1;
const rbTypeFileIndex = 2;
const rbTypeThumbReq = 3;
const rbTypeThumbData = 4;
const rbRoleViewer = 0;

Uint8List rbRelayJoinBytes(String sessionId, int role) {
  final id = _uuidBytes(sessionId);
  final buf = BytesBuilder();
  buf.add(utf8.encode('RBRL'));
  buf.add([1, ...id, role]);
  return buf.toBytes();
}

List<int> _uuidBytes(String sessionId) {
  final hex = sessionId.replaceAll('-', '');
  return List.generate(16, (i) => int.parse(hex.substring(i * 2, i * 2 + 2), radix: 16));
}

Uint8List rbEncodeWireFrame(int typ, Uint8List payload) {
  final hdr = ByteData(6);
  hdr.setUint8(0, rbWireVer);
  hdr.setUint8(1, typ);
  hdr.setUint32(2, payload.length, Endian.big);
  final frame = BytesBuilder();
  frame.add(hdr.buffer.asUint8List());
  frame.add(payload);
  return frame.toBytes();
}

(int, Uint8List) rbDecodeWireFrame(Uint8List raw) {
  final hdr = ByteData.sublistView(raw, 0, 6);
  if (hdr.getUint8(0) != rbWireVer) {
    throw StateError('bad wire version');
  }
  final n = hdr.getUint32(2, Endian.big);
  return (hdr.getUint8(1), Uint8List.sublistView(raw, 6, 6 + n));
}

Uint8List rbEncodeRelayFrame(Uint8List payload) {
  final sz = ByteData(4)..setUint32(0, payload.length, Endian.big);
  return Uint8List.fromList([...sz.buffer.asUint8List(), ...payload]);
}

int rbDecodeRelaySize(Uint8List raw) => ByteData.sublistView(raw).getUint32(0, Endian.big);
