#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"
# 需：rustup target add x86_64-unknown-linux-gnu && cargo install cargo-zigbuild
cargo zigbuild --release \
  --no-default-features --features daemon \
  --bin rb-daemon \
  --target x86_64-unknown-linux-gnu
OUT="${CARGO_TARGET_DIR:-$ROOT/target}/x86_64-unknown-linux-gnu/release/rb-daemon"
echo "=> Linux x64: $OUT"
file "$OUT"
