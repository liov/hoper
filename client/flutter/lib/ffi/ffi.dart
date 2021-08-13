import 'dart:ffi';
import 'package:ffi/ffi.dart';
import 'dart:io'; // For Platform.isX

DynamicLibrary findDynamicLibrary(String name, String dir) {
  if (!dir.endsWith('/')) dir = dir + '/';
  if (Platform.isAndroid) return DynamicLibrary.open('lib$name.so');
  if (Platform.isLinux) return DynamicLibrary.open('${dir}lib$name.so');
  if (Platform.isMacOS) return DynamicLibrary.open('${dir}lib$name.dylib');
  if (Platform.isWindows) return DynamicLibrary.open('$dir$name.dll');
  return DynamicLibrary.process();
}

final DynamicLibrary nativeAddLib = findDynamicLibrary("rust",'libraries');


final Pointer<Utf8> Function(Pointer<Utf8> x) nativeGreeting =
nativeAddLib
    .lookup<NativeFunction<Pointer<Utf8> Function(Pointer<Utf8>)>>("rust_greeting")
    .asFunction();

String greeting(){
  final String myString = "😎👿💬";
  final Pointer<Utf8> charPointer = nativeGreeting(myString.toNativeUtf8());
  return charPointer.toDartString();
}