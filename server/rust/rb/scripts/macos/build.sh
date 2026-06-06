#!/usr/bin/env bash
# 远程相册：编译到 crate 内 target/（避免沙箱 CARGO_TARGET_DIR 指向缓存导致跑旧 rb）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
export CARGO_TARGET_DIR="${CARGO_TARGET_DIR:-$ROOT/target}"
cd "$ROOT"
echo "CARGO_TARGET_DIR=$CARGO_TARGET_DIR"
cargo build --release --features host,daemon
# Flutter 静态库（transport，无 FFmpeg）；移动平台: bash scripts/build_flutter_lib.sh --mobile
RB_LIB_FEATURES=transport-slim bash "$ROOT/scripts/build_flutter_lib.sh" "${RB_LIB_MODE:-}"
echo "rb:      $CARGO_TARGET_DIR/release/rb"
echo "daemon:   $CARGO_TARGET_DIR/release/rb-daemon"
echo "librb.a: $ROOT/../../../client/app/staticLibs/"
