import 'dart:ffi';
import 'package:ffi/ffi.dart';
import 'dart:io'; // For Platform.isX

DynamicLibrary findDynamicLibrary(String name, String dir) {
  if (!dir.endsWith('/')) dir = dir + '/';
  if (Platform.isLinux) return DynamicLibrary.open('${dir}lib$name.so');
  if (Platform.isMacOS) return DynamicLibrary.open('${dir}lib$name.dylib');
  if (Platform.isWindows) return DynamicLibrary.open('$dir$name.dll');
  return DynamicLibrary.process();
}


final DynamicLibrary nativeRustLib = findDynamicLibrary("rust_lib",r'D:/code/hoper/client/flutter/lib-dyn/rust-lib/target/release/');

final void Function(int x) funServe =
nativeRustLib
    .lookup<NativeFunction<Void Function(Int32)>>("server")
    .asFunction();


main(){
  funServe(3000);
}