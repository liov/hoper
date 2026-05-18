#!/usr/bin/env bash
# 构建供 Flutter 加载的 librfv（ICE viewer + rb_agent_run + ffmpeg 数据面）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="${ROOT}/../../../client/app/libraries"
mkdir -p "$OUT"
cd "$ROOT"
cargo build --release --features host
case "$(uname -s)" in
  Darwin) cp -f target/release/librfv.dylib "$OUT/librfv.dylib" ;;
  Linux) cp -f target/release/librfv.so "$OUT/librfv.so" ;;
  *) echo "copy librfv manually for this OS" >&2; exit 1 ;;
esac
echo "ok: $OUT"
