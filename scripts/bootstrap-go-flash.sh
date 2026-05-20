#!/bin/sh
# Clone or update deps/go-flash for local and CI builds.
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
DEPS_DIR="$REPO_ROOT/deps/go-flash"
GO_FLASH_REPO="${GO_FLASH_REPO:-https://gitlab.com/synthread/libs/go-flash.git}"
GO_FLASH_REF="${GO_FLASH_REF:-main}"
GO_FLASH_OVERLAY="${GO_FLASH_OVERLAY:-$REPO_ROOT/support/usb-validation/n32/go-flash-overlay}"

_apply_overlay() {
  if [ ! -f "$GO_FLASH_OVERLAY/go.mod" ]; then
    return 0
  fi
  echo "Applying go-flash overlay from $GO_FLASH_OVERLAY"
  rsync -a --exclude .git "$GO_FLASH_OVERLAY/" "$DEPS_DIR/"
}

if [ -d "$DEPS_DIR/.git" ]; then
  echo "Updating existing go-flash at $DEPS_DIR"
  git -C "$DEPS_DIR" fetch origin --tags
  git -C "$DEPS_DIR" checkout "$GO_FLASH_REF" 2>/dev/null || git -C "$DEPS_DIR" checkout "origin/$GO_FLASH_REF"
  git -C "$DEPS_DIR" pull --ff-only origin "$GO_FLASH_REF" 2>/dev/null || true
  _apply_overlay
  exit 0
fi

if [ -f "$DEPS_DIR/go.mod" ] && [ ! -d "$DEPS_DIR/.git" ]; then
  echo "go-flash already present at $DEPS_DIR (non-git)"
  _apply_overlay
  exit 0
fi

mkdir -p "$REPO_ROOT/deps"
echo "Cloning go-flash into $DEPS_DIR"
git clone --depth 1 --branch "$GO_FLASH_REF" "$GO_FLASH_REPO" "$DEPS_DIR" 2>/dev/null \
  || git clone --depth 1 "$GO_FLASH_REPO" "$DEPS_DIR"

_apply_overlay
