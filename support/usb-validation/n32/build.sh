#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/../../.." && pwd)
OUT=${OUT:-$REPO_ROOT/build/n32probe_usb}

if [ -z "${GO:-}" ]; then
  if command -v go >/dev/null 2>&1; then
    GO=$(command -v go)
  elif [ -x "${HOME}/.config/opencode/tools/go/bin/go" ]; then
    GO="${HOME}/.config/opencode/tools/go/bin/go"
  else
    echo "go not found; set GO or install Go (OpenCode tools path is optional)" >&2
    exit 1
  fi
fi
export GO

"$REPO_ROOT/scripts/bootstrap-go-flash.sh"

mkdir -p "$OUT"
cp "$SCRIPT_DIR/flashforge_init.sh" "$OUT/flashforge_init.sh"
chmod +x "$OUT/flashforge_init.sh"

cd "$REPO_ROOT/deps/go-flash"

echo "building ARM binaries"
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 "$GO" build -o "$OUT/n32probe_arm" ./cmd/n32probe
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 "$GO" build -o "$OUT/fbfill_arm" ./cmd/fbfill

echo "building MIPS little-endian binaries"
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat "$GO" build -o "$OUT/n32probe_mips" ./cmd/n32probe
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat "$GO" build -o "$OUT/fbfill_mips" ./cmd/fbfill

chmod +x "$OUT"/n32probe_* "$OUT"/fbfill_*
(cd "$OUT" && shasum -a 256 flashforge_init.sh n32probe_arm fbfill_arm n32probe_mips fbfill_mips > SHA256SUMS)

echo "package written to $OUT"
