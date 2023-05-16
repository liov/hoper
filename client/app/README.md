# dart

## 简介
最初是个实验性的东西，搞得大而杂，有flutter，有原生，flutter里还有绑定lua和rust ffi，原生也有lua绑定，
全面但牺牲了稳定性，有很多用不到的也上了，没必要

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


## protobuf
dart pub global activate protoc_plugin
export PATH="$PATH:$HOME/.pub-cache/bin" //$HOME\AppData\Local\Pub\Cache\bin
export PATH="$PATH:$flutterSDK/bin/cache\dart-sdk\bin"
dart run generate.dart

怪不得[flutter_lua](https://github.com/drydart/flutter_lua)插件有libgojni.so 且只支持lua5.2,
原来用的是[go实现的lua虚拟机](https://github.com/Shopify/go-lua),比clua慢6倍
这有个[支持lua5.4的的clua](https://github.com/tgarm/flutter-luavm)

### pigeon
flutter pub run pigeon --input lib/pigeons/route.dart \
  --dart_out lib/pigeons/pigeon.dart \
  --objc_header_out ios/Runner/route.h \
  --objc_source_out ios/Runner/route.m \
  --java_out android/app/src/main/java/io/flutter/plugins/Route.java \
  --java_package "io.flutter.plugins"
### 动态库

#### android
##### 创建
rustup target add x86_64-linux-android armv7-linux-androideabi aarch64-linux-android
##### 编译
cargo build --release --target=armv7-linux-androideabi

### 图标和开屏
flutter pub run flutter_native_splash:create // 天坑，不看源码还不知道，flutter_native_splash是根据build.gradle判断编译SKD版本的，判断方法简单粗暴截取转整型，后面有注释识别不了
flutter pub run flutter_launcher_icons:main
### json
flutter pub run build_runner build --delete-conflicting-outputs

Get太灵活了，写法太多反而不好

flutter create --platforms=windows .
flutter create --platforms=web .

### 打包

#### android
##### key
keytool -genkey -v -keystore D:/key.jks -keyalg RSA -keysize 2048 -validity 10000 -alias key
flutter build apk --release --target-platform android-arm64