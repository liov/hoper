#!/usr/bin/env bash
# client/app 一键编译（依赖、可选 proto/rb 静态库、Flutter 打包）
#
#   bash scripts/build.sh                    # 本机：pub get + flutter build windows（debug）
#   bash scripts/build.sh --android          # Android arm64 APK（debug，团队 debug 签名）
#   bash scripts/build.sh --android --release
#   bash scripts/build.sh --android --rb       # 先编 dynLibs/android/.../librb.so 再 APK
#   bash scripts/build.sh --android --no-proto --no-codegen  # 仅改 Dart/UI 时加快
#   bash scripts/build.sh --proto            # protogen dart
#   bash scripts/build.sh --codegen          # build_runner（riverpod .g.dart）
#   bash scripts/build.sh --pub-only         # 仅 flutter pub get
set -euo pipefail

APP="$(cd "$(dirname "$0")/.." && pwd)"
HOPER="$(cd "$APP/../.." && pwd)"
RB="$HOPER/server/rust/rb"
PROTO_IN="$HOPER/proto"
PROTO_PLUGIN="$HOPER/thirdparty/protobuf/_proto"

DO_PROTO=1
DO_CODEGEN=1
DO_RB=0
DO_PUB_ONLY=0
FLUTTER_BUILD=(build windows --debug)
ANDROID_RELEASE=0

usage() {
  sed -n '2,12p' "$0" | sed 's/^# \{0,1\}//'
  exit "${1:-0}"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help) usage 0 ;;
    --proto) DO_PROTO=1 ;;
    --no-proto) DO_PROTO=0 ;;
    --codegen) DO_CODEGEN=1 ;;
    --no-codegen) DO_CODEGEN=0 ;;
    --rb) DO_RB=1 ;;
    --pub-only) DO_PUB_ONLY=1 ;;
    --android)
      FLUTTER_BUILD=(build apk --debug --target-platform android-arm64)
      ;;
    --release)
      ANDROID_RELEASE=1
      FLUTTER_BUILD=(build apk --release --target-platform android-arm64)
      ;;
    --windows) FLUTTER_BUILD=(build windows --debug) ;;
    --apk) FLUTTER_BUILD=(build apk --debug --target-platform android-arm64) ;;
    *) echo "unknown arg: $1" >&2; usage 1 ;;
  esac
  shift
done

flutter_sdk() {
  local props="$APP/android/local.properties" line
  if [[ -f "$props" ]]; then
    line="$(grep -E '^flutter\.sdk=' "$props" 2>/dev/null | tail -1 || true)"
    if [[ -n "$line" ]]; then
      local sdk="${line#flutter.sdk=}"
      sdk="${sdk//\\//}"
      if [[ -d "$sdk" ]]; then
        echo "$sdk"
        return 0
      fi
    fi
  fi
  if [[ -n "${FLUTTER_ROOT:-}" && -d "$FLUTTER_ROOT" ]]; then
    echo "$FLUTTER_ROOT"
    return 0
  fi
  command -v flutter >/dev/null 2>&1 && command -v flutter | xargs dirname | xargs dirname && return 0
  return 1
}

FLUTTER_BIN="$(flutter_sdk || true)"
if [[ -z "$FLUTTER_BIN" ]]; then
  echo "[app-build] ERROR: 未找到 Flutter SDK。请设置 android/local.properties 的 flutter.sdk 或 FLUTTER_ROOT" >&2
  exit 1
fi
export PATH="$FLUTTER_BIN/bin:$FLUTTER_BIN/bin/cache/dart-sdk/bin:${PATH:-}"
echo "[app-build] FLUTTER=$FLUTTER_BIN"
cd "$APP"

if [[ "$DO_PROTO" -eq 1 ]]; then
  echo "[app-build] protogen dart …"
  if ! command -v protogen >/dev/null 2>&1; then
    echo "[app-build] WARN: protogen 不在 PATH，跳过（需: protogen dart -p … -i … -o lib/gen/pb）" >&2
  else
    protogen dart -p "$PROTO_PLUGIN" -i "$PROTO_IN" -o lib/gen/pb
  fi
fi

if [[ "$DO_RB" -eq 1 ]]; then
  echo "[app-build] build_flutter_lib …"
  if [[ "${FLUTTER_BUILD[1]:-}" == apk ]]; then
    bash "$RB/scripts/build_flutter_lib.sh" --android
  else
    bash "$RB/scripts/build_flutter_lib.sh"
  fi
fi

echo "[app-build] flutter pub get …"
flutter pub get --no-example

if [[ "$DO_CODEGEN" -eq 1 ]]; then
  echo "[app-build] build_runner …"
  dart run build_runner build --delete-conflicting-outputs
fi

if [[ "$DO_PUB_ONLY" -eq 1 ]]; then
  echo "[app-build] done (pub only)"
  exit 0
fi

# Android：需 android/debug-key.properties + team-debug.keystore（见 debug-key.properties.example）
if [[ "${FLUTTER_BUILD[1]:-}" == apk ]]; then
  if [[ ! -f "$APP/android/debug-key.properties" ]]; then
    echo "[app-build] WARN: 无 android/debug-key.properties，将使用本机默认 debug.keystore" >&2
  fi
  if [[ ! -f "$APP/dynLibs/android/arm64-v8a/librb.so" ]]; then
    echo "[app-build] WARN: 无 dynLibs/android/arm64-v8a/librb.so，远程相册 ICE 可能不可用。加 --rb 先编 Rust" >&2
  fi
fi

echo "[app-build] flutter ${FLUTTER_BUILD[*]} …"
flutter "${FLUTTER_BUILD[@]}"

echo "[app-build] 产物:"
case "${FLUTTER_BUILD[1]:-}" in
  apk)
    if [[ "$ANDROID_RELEASE" -eq 1 || "${FLUTTER_BUILD[*]}" == *release* ]]; then
      ls -la "$APP/build/app/outputs/flutter-apk/app-release.apk" 2>/dev/null || true
    else
      ls -la "$APP/build/app/outputs/flutter-apk/app-debug.apk" 2>/dev/null || true
    fi
    ;;
  windows)
    ls -la "$APP/build/windows/"*/runner/Debug/*.exe 2>/dev/null || true
    ;;
esac
