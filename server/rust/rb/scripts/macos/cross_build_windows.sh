#!/usr/bin/env bash
# Windows x64 Agent：rb.exe（macOS/Linux 上 zig 交叉编译）
# 依赖：rustup target add x86_64-pc-windows-gnu && brew install mingw-w64
# FFmpeg：默认从 MSYS2 UCRT64 包拉取到 .ffmpeg-win64（可 export FFMPEG_DIR= 覆盖）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
export CARGO_TARGET_DIR="$ROOT/target"
TARGET=x86_64-pc-windows-gnu
OUT="$CARGO_TARGET_DIR/$TARGET/release/rb.exe"
FFMPEG_DIR="${FFMPEG_DIR:-$ROOT/.ffmpeg-win64}"
MSYS_PKG=mingw-w64-ucrt-x86_64-ffmpeg-7.1-2-any.pkg.tar.zst

ensure_ffmpeg_msys() {
  if [[ -f "$FFMPEG_DIR/include/libavcodec/avcodec.h" && -f "$FFMPEG_DIR/lib/libavcodec.a" ]]; then
    return 0
  fi
  echo "Fetching MSYS2 FFmpeg dev → $FFMPEG_DIR"
  mkdir -p "$ROOT/.cache/msys2-ffmpeg"
  local cache="$ROOT/.cache/msys2-ffmpeg/$MSYS_PKG"
  if [[ ! -f "$cache" ]]; then
    curl -fsSL -o "$cache" "https://repo.msys2.org/mingw/ucrt64/$MSYS_PKG"
  fi
  rm -rf "$ROOT/.cache/msys2-ffmpeg/extract"
  mkdir -p "$ROOT/.cache/msys2-ffmpeg/extract"
  tar -xf "$cache" -C "$ROOT/.cache/msys2-ffmpeg/extract"
  rm -rf "$FFMPEG_DIR"
  mkdir -p "$FFMPEG_DIR"
  cp -R "$ROOT/.cache/msys2-ffmpeg/extract/ucrt64/include" "$FFMPEG_DIR/"
  cp -R "$ROOT/.cache/msys2-ffmpeg/extract/ucrt64/lib" "$FFMPEG_DIR/"
  cp -R "$ROOT/.cache/msys2-ffmpeg/extract/ucrt64/bin" "$FFMPEG_DIR/" 2>/dev/null || true
  # 仅保留 .dll.a 导入库，避免链进巨型静态 .a 及大量传递依赖
  for f in "$FFMPEG_DIR/lib"/*.a; do
    [[ "$f" == *.dll.a ]] && continue
    rm -f "$f"
  done
}

cd "$ROOT"
rustup target add "$TARGET" 2>/dev/null || true
ensure_ffmpeg_msys
export FFMPEG_DIR
export PKG_CONFIG_ALLOW_CROSS=1
export PKG_CONFIG_PATH="$FFMPEG_DIR/lib/pkgconfig"

echo "CARGO_TARGET_DIR=$CARGO_TARGET_DIR"
echo "FFMPEG_DIR=$FFMPEG_DIR"
cargo build --release \
  --no-default-features --features host \
  --bin rb --target "$TARGET"

DIST="$ROOT/dist/windows-x64"
mkdir -p "$DIST"
cp -f "$OUT" "$DIST/"
cp -f "$FFMPEG_DIR/bin"/av*.dll "$FFMPEG_DIR/bin"/sw*.dll "$FFMPEG_DIR/bin"/postproc*.dll "$DIST/" 2>/dev/null || true
echo "=> Windows x64: $OUT"
echo "=> 发布目录（含 DLL）: $DIST/"
file "$OUT" 2>/dev/null || ls -la "$OUT"
