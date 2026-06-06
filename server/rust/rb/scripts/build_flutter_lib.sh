#!/usr/bin/env bash
# Flutter librb 静态库 → client/app/staticLibs/
#   bash scripts/build_flutter_lib.sh              # 本机 → staticLibs/macos|linux|windows/<arch>/
#   bash scripts/build_flutter_lib.sh --android    # staticLibs/android/<abi>/
#   bash scripts/build_flutter_lib.sh --ios | --ios-sim
#   bash scripts/build_flutter_lib.sh --mobile     # macOS: iOS+Android；其它 host: 仅 Android
#   bash scripts/build_flutter_lib.sh --all
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
APP="$(cd "$ROOT/../../../client/app" && pwd)"
export CARGO_TARGET_DIR="${CARGO_TARGET_DIR:-$ROOT/target}"
# Flutter 静态库仅需 P2P/ICE gRPC，勿用 host（会链入 FFmpeg，macOS/Android 链接失败）
# Flutter 静态库：Viewer ICE + HTTP/2 客户端（无缩略图编码）
RB_LIB_FEATURES="${RB_LIB_FEATURES:-viewer-ffi-slim}"
MOBILE_FEATURES="${MOBILE_FEATURES:-viewer-ffi-slim}"
# Android/iOS 用 release-small（opt-level=z）；桌面 host 仍用 release
RB_CARGO_PROFILE="${RB_CARGO_PROFILE:-release}"
# Android 默认 Rust cdylib → dynLibs/android/<abi>/librb.so（无 CMake）
RB_ANDROID_FORMAT="${RB_ANDROID_FORMAT:-cdylib}"
MODE="${1:-host}"
# 默认仅 arm64-v8a；模拟器 x86_64：RB_ANDROID_ABIS=x86 bash scripts/build_flutter_lib.sh --android
ANDROID_ABIS=(
  "aarch64-linux-android:arm64-v8a"
)
if [[ "${RB_ANDROID_ABIS:-}" == x86 || "${RB_ANDROID_ABIS:-}" == all ]]; then
  ANDROID_ABIS+=("x86_64-linux-android:x86_64")
fi
IOS_TARGET="${IOS_TARGET:-aarch64-apple-ios}"

# 与 client/app/macos Podfile `platform :osx` 一致
RB_MACOS_DEPLOYMENT_TARGET="${RB_MACOS_DEPLOYMENT_TARGET:-10.15}"

log() { echo "[rb-static] $*"; }
die() { echo "[rb-static] ERROR: $*" >&2; exit 1; }

is_windows_host() {
  case "$(uname -s)" in
    MINGW*|MSYS*|CYGWIN*) return 0 ;;
  esac
  [[ "${OS:-}" == Windows_NT ]]
}

normalize_path() {
  echo "${1//\\//}"
}

# MSYS/Git Bash 下 C:/D: 路径转为 /c/、/d/，否则 -d / ls 失败
to_bash_path() {
  local p="$(normalize_path "$1")"
  if is_windows_host && [[ "$p" =~ ^([A-Za-z]):/(.*)$ ]]; then
    echo "/${BASH_REMATCH[1],,}/${BASH_REMATCH[2]}"
  else
    echo "$p"
  fi
}

read_android_sdk_from_local_props() {
  local props="$APP/android/local.properties" line sdk
  [[ -f "$props" ]] || return 1
  line="$(grep -E '^sdk\.dir=' "$props" 2>/dev/null | tail -1 || true)"
  [[ -n "$line" ]] || return 1
  sdk="${line#sdk.dir=}"
  sdk="${sdk//\\//}"
  to_bash_path "$sdk"
}

# 接受 NDK 根目录或 SDK 根目录（含 ndk/ 子目录）
resolve_ndk_home() {
  local raw="$1" p ndk
  [[ -n "$raw" ]] || return 1
  p="$(to_bash_path "$raw")"
  if [[ -d "$p/toolchains/llvm/prebuilt" ]]; then
    echo "$p"
    return 0
  fi
  if [[ -d "$p/ndk" ]]; then
    ndk="$(ls -d "$p/ndk/"* 2>/dev/null | sort -V | tail -1 || true)"
    if [[ -n "$ndk" && -d "$ndk/toolchains/llvm/prebuilt" ]]; then
      echo "$ndk"
      return 0
    fi
  fi
  return 1
}

