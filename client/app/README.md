# dart

## 简介
最初是个实验性的东西，搞得大而杂，有flutter，有原生，flutter里还有绑定lua和rust ffi，原生也有lua绑定，
全面但牺牲了稳定性，有很多用不到的也上了，没必要

## Getting Started
安装flutter
打开项目目录下执行：
flutter packages upgrade
flutter pub outdated
flutter pub upgrade --major-versions
flutter pub get --no-example

### protobuf
flutter pub global activate protoc_plugin
export PATH="$PATH:$HOME/.pub-cache/bin" (win:%USERPROFILE%\AppData\Local\Pub\Cache\bin)
export PATH="$PATH:$flutterSDK/bin/cache/dart-sdk/bin"
protogen dart -p ../../thirdparty/protobuf/_proto -i ../../proto -o lib/gen/pb


### 静态库（rb / 远程相册 ICE）

在 `server/rust/rb` 下：

```bash
bash scripts/build_flutter_lib.sh          # staticLibs/macos|linux|windows/<arch>/
bash scripts/build_flutter_lib.sh --mobile # staticLibs/android/<abi>/ + ios/device/（IOS_BUILD_SIM=1 另编 sim/）
bash scripts/build_flutter_lib.sh --all    # 本机 + mobile
```

- Android：Rust `cdylib` → `dynLibs/android/arm64-v8a/librb.so`，Gradle `jniLibs` 打包；仅 **arm64-v8a**。
- Dart：`DynamicLibrary.executable()`（iOS/桌面）或 `open('librb.so')`（Android）。
- `build_flutter_lib.sh` 默认 `RB_LIB_FEATURES=viewer-ffi`（Viewer：ICE + HTTP/2 客户端，无 image 编码）；缩略图由 PC Agent 的 `rb`（`cargo build --features host`）编码。

bash scripts/build_flutter_lib.sh --android   # arm64-v8a → dynLibs/android/arm64-v8a/librb.so（Rust cdylib）
# RB_ANDROID_ABIS=x86 bash scripts/build_flutter_lib.sh --android  # 额外编 x86_64 模拟器

### 图标和开屏
dart run flutter_native_splash:create  --path=flutter_native_splash.yaml// 天坑，不看源码还不知道，flutter_native_splash是根据build.gradle判断编译SKD版本的，判断方法简单粗暴截取转整型，后面有注释识别不了
dart run flutter_launcher_icons:generate
### json
dart run build_runner build --delete-conflicting-outputs

## platforms
### windows

flutter create --platforms=windows .
### web
flutter create --platforms=web .

### macos
flutter create --platforms=macos .

### linux
flutter create --platforms=linux .

### ios
flutter create --platforms=ios .

### android
flutter create --platforms=android .

### 打包

#### android
##### key
keytool -genkey -v -keystore D:/key.jks -storetype JKS -keyalg RSA -keysize 2048 -validity 10000 -alias key
flutter build apk --release --target-platform android-arm64

##### 一键脚本（推荐）
```bash
cd client/app
bash scripts/build.sh --android --rb    # 编 librb.so + debug APK（团队 debug 签名）
bash scripts/build.sh --android --release --rb
bash scripts/build.sh --windows         # Windows 桌面 debug
bash scripts/build.sh --pub-only        # 仅依赖
```
PowerShell：`.\scripts\build.ps1 -Android -Rb`