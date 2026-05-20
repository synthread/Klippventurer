#!/bin/sh
# Deploy a built n32probe USB package to a dev printer over SSH.
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
REPO_ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/../../.." && pwd)
PACKAGE=${PACKAGE:-$REPO_ROOT/build/n32probe_usb}
REMOTE_DIR=${REMOTE_DIR:-/tmp/n32probe_usb}
PRINTER_USER=${PRINTER_USER:-root}

if [ -z "${PRINTER_HOST:-}" ]; then
  echo "Set PRINTER_HOST (or configure a Host in ~/.ssh/config)" >&2
  exit 1
fi

PRINTER_SSH=${PRINTER_SSH:-$PRINTER_USER@$PRINTER_HOST}

if [ ! -f "$PACKAGE/flashforge_init.sh" ]; then
  echo "Package missing at $PACKAGE — run support/usb-validation/n32/build.sh first" >&2
  exit 1
fi

echo "Deploying $PACKAGE -> $PRINTER_SSH:$REMOTE_DIR"
ssh "$PRINTER_SSH" "mkdir -p $REMOTE_DIR"
rsync -av --delete "$PACKAGE/" "$PRINTER_SSH:$REMOTE_DIR/"
ssh "$PRINTER_SSH" "chmod +x $REMOTE_DIR/flashforge_init.sh $REMOTE_DIR/n32probe_* $REMOTE_DIR/fbfill_*" 2>/dev/null || true

if [ "${DEPLOY_RUN:-}" = "1" ]; then
  echo "Running flashforge_init.sh on printer (DEPLOY_RUN=1)"
  ssh -t "$PRINTER_SSH" "cd $REMOTE_DIR && ./flashforge_init.sh"
else
  echo "Upload complete. Run on printer:"
  echo "  ssh -t $PRINTER_SSH 'cd $REMOTE_DIR && ./flashforge_init.sh'"
  echo "Or re-run with DEPLOY_RUN=1"
fi
