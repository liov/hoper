import 'dart:ffi';
import 'package:ffi/ffi.dart';
import 'dart:io'; // For Platform.isX

String findDynamicLibraryFile(String name, String dir) {
  if (!dir.endsWith('/')) dir = dir + '/';
  if (Platform.isAndroid) return 'lib$name.so';
  if (Platform.isLinux) return '${dir}lib$name.so';
  if (Platform.isMacOS) return '${dir}lib$name.dylib';
  if (Platform.isWindows) return '$dir$name.dll';
  throw Exception("Platform not implemented");
}

final DynamicLibrary nativeAddLib = Platform.isAndroid
    ? DynamicLibrary.open(findDynamicLibraryFile("rust",'libraries'))
    : DynamicLibrary.process();

final Pointer<Utf8> Function(Pointer<Utf8> x) nativeGreeting =
nativeAddLib
    .lookup<NativeFunction<Pointer<Utf8> Function(Pointer<Utf8>)>>("rust_greeting")
    .asFunction();

String greeting(){
  final String myString = "😎👿💬";
  final Pointer<Utf8> charPointer = Utf8.toUtf8(myString);
  return Utf8.fromUtf8(nativeGreeting(charPointer));
}