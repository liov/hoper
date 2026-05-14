import 'dart:typed_data';

abstract class RbWireTransport {
  Future<void> writeFrame(int typ, Uint8List payload);
  Future<(int, Uint8List)> readFrame();
  Future<void> close();
}
