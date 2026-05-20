#!/bin/sh

cd "$(dirname "$0")" || exit 1

LOG=n32probe.log
exec >"$LOG" 2>&1

echo "n32probe USB validation package"
date 2>/dev/null || true
uname -a 2>/dev/null || true
ARCH=$(uname -m 2>/dev/null || echo unknown)
echo "ARCH=$ARCH"
echo "PWD=$(pwd)"
ls -la 2>/dev/null || true
echo "serial devices:"
ls -l /dev/ttyS* /dev/ttyAMA* /dev/ttyUSB* /dev/ttyACM* 2>/dev/null || true
echo "gpio sysfs:"
ls -ld /sys/class/gpio /sys/class/gpio/gpio* 2>/dev/null || true

paint_fb() {
  COLOR=$1
  if [ -x "$FB" ]; then
    FB_ARGS="--color $COLOR --width ${N32PROBE_FB_WIDTH:-0} --height ${N32PROBE_FB_HEIGHT:-0} --bpp ${N32PROBE_FB_BPP:-0}"
    if [ -n "${N32PROBE_FB_ORDER:-}" ]; then
      FB_ARGS="$FB_ARGS --order ${N32PROBE_FB_ORDER}"
    fi
    "$FB" $FB_ARGS
    return $?
  fi

  if [ ! -w /dev/fb0 ]; then
    echo "framebuffer helper missing and /dev/fb0 is not writable"
    return 1
  fi

  BYTE="\377\000\000\000"
  if [ "$COLOR" = "yellow" ]; then
    BYTE="\377\377\000\000"
  fi
  if [ "$COLOR" = "green" ]; then
    BYTE="\000\377\000\000"
  fi

  # Emergency fallback when fbfill is missing. This writes a conservative
  # 800x480x32bpp-sized solid color buffer using POSIX shell printf.
  i=0
  : >/tmp/n32probe_fbfill.raw
  while [ "$i" -lt 384000 ]; do
    printf "$BYTE" >>/tmp/n32probe_fbfill.raw
    i=$((i + 1))
  done
  cat /tmp/n32probe_fbfill.raw >/dev/fb0
}

FB=./fbfill_arm
PROBE=./n32probe_arm
case "$ARCH" in
  mips|mipsel|mips*)
    FB=./fbfill_mips
    PROBE=./n32probe_mips
    ;;
  arm*|aarch64)
    FB=./fbfill_arm
    PROBE=./n32probe_arm
    ;;
  *)
    echo "unsupported architecture: $ARCH"
    paint_fb red || true
    sync
    sleep 3600
    exit 1
    ;;
esac

if [ ! -x "$PROBE" ]; then
  echo "missing probe binary: $PROBE"
  paint_fb red || true
  sync
  sleep 3600
  exit 1
fi

if [ ! -x "$FB" ]; then
  echo "missing framebuffer helper: $FB"
  paint_fb red || true
  sync
  sleep 3600
  exit 1
fi

MACHINE=${N32PROBE_MACHINE:-auto}
case "$ARCH" in
  mips|mipsel|mips*)
    : "${N32PROBE_MACHINE:=adv3-mips}"
    MACHINE=$N32PROBE_MACHINE
    ;;
esac

echo "Running probe: machine=$MACHINE"
"$PROBE" \
  --machine "$MACHINE" \
  --family n32g45x \
  --tty "${N32PROBE_TTY:-/dev/ttyS1}" \
  --baud "${N32PROBE_BAUD:-115200}" \
  --boot-entry "${N32PROBE_BOOT_ENTRY:-auto}" \
  --boot-delay "${N32PROBE_BOOT_DELAY:-3s}" \
  --scan-ttys "${N32PROBE_SCAN_TTYS:-/dev/ttyS0,/dev/ttyS1,/dev/ttyS2,/dev/ttyAMA0}" \
  --allow-unknown-chip-id \
  --json > n32probe-result.json
STATUS=$?

cat n32probe-result.json

sync
if [ "$STATUS" -eq 0 ]; then
  echo "validation PASS; log ready; writing green framebuffer"
  paint_fb green || STATUS=1
else
  echo "validation WARN/FAIL; log ready; writing yellow framebuffer"
  paint_fb yellow || true
fi
sync

echo "holding init process; status=$STATUS"
sleep 3600
exit "$STATUS"