default_android_sdk() {
  local sdk props_sdk=""
  props_sdk="$(read_android_sdk_from_local_props || true)"
  if [[ -n "$props_sdk" ]]; then
    echo "$props_sdk"
    return
  fi
  case "$(uname -s)" in
    Darwin) sdk="$HOME/Library/Android/sdk" ;;
    MINGW*|MSYS*|CYGWIN*)
      if [[ -n "${LOCALAPPDATA:-}" ]]; then
        sdk="$(normalize_path "$LOCALAPPDATA/Android/Sdk")"
      fi
      ;;
    *) sdk="$HOME/Android/Sdk" ;;
  esac
  to_bash_path "$sdk"
}

find_ndk() {
  local sdk ndk cand
  for cand in "${ANDROID_NDK_HOME:-}" "${ANDROID_NDK_ROOT:-}" "${ANDROID_NDK_LATEST_HOME:-}"; do
    if ndk="$(resolve_ndk_home "$cand" 2>/dev/null)"; then
      echo "$ndk"
      return
    fi
  done
  sdk="$(to_bash_path "${ANDROID_HOME:-${ANDROID_SDK_ROOT:-}}")"
  if [[ -z "$sdk" || ! -d "$sdk" ]]; then
    sdk="$(default_android_sdk)"
  fi
  ndk="$(ls -d "$sdk/ndk/"* 2>/dev/null | sort -V | tail -1 || true)"
  if [[ -n "$ndk" && -d "$ndk/toolchains/llvm/prebuilt" ]]; then
    echo "$(to_bash_path "$ndk")"
    return
  fi
  die "未找到 NDK（设置 ANDROID_NDK_HOME=<sdk>/ndk/<ver> 或 ANDROID_HOME，当前 sdk=$sdk）"
}

apply_macos_deployment_target() {
  export MACOSX_DEPLOYMENT_TARGET="$RB_MACOS_DEPLOYMENT_TARGET"
  local min="-mmacosx-version-min=${RB_MACOS_DEPLOYMENT_TARGET}"
  export RUSTFLAGS="${RUSTFLAGS:+$RUSTFLAGS }-C link-arg=${min}"
  export CFLAGS="${CFLAGS:+$CFLAGS }${min}"
  export CXXFLAGS="${CXXFLAGS:+$CXXFLAGS }${min}"
}

host_arch_slug() {
  case "$(uname -m)" in
    arm64|aarch64) echo arm64 ;;
    x86_64|amd64) echo x86_64 ;;
    i386|i686) echo x86 ;;
    *) echo "$(uname -m)" ;;
  esac
}

win_arch_slug() {
  case "$(uname -m)" in
    x86_64|amd64) echo x64 ;;
    aarch64|arm64) echo arm64 ;;
    i386|i686) echo x86 ;;
    *) echo x64 ;;
  esac
}

ndk_prebuilt_host() {
  local ndk="$1" base="$ndk/toolchains/llvm/prebuilt"
  for d in windows-x86_64 darwin-aarch64 darwin-x86_64 linux-x86_64; do
    [[ -d "$base/$d" ]] && echo "$d" && return
  done
  die "NDK llvm prebuilt 未找到: $base"
}

target_env_key() {
  echo "$1" | tr '[:lower:]' '[:upper:]' | tr '-' '_'
}

target_cc_key() {
  echo "$1" | tr '-' '_'
}

# NDK clang：Windows 上 Cargo 必须指向 .cmd，否则 os error 193
ndk_clang() {
  local tc="$1" abi="$2" api="$3" kind="${4:-clang}"
  local base="$tc/${abi}${api}-${kind}"
  if is_windows_host; then
    [[ -f "${base}.cmd" ]] && echo "${base}.cmd" && return
    [[ -f "${base}.exe" ]] && echo "${base}.exe" && return
  fi
  echo "$base"
}

