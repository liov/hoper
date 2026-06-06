import 'dart:convert';
import 'dart:typed_data';

const rbRoleViewer = 0;
const rbRoleAgent = 1;

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
