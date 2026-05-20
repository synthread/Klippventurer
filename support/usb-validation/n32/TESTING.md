# N32 USB validation test rundown

## Package to test

Build output:

```text
build/n32probe_usb/
```

Copy all files from that directory to the USB root:

- `flashforge_init.sh`
- `n32probe_arm`
- `fbfill_arm`
- `n32probe_mips`
- `fbfill_mips`
- `SHA256SUMS`

## Expected visual result

- Solid green framebuffer: all required validation passed and logs were synced.
- Solid red framebuffer: validation failed or log/sync failed.
- The script holds the init process; it does not shut down the printer.

## Logs to collect

After each run, collect from the USB package directory:

- `n32probe.log`
- `n32probe-result.json`
- `SHA256SUMS`

## N32 Adventurer 3 test

Goal: validate the legacy-MIPS style package still runs and can identify the N32 microcontroller on the mainboard safely. Adventurer 3 is not treated as a 5M-style eboard/mainboard split.

1. Copy package files to USB root.
2. If needed, force the profile by adding this environment before invoking manually, or by editing the script for this run:

   ```sh
   N32PROBE_MACHINE=adv3-mips
   N32PROBE_BOOT_ENTRY=gpio
   ```

3. Run via normal FlashForge USB init path.
4. Expected:
   - `n32probe_mips` is selected on MIPS printers.
   - GPIO boot-entry is attempted only for `adv3-mips` / `gpio` mode.
   - Probe performs sync, GET, GET ID only.
   - Green screen if chip responds; red if not.
5. Send back `n32probe.log` and `n32probe-result.json`.

## ARM Adventurer 3 Pro 2 / Adventurer 4-family test

Goal: validate the ARM N32 microcontroller on the mainboard via vendor-style serial bootloader request. Adventurer 4-family machines are not treated as 5M-style eboard/mainboard split unless later hardware evidence proves otherwise.

1. Copy package files to USB root.
2. Default ARM behavior should be sufficient:

   ```sh
   N32PROBE_MACHINE=arm-n32-single
   N32PROBE_BOOT_ENTRY=serial-request
   N32PROBE_TTY=/dev/ttyS1
   N32PROBE_BAUD=115200
   ```

3. Expected:
   - `n32probe_arm` is selected.
   - script sends `~ \x1c Request Serial Bootloader!! ~` to `/dev/ttyS1` at 230400.
   - probe reopens `/dev/ttyS1` at 115200 and runs sync, GET, GET ID only.
   - green if probe succeeds; red if not.

## Adventurer 5M / 5M Pro test

Goal: validate both MCUs before framebuffer status.

Force 5M profile if auto-detection does not catch it:

```sh
N32PROBE_MACHINE=ad5m
```

Expected mainboard behavior:

- Uses best-effort vendor-native contact/reset command if available:

  ```sh
  NationsCommand -c --pn /dev/ttyS5 -r
  ```

- Does not pass firmware paths or `--fn`.
- Treats exit status 0 as mainboard contact/reset pass.

Expected eboard behavior:

- Sends serial bootloader request to `/dev/ttyS1` at 230400.
- Probes `/dev/ttyS1` at 115200 with sync, GET, GET ID only.
- Eboard probe must pass for overall green.

If mainboard contact passes but eboard probe fails, expect red and logs showing which sub-check failed.

## What we need back

For each machine tested, report:

- exact printer model and firmware version if known
- framebuffer color shown
- whether the printer remained responsive after power cycle
- `n32probe.log`
- `n32probe-result.json`
- any observed sounds/motion/display changes during bootloader request

Do not retry flash/update commands manually. This package is intended to be read-only.
