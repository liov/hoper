# dart

## 简介
最初是个实验性的东西，搞得大而杂，有flutter，有原生，flutter里还有绑定lua和rust ffi，原生也有lua绑定，
全面但牺牲了稳定性，有很多用不到的也上了，没必要

## Getting Started
安装flutter
打开项目目录下执行：
pub get

### protobuf
dart pub global activate protoc_plugin
export PATH="$PATH:$HOME/.pub-cache/bin" (win:%USERPROFILE%\AppData\Local\Pub\Cache\bin)
export PATH="$PATH:$flutterSDK/bin/cache/dart-sdk/bin"
dart run generate.dart


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

## platforms
### windows

flutter create --platforms=windows .
### web
flutter create --platforms=web .


### 打包

#### android
##### key
keytool -genkey -v -keystore D:/key.jks -storetype JKS -keyalg RSA -keysize 2048 -validity 10000 -alias key
flutter build apk --release --target-platform android-arm64