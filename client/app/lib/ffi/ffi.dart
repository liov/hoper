import 'dart:ffi';
import 'dart:isolate';
import 'package:ffi/ffi.dart';
import 'package:applib/util/ffi.dart';
import 'dart:io'; // For Platform.isX
import 'dart:math'; // For Platform.isX

final DynamicLibrary nativeGreetingLib = findDynamicLibrary("greeting",'libraries');
final DynamicLibrary nativeHelloLib = findDynamicLibrary("hello",'libraries');
final DynamicLibrary nativeRustLib = findDynamicLibrary("rust",'libraries');


final Pointer<Utf8> Function(Pointer<Utf8> x) rustGreeting =
nativeGreetingLib
    .lookup<NativeFunction<Pointer<Utf8> Function(Pointer<Utf8>)>>("rust_greeting")
    .asFunction();

String callRustGreeting(){
  const String myString = "😎👿💬";
  final Pointer<Utf8> charPointer = rustGreeting(myString.toNativeUtf8());
  return charPointer.toDartString();
}

final Pointer<Utf8> Function(Pointer<Utf8> x) goPrint =
nativeHelloLib
    .lookup<NativeFunction<Pointer<Utf8> Function(Pointer<Utf8>)>>("goprint")
    .asFunction();

String callGoPrint(String s){
  final Pointer<Utf8> charPointer = goPrint(s.toNativeUtf8());
  return charPointer.toDartString();
}

final void Function(int x) funServe =
nativeRustLib
    .lookup<NativeFunction<Void Function(Int32)>>("server")
    .asFunction();

void serve(int port){
  Isolate.spawn(funServe, 3000);
}

final int Function(int x) test =
nativeRustLib
    .lookup<NativeFunction<Int32 Function(Int32)>>("test")
    .asFunction();
