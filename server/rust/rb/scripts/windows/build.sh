#!/usr/bin/env bash
set -euo pipefail

SDK=/d/sdk
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

export HTTP_PROXY="${HTTP_PROXY:-${http_proxy:-}}"
export HTTPS_PROXY="${HTTPS_PROXY:-${https_proxy:-}}"

export FFMPEG_DIR=D:/sdk/msys64/ucrt64
export PKG_CONFIG_PATH=D:/sdk/msys64/ucrt64/lib/pkgconfig
export PKG_CONFIG_ALLOW_CROSS=1
export LIBCLANG_PATH=D:/sdk/llvm/21.1.5/bin
export BINDGEN_EXTRA_CLANG_ARGS="--target=x86_64-w64-windows-gnu --sysroot=D:/sdk/msys64/ucrt64 -resource-dir=D:/sdk/msys64/ucrt64/lib/clang/22"
# MSYS bash 的 PATH 须用 /d/ 前缀，D:/ 形式找不到 rustc.exe
export PATH="${SDK}/msys64/ucrt64/bin:${SDK}/llvm/21.1.5/bin:${SDK}/.cargo/bin:${PATH:-}"
export CARGO_HOME=D:/sdk/.cargo
export RUSTUP_HOME=D:/sdk/rust
export CARGO_TARGET_DIR="${ROOT}/target"

cd "$ROOT"
echo "CARGO_TARGET_DIR=$CARGO_TARGET_DIR"
echo "FFMPEG_DIR=$FFMPEG_DIR"
echo "LIBCLANG_PATH=$LIBCLANG_PATH"
echo "rustc: $(rustc --version)"
cargo build --release --features host,daemon
