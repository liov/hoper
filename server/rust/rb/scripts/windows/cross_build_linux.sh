#!/usr/bin/env bash
# Linux x64 rb-daemon：Windows MSYS bash 上 zig 交叉编译（生产部署用）
# 依赖：rustup target add x86_64-unknown-linux-gnu && cargo install cargo-zigbuild && zig 在 PATH
# 环境变量与 build.sh 对齐；SDK 默认 /d/sdk，可 export SDK= 覆盖
set -euo pipefail

SDK="${SDK:-/d/sdk}"
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
TARGET=x86_64-unknown-linux-gnu
OUT="${CARGO_TARGET_DIR:-$ROOT/target}/$TARGET/release/rb-daemon"

export HTTP_PROXY="${HTTP_PROXY:-${http_proxy:-}}"
export HTTPS_PROXY="${HTTPS_PROXY:-${https_proxy:-}}"
# MSYS bash 的 PATH 须用 /d/ 前缀，D:/ 形式找不到 rustc.exe；ucrt64/bin 含 pacman 安装的 zig
export PATH="${SDK}/msys64/ucrt64/bin:${SDK}/.cargo/bin:${PATH:-}"
export CARGO_HOME="${CARGO_HOME:-D:/sdk/.cargo}"
export RUSTUP_HOME="${RUSTUP_HOME:-D:/sdk/rust}"
export CARGO_TARGET_DIR="${CARGO_TARGET_DIR:-$ROOT/target}"

cd "$ROOT"
rustup target add "$TARGET" 2>/dev/null || true
command -v zig >/dev/null || { echo "error: zig not in PATH (pacman -S mingw-w64-ucrt-x86_64-zig 或加入 zig 目录)"; exit 1; }

echo "CARGO_TARGET_DIR=$CARGO_TARGET_DIR"
echo "rustc: $(rustc --version)"
echo "zig: $(zig version)"
cargo zigbuild --version >/dev/null 2>&1 || { echo "error: cargo install cargo-zigbuild"; exit 1; }

cargo zigbuild --release \
  --no-default-features --features daemon \
  --bin rb-daemon \
  --target "$TARGET"

echo "=> Linux x64: $OUT"
file "$OUT" 2>/dev/null || ls -la "$OUT"
