# Dev deploy over SSH

## Purpose

Faster iteration on **trusted development printers** without swapping USB media each run. This is a developer workflow only; end-user installs remain USB-based until the `v0.3.3` installer prototype milestone.

Constraints: [ADR 0002](../../docs/specs/adr-0002-package-layout-and-rootfs-reuse.md), [ADR 0005](../../docs/specs/adr-0005-runtime-overlay-and-recovery.md) — no broad eMMC mutation; default to read-only probe behavior.

## Prerequisites

- Printer reachable on the LAN (SSH server running; stock or post-bridge).
- SSH key enrolled for the printer user (store key reference in 1Password; **never** commit keys).
- Local build completed:

  ```sh
  support/usb-validation/n32/build.sh
  ```

- Optional env:

  | Variable | Default | Meaning |
  | --- | --- | --- |
  | `PRINTER_HOST` | _(required)_ | SSH host alias or IP |
  | `PRINTER_USER` | `root` | SSH user |
  | `REMOTE_DIR` | `/tmp/n32probe_usb` | Staging directory on printer |
  | `PRINTER_SSH` | `$PRINTER_USER@$PRINTER_HOST` | Full SSH target |

## Quick deploy

```sh
export PRINTER_HOST=adv3-dev   # ~/.ssh/config Host
support/usb-validation/n32/deploy-ssh.sh
```

## Manual steps

```sh
PACKAGE=build/n32probe_usb
REMOTE_DIR=/tmp/n32probe_usb

ssh "$PRINTER_SSH" "mkdir -p $REMOTE_DIR"
rsync -av --delete "$PACKAGE/" "$PRINTER_SSH:$REMOTE_DIR/"
ssh "$PRINTER_SSH" "chmod +x $REMOTE_DIR/flashforge_init.sh $REMOTE_DIR/n32probe_* $REMOTE_DIR/fbfill_*"
ssh "$PRINTER_SSH" "cd $REMOTE_DIR && ./flashforge_init.sh"
```

## Retrieve logs

```sh
scp "$PRINTER_SSH:$REMOTE_DIR/n32probe.log" ./n32probe.log
scp "$PRINTER_SSH:$REMOTE_DIR/n32probe-result.json" ./n32probe-result.json
```

Share logs per [support/usb-validation/n32/TESTING.md](../../support/usb-validation/n32/TESTING.md).

## Safety

- Treat probe packages as **read-only** unless explicitly running a flash-capable build.
- Do not write to eMMC paths outside agreed staging dirs.
- Prefer a dedicated dev printer; do not use production/user machines without explicit consent.

## Related

- USB validation: [support/usb-validation/n32/README.md](../../support/usb-validation/n32/README.md)
- OTA preview channel design: [.lens/schemas/ota-preview-channel.md](../schemas/ota-preview-channel.md)
- Legacy Pi + UART path: [docs/installation.md](../../docs/installation.md) (v0.2.x only)