android_build_env() {
  local triple="$1" ndk api abi key
  ndk="$(find_ndk)"; api="${ANDROID_API_LEVEL:-24}"
  local tc="$ndk/toolchains/llvm/prebuilt/$(ndk_prebuilt_host "$ndk")/bin"
  case "$triple" in
    aarch64-linux-android) abi=aarch64-linux-android ;;
    x86_64-linux-android) abi=x86_64-linux-android ;;
    *) die "未知 triple: $triple（仅 arm64-v8a / 可选 x86_64）" ;;
  esac
  local clang cxx cckey cargkey
  clang="$(ndk_clang "$tc" "$abi" "$api" clang)"
  cxx="$(ndk_clang "$tc" "$abi" "$api" "clang++")"
  cckey="$(target_cc_key "$triple")"
  cargkey="$(target_env_key "$triple")"
  export "CARGO_TARGET_${cargkey}_LINKER"="$clang"
  export "CC_${cckey}"="$clang"
  export "CXX_${cckey}"="$cxx"
  export "AR_${cckey}"="$tc/llvm-ar"
  export "RANLIB_${cckey}"="$tc/llvm-ranlib"
  local ffmpeg_key="FFMPEG_DIR_${triple//-/_}"
  local ffmpeg_dir="${!ffmpeg_key:-${FFMPEG_DIR_ANDROID:-}}"
  if [[ -n "$ffmpeg_dir" ]]; then
    export PKG_CONFIG_ALLOW_CROSS=1 PKG_CONFIG_PATH="$ffmpeg_dir/lib/pkgconfig" FFMPEG_DIR="$ffmpeg_dir"
  fi
}

ios_subdir() {
  case "$1" in
    *-sim) echo sim ;;
    *) echo device ;;
  esac
}

ios_build_env() {
  local triple="$1" sdk=iphoneos sysroot clang ios_flags cargkey cckey
  command -v xcrun >/dev/null 2>&1 || die "iOS 编译需 macOS 且已安装 Xcode（xcrun）"
  [[ "$triple" == *-sim* ]] && sdk=iphonesimulator
  sysroot="$(xcrun --sdk "$sdk" --show-sdk-path)"
  clang="$(xcrun --sdk "$sdk" --find clang 2>/dev/null || xcrun --find clang)"
  ios_flags="-arch arm64 -isysroot ${sysroot}"
  cargkey="$(target_env_key "$triple")"
  cckey="$(target_cc_key "$triple")"
  export "CC_${cckey}"="$clang"
  export "CXX_${cckey}"="$clang"
  export "CFLAGS_${cckey}=${ios_flags}"
  export "CXXFLAGS_${cckey}=${ios_flags}"
  export "CARGO_TARGET_${cargkey}_LINKER=${clang}"
  export "CARGO_TARGET_${cargkey}_RUSTFLAGS=-C link-arg=-isysroot -C link-arg=${sysroot} -C link-arg=-arch -C link-arg=arm64"
  local ffmpeg_key="FFMPEG_DIR_${triple//-/_}"
  local ffmpeg_dir="${!ffmpeg_key:-${FFMPEG_DIR_IOS:-}}"
  if [[ -n "$ffmpeg_dir" ]]; then
    export PKG_CONFIG_ALLOW_CROSS=1 PKG_CONFIG_PATH="$ffmpeg_dir/lib/pkgconfig" FFMPEG_DIR="$ffmpeg_dir"
  fi
}

ensure_rust_targets() {
  local t; for t in "$@"; do
    rustup target list --installed | grep -q "^${t}\$" || { log "rustup target add $t"; rustup target add "$t"; }
  done
}

