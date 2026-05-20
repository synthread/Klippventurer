# N32 USB validation package

This package builds a FlashForge USB-init validation bundle for N32 Adventurer printers. It is read-only for MCU flash: it probes bootloaders and records logs, but does not erase, write, protect, or unprotect flash.

## Build

Requires `go` and `deps/go-flash` (bootstrapped automatically):

```sh
scripts/bootstrap-go-flash.sh   # optional; build.sh runs this
support/usb-validation/n32/build.sh
```

Override `GO` or `GO_FLASH_REPO` / `GO_FLASH_REF` if needed.

N32 probe support is overlaid from [go-flash-overlay/](go-flash-overlay/) (library + `cmd/`) until merged into [synthread/libs/go-flash](https://gitlab.com/synthread/libs/go-flash) `main`.

Output is written to ignored local directory:

```text
build/n32probe_usb/
```

Copy the contents of that directory to the root of a USB drive.

## Runtime behavior

- Logs to `n32probe.log` and `n32probe-result.json` on the USB package directory.
- Writes solid green pixels to `/dev/fb0` only after validation succeeds and logs are synced.
- Writes solid red pixels to `/dev/fb0` on validation failure.
- Holds the init process instead of shutting down.

## Target behavior

- MIPS Adventurer 3 defaults to legacy GPIO boot-entry plus `/dev/ttyS1` probe.
- ARM Adventurers default to N32 and `/dev/ttyS1`.
- ARM bootloader entry uses the FlashForge serial request at 230400 baud before probing at 115200.
- Adventurer 3 and 4-family machines validate the microcontroller on the mainboard; they are not treated as 5M-style dual-MCU/eboard systems.
- Adventurer 5M can be forced with `N32PROBE_MACHINE=ad5m`; it validates the `/dev/ttyS5` mainboard using the stock read-only `NationsCommand -c --pn /dev/ttyS5 -r` contact/reset path when available, then validates the `/dev/ttyS1` eboard probe.

## Useful overrides

```sh
N32PROBE_MACHINE=ad5m
N32PROBE_TTY=/dev/ttyS1
N32PROBE_BAUD=115200
N32PROBE_BOOT_ENTRY=serial-request
```

Share back `n32probe.log`, `n32probe-result.json`, and `SHA256SUMS` after testing. Do not include unrelated printer logs unless requested.
