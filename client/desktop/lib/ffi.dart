import 'dart:ffi';
import 'dart:io';
import 'dart:math'; // For Platform.isX

DynamicLibrary findDynamicLibrary(String name, String dir) {
  if (Platform.isAndroid) {
    try {
      return DynamicLibrary.open('lib$name.so');
      // ignore: avoid_catching_errors
    } on ArgumentError {
      final appIdAsBytes = File('/proc/self/cmdline').readAsBytesSync();

      // app id ends with the first \0 character in here.
      final endOfAppId = max(appIdAsBytes.indexOf(0), 0);
      final appId = String.fromCharCodes(appIdAsBytes.sublist(0, endOfAppId));

      return DynamicLibrary.open('/data/data/$appId/lib/lib$name.so');
    }
  }
  if (!dir.endsWith('/')) dir = dir + '/';
  if (Platform.isLinux) return DynamicLibrary.open('${dir}lib$name.so');
  if (Platform.isMacOS) return DynamicLibrary.open('${dir}lib$name.dylib');
  if (Platform.isWindows) return DynamicLibrary.open('$dir$name.dll');
  return DynamicLibrary.process();
}


final DynamicLibrary nativeRustLib = findDynamicLibrary("rust",r'libraries');

final void Function(int x) funServe =
nativeRustLib
    .lookup<NativeFunction<Void Function(Int32)>>("server")
    .asFunction();

void server(int x){
  funServe(x);
}