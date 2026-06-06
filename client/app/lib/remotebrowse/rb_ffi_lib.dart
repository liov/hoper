import 'dart:ffi';
import 'dart:io';

/// Android：Rust 直编 `dynLibs/android/<abi>/librb.so`；桌面/iOS 静态链入可执行文件。
DynamicLibrary findRbStaticLibrary() {
  if (Platform.isIOS ||
      Platform.isMacOS ||
      Platform.isLinux ||
      Platform.isWindows) {
    return DynamicLibrary.executable();
  }
  if (Platform.isAndroid) {
    try {
      return DynamicLibrary.open('librb.so');
      // ignore: avoid_catching_errors
    } on ArgumentError {
      return DynamicLibrary.executable();
    }
  }
  return DynamicLibrary.process();
}