cargo_rb() {
  local triple="${1:-}" log_suffix="${1:-host}" features="${3:-$RB_LIB_FEATURES}"
  local profile="${4:-$RB_CARGO_PROFILE}"
  local log="$CARGO_TARGET_DIR/build-${log_suffix}.log"
  mkdir -p "$CARGO_TARGET_DIR"
  ( cd "$ROOT"
    if [[ -n "$triple" ]]; then
      cargo build --profile "$profile" --target "$triple" --no-default-features --features "$features"
    else
      cargo build --profile "$profile" --no-default-features --features "$features"
    fi
  ) >"$log" 2>&1 || { cat "$log" >&2; die "cargo 失败: $log_suffix"; }
}

# 仅编 cdylib（不编 staticlib），Android 动态库
cargo_rb_cdylib() {
  local triple="$1" log_suffix="$1" features="${2:-$MOBILE_FEATURES}" profile="${3:-$RB_CARGO_PROFILE}"
  local log="$CARGO_TARGET_DIR/build-${log_suffix}-cdylib.log"
  mkdir -p "$CARGO_TARGET_DIR"
  ( cd "$ROOT"
    cargo rustc --profile "$profile" --target "$triple" --lib \
      --no-default-features --features "$features" --crate-type cdylib
  ) >"$log" 2>&1 || { cat "$log" >&2; die "cargo cdylib 失败: $log_suffix"; }
}

static_artifact() {
  local triple="${1:-}" profile="${2:-release}"
  if [[ -z "$triple" ]]; then
    if is_windows_host; then
      echo "$CARGO_TARGET_DIR/$profile/rb.lib"
    else
      case "$(uname -s)" in
        Darwin|Linux) echo "$CARGO_TARGET_DIR/$profile/librb.a" ;;
        *) die "未知 host OS: $(uname -s)" ;;
      esac
    fi
  else
    echo "$CARGO_TARGET_DIR/$triple/$profile/librb.a"
  fi
}

cdylib_artifact() {
  local triple="$1" profile="${2:-release-small}"
  echo "$CARGO_TARGET_DIR/$triple/$profile/librb.so"
}

install_cdylib() {
  local src="$1" dest_dir="$2"
  [[ -f "$src" ]] || die "缺少 $src"
  mkdir -p "$dest_dir"
  cp -f "$src" "$dest_dir/librb.so"
}

install_static() {
  local src="$1" dest_dir="$2" dest_name="${3:-librb.a}"
  [[ -f "$src" ]] || die "缺少 $src"
  mkdir -p "$dest_dir"
  cp -f "$src" "$dest_dir/$dest_name"
}

build_host() {
  local arch a
  if is_windows_host; then
    arch="$(win_arch_slug)"
    log "windows $arch (transport, rb.lib)"
    cargo_rb ""
    a="$(static_artifact)"
    install_static "$a" "$APP/staticLibs/windows/$arch" "rb.lib"
    log "ok staticLibs/windows/$arch/rb.lib"
    return
  fi
  arch="$(host_arch_slug)"
  if [[ "$(uname -s)" == Darwin ]]; then
    apply_macos_deployment_target
    log "macOS deployment target $RB_MACOS_DEPLOYMENT_TARGET"
  fi
  cargo_rb ""
  a="$(static_artifact)"
  case "$(uname -s)" in
    Darwin)
      install_static "$a" "$APP/staticLibs/macos/$arch"
      log "ok staticLibs/macos/$arch/librb.a"
      ;;
    Linux)
      install_static "$a" "$APP/staticLibs/linux/$arch"
      log "ok staticLibs/linux/$arch/librb.a"
      ;;
    *) die "host 不支持: $(uname -s)" ;;
  esac
}

build_android_one() {
  local triple="$1" abi="$2"
  android_build_env "$triple"
  if [[ "$RB_ANDROID_FORMAT" == cdylib ]]; then
    log "android $abi ($triple) cdylib → dynLibs/android/$abi/librb.so"
    cargo_rb_cdylib "$triple" "$MOBILE_FEATURES" "$RB_CARGO_PROFILE"
    install_cdylib "$(cdylib_artifact "$triple" "$RB_CARGO_PROFILE")" "$APP/dynLibs/android/$abi"
    log "ok dynLibs/android/$abi/librb.so"
    return
  fi
  log "android $abi ($triple) staticlib → staticLibs/android/$abi/librb.a"
  cargo_rb "$triple" "$triple" "$MOBILE_FEATURES" "$RB_CARGO_PROFILE"
  install_static "$(static_artifact "$triple" "$RB_CARGO_PROFILE")" "$APP/staticLibs/android/$abi"
  log "ok staticLibs/android/$abi/librb.a"
}

