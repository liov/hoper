# dart

## 简介
最初是个实验性的东西，搞得大而杂，有flutter，有原生，flutter里还有绑定lua和rust ffi，原生也有lua绑定，
全面但牺牲了稳定性，有很多用不到的也上了，没必要

A new Flutter project.

## Getting Started
安装flutter
打开项目目录下执行：
pub get

### dev step(2020年开始跳过这一步)

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


### protobuf
dart pub global activate protoc_plugin
export PATH="$PATH:$HOME/.pub-cache/bin" (win:%USERPROFILE%\AppData\Local\Pub\Cache\bin)
export PATH="$PATH:$flutterSDK/bin/cache/dart-sdk/bin"
dart run generate.dart
protoc --dart_out=grpc:lib/src/generated -I../../../protobuf ../../../protobuf/helloworld.proto

### lua(废弃)
怪不得[flutter_lua](https://github.com/drydart/flutter_lua)插件有libgojni.so 且只支持lua5.2,
原来用的是[go实现的lua虚拟机](https://github.com/Shopify/go-lua),比clua慢6倍
这有个[支持lua5.4的的clua](https://github.com/tgarm/flutter-luavm)

### pigeon
dart run pigeon --input lib/pigeons/route.dart \
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
dart run flutter_native_splash:create // 天坑，不看源码还不知道，flutter_native_splash是根据build.gradle判断编译SKD版本的，判断方法简单粗暴截取转整型，后面有注释识别不了
dart run flutter_launcher_icons:main
### json
dart run build_runner build --delete-conflicting-outputs

Get太灵活了，写法太多反而不好

## platforms
### windows
https://dist.nuget.org/win-x86-commandline/latest/nuget.exe
flutter create --platforms=windows .
Should work as is in debug mode (sqlite3.dll is bundled).

In release mode, add sqlite3.dll in same folder as your executable.

sqfliteFfiInit is provided as an implementation reference for loading the sqlite library. Please look at sqlite3 if you want to override the behavior.
# Linux
libsqlite3 and libsqlite3-dev linux packages are required.

One time setup for Ubuntu (to run as root):

dart tool/linux_setup.dart
or

sudo apt-get -y install libsqlite3-0 libsqlite3-dev
# MacOS
Should work as is.

### web
flutter create --platforms=web .
Look at package sqflite_common_ffi_web for experimental Web support.

### 打包

#### android
##### key
keytool -genkey -v -keystore D:/key.jks -keyalg RSA -keysize 2048 -validity 10000 -alias key
flutter build apk --release --target-platform android-arm64