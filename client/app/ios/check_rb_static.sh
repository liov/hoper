#!/bin/sh
# 校验 staticLibs/ios 静态库存在（链接由 rb.xcconfig 完成）
set -e
BASE="${SRCROOT}/../staticLibs/ios"
if echo "${PLATFORM_NAME:-}" | grep -qi simulator; then
  A="${BASE}/sim/librb.a"
else
  A="${BASE}/device/librb.a"
fi
if [ ! -f "$A" ]; then
  echo "error: 缺少 $A，请先: cd server/rust/rb && bash scripts/build_flutter_lib.sh --mobile" >&2
  exit 1
fi