build_android_parallel() {
  local pids=() line triple abi ok=0 p
  for line in "${ANDROID_ABIS[@]}"; do
    triple="${line%%:*}"; abi="${line##*:}"
    ( build_android_one "$triple" "$abi" ) &
    pids+=($!)
  done
  for p in "${pids[@]}"; do wait "$p" || ok=1; done
  [[ "$ok" -eq 0 ]] || die "Android 构建失败"
}

build_ios() {
  local triple="${1:-$IOS_TARGET}" sub out
  [[ "$(uname -s)" == Darwin ]] || die "iOS 需在 macOS 编译"
  sub="$(ios_subdir "$triple")"
  out="$APP/staticLibs/ios/$sub"
  ios_build_env "$triple"
  log "ios ($triple) → staticLibs/ios/$sub/librb.a"
  cargo_rb "$triple" "$triple" "$MOBILE_FEATURES"
  install_static "$(static_artifact "$triple")" "$out"
  log "ok staticLibs/ios/$sub/librb.a"
}

build_mobile_parallel() {
  local triples=() line
  for line in "${ANDROID_ABIS[@]}"; do triples+=("${line%%:*}"); done
  if [[ "$(uname -s)" == Darwin ]]; then
    triples=("$IOS_TARGET" "${triples[@]}")
    ensure_rust_targets "${triples[@]}"
    log "并行: iOS device + Android(${#ANDROID_ABIS[@]})"
    local ios_pids=() ok=0 p
    ( build_ios "$IOS_TARGET" ) &
    ios_pids+=($!)
    if [[ "${IOS_BUILD_SIM:-}" == 1 ]]; then
      ( build_ios "aarch64-apple-ios-sim" ) &
      ios_pids+=($!)
    fi
    ( build_android_parallel ) &
    ap=$!
    for p in "${ios_pids[@]}"; do wait "$p" || ok=1; done
    wait "$ap" || ok=1
    [[ "$ok" -eq 0 ]] || die "mobile 失败"
    log "完成 → staticLibs/android/* + staticLibs/ios/device|sim/librb.a"
    return
  fi
  ensure_rust_targets "${triples[@]}"
  log "并行: Android(${#ANDROID_ABIS[@]})（iOS 需 macOS）"
  build_android_parallel
  log "完成 → dynLibs/android/*/librb.so"
}

build_android_only() {
  local triples=() line
  for line in "${ANDROID_ABIS[@]}"; do triples+=("${line%%:*}"); done
  ensure_rust_targets "${triples[@]}"
  build_android_parallel
  log "完成 → dynLibs/android/*/librb.so"
}

case "$MODE" in
  host | "") build_host ;;
  --host) build_host ;;
  --android | android) build_android_only ;;
  --ios | ios) build_ios "$IOS_TARGET" ;;
  --ios-sim | ios-sim) build_ios "aarch64-apple-ios-sim" ;;
  --mobile | mobile) build_mobile_parallel ;;
  --all | all)
    if [[ "$(uname -s)" == Darwin ]]; then ( build_host ) & hp=$!; build_mobile_parallel; wait "$hp"
    else build_host; build_mobile_parallel; fi ;;
  -h | --help)
    sed -n '2,7p' "$0"
    echo "RB_LIB_FEATURES=viewer-ffi-slim（默认）；桌面含 Agent 用 transport-slim 或 host"
    echo "MOBILE_FEATURES=viewer-ffi-slim；Windows/Linux host 可 --android；iOS 需 macOS"
    echo "Windows host → staticLibs/windows/x64/rb.lib（MSVC 或 .cargo gnu linker）"
    ;;
  *) die "参数: host | --android | --ios | --ios-sim | --mobile | --all" ;;
esac
exit 0
