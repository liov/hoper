import 'dart:ffi';
import 'dart:io';
import 'dart:math';

DynamicLibrary findDynamicLibrary(String name, [String? dir]) {
  if (Platform.isAndroid) {
    try {
      return DynamicLibrary.open('lib$name.so');
      // ignore: avoid_catching_errors
    } on ArgumentError {
      final appIdAsBytes = File('/proc/self/cmdline').readAsBytesSync();
      final endOfAppId = max(appIdAsBytes.indexOf(0), 0);
      final appId = String.fromCharCodes(appIdAsBytes.sublist(0, endOfAppId));
      return DynamicLibrary.open('/data/data/$appId/lib/lib$name.so');
    }
  }
  if (Platform.isIOS) {
    for (final path in ['Frameworks/lib$name.dylib', 'lib$name.dylib']) {
      try {
        return DynamicLibrary.open(path);
        // ignore: avoid_catching_errors
      } on ArgumentError {
        continue;
      }
    }
    return DynamicLibrary.process();
  }
  final base = dir ?? _desktopDynLibsDir();
  final root = base.endsWith('/') ? base : '$base/';
  if (Platform.isLinux) return DynamicLibrary.open('${root}lib$name.so');
  if (Platform.isMacOS) return DynamicLibrary.open('${root}lib$name.dylib');
  if (Platform.isWindows) return DynamicLibrary.open('$root$name.dll');
  return DynamicLibrary.process();
}

String _desktopDynLibsDir() {
  final arch = _hostArchSlug();
  if (Platform.isMacOS) return 'dynLibs/macos/$arch';
  if (Platform.isLinux) return 'dynLibs/linux/$arch';
  if (Platform.isWindows) return 'dynLibs/windows/$arch';
  return 'dynLibs';
}

String _hostArchSlug() {
  final n = Abi.current().toString();
  if (n.contains('Arm64')) return 'arm64';
  if (n.contains('Riscv64')) return 'riscv64';
  if (n.contains('Riscv32')) return 'riscv32';
  if (n.contains('Arm')) return 'arm32';
  if (Platform.isWindows) {
    if (n.contains('IA32')) return 'x86';
    return 'x64';
  }
  if (n.contains('X64')) return 'x86_64';
  if (n.contains('IA32')) return 'x86';
  return 'unknown';
}
