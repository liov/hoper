# dart

A new Flutter project.

## Getting Started

This project is a starting point for a Flutter application.

A few resources to get you started if this is your first Flutter project:

- [Lab: Write your first Flutter app](https://flutter.dev/docs/get-started/codelab)
- [Cookbook: Useful Flutter samples](https://flutter.dev/docs/cookbook)

For help getting started with Flutter, view our
[online documentation](https://flutter.dev/docs), which offers tutorials,
samples, guidance on mobile development, and a full 

# dev step
choco install dart-sdk
pub global activate protoc_plugin
protoc --dart_out=grpc:lib/src/generated -I../../../../proto ../../../../proto/helloworld.proto
pub get
dart xxx

export PUB_HOSTED_URL=https://pub.flutter-io.cn
export FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn

maven { url 'http://download.flutter.io' }

-- 新版本无需替换
```groovy
String storageUrl = System.getenv('FLUTTER_STORAGE_BASE_URL') ?: "https://storage.googleapis.com"

repositories {
    google()
    jcenter()
    maven {
        url "$storageUrl/download.flutter.io"
    }
}
```
修改flutter安装目录下三个文件：flutter/packages/flutter_tools/gradle/resolve_dependencies.gradle

flutter/packages/flutter_tools/gradle/aar_init_script.gradle

flutter/packages/flutter_tools/gradle/flutter.gradle

将其中的：https://storage.googleapis.com/download.flutter.io 替换为：http://download.flutter.io 后重新编译
https://storage.flutter-io.cn/download.flutter.io

keytool -genkey -v -keystore D:/key.jks -keyalg RSA -keysize 2048 -validity 10000 -alias key

pub global activate protoc_plugin

export PATH="$PATH:$HOME/.pub-cache/bin" //$HOME\AppData\Local\Pub\Cache\bin

protoc --dart_out=grpc:lib/generated --proto_path=../../std_proto  -Iprotos ../../std_proto/user/user.enum.proto

怪不得[flutter_lua](https://github.com/drydart/flutter_lua)插件有libgojni.so 且只支持lua5.2,
原来用的是[go实现的lua虚拟机](https://github.com/Shopify/go-lua),比clua慢6倍
这有个[支持lua5.4的的clua](https://github.com/tgarm/flutter-luavm)


flutter pub run pigeon --input pigeons/route.dart
rustup target add x86_64-linux-android armv7-linux-androideabi aarch64-linux-android
cargo build --release --target=armv7-linux-androideabi
flutter build apk --release --target-platform android-arm64

#todo
集成[MLN](https://github.com/momotech/MLN) Demo都跑不起来，没必要为了一个热更新在客户端上浪费太多时间
Rust FFI
Webview

放弃flutter_boost整合原生+MLN为app提供热更新的能力，MLN小众，坑看起来也不少，lua bind 原生看起来也有学习成本,动态部分用webview替代实现
RustFFI调试是个问题

不是专业的别折腾，且使用体验并不好,走主流方案
不过LN是个不错的东西