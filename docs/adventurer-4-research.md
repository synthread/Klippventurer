# Adventurer 4 Research Notes

These notes track evidence for porting stock Adventurer 4 firmware behavior to a native Klipper config. Treat anything marked tentative as unsafe until verified against Ghidra decompilation, board tracing, or a live eMMC dump.

## Firmware Inputs

Local firmware archives currently available under `.firmware-analysis/tmp/`:

- `ad4lite-software-2.0.5.tar.xz`
- `ad4lite-control-2.1.0.tar.xz`
- `ad4lite-library-1.0.0.tar.xz`
- `ad4lite-kernel-1.0.0.tar.xz`
- `ad4pro114.tgz`
- `ad4pro116.tgz`

Extracted working directories live under `.firmware-analysis/extracted/`.

Use `partialzip` for any further remote firmware packages when only selected archive members are needed.

## Confirmed From Static Analysis

- Both Lite and Pro stock UI binaries are unstripped ARM Linux ELF binaries with symbols and debug path strings.
- Lite software binary: `.firmware-analysis/extracted/ad4lite-software-2.0.5/adventurer4-arm`.
- Pro software binary: `.firmware-analysis/extracted/ad4pro116-software-1.1.6/adventurer4-arm`.
- Lite MCU control image: `Adventurer4-IAP-2.1.0-20210913.hex`.
- Pro MCU control images: `Adventurer4-NA-IAP.hex` and `Adventurer4-HD-IAP-2.1.0-20220411.hex`.
- Stock UI reports Adventurer 4 build volume as `220X200X250` / `220 X 200 X 250`.
- Stock app reads/writes `/opt/finder_rush/config/config-AD4.json`.
- `config-AD4.json` is not present in the extracted update payloads, so it is likely generated on-device, factory-installed on eMMC, or restored from another non-update partition.
- Lite has nine-point auto calibration symbols/strings: `Config::getAutoPoint1()` through `Config::getAutoPoint9()`.
- Pro has thirty-point auto calibration symbols/strings: `Config::getAutoPoint1()` through `Config::getAutoPoint30()` plus `30 Point Auto Calibration`.
- Both variants expose Z/nozzle calibration state names: `buildZOffset`, `m_ZOffset=`, `extruderCalibrationValue`, and `nozzleTempDifferentValue`.

## Tentative Printer.cfg Deltas From Adventurer 3

Current `config/printer.cfg` is Adventurer 3/SZ16-oriented. Tentative Adventurer 4 changes compared to that baseline:

- Build area changes from roughly `161 x 151 x 155` to `220 x 200 x 250`.
- `[stepper_x] position_max` should become approximately `220` after homing-offset validation.
- `[stepper_y] position_max` should become approximately `200` after endstop-origin validation.
- `[stepper_z] position_max` should become approximately `250` after validating usable Z travel and top-end safety margin.
- `[bed_mesh] mesh_max` should expand near the AD4 bed limits, but keep edge margins until probe/nozzle offsets are known.
- Pro and Lite should not share one mesh profile: Lite appears to use 9-point calibration, Pro appears to use 30-point calibration.
- `nozzle_diameter` should not inherit AD3's `0.3` by default; AD4 stock hardware is expected to use a different nozzle family and must be confirmed per toolhead.
- MCU pin mappings, TMC mux/select pins, thermistor curve, fan pins, LED pins, and endstop pins remain unconfirmed for AD4 and should not be copied blindly from AD3.

## Include-Based Config Generator Plan

Replace the single monolithic `config/printer.cfg` with generated printer entrypoints that compose reusable include fragments.

Proposed layout:

```text
config/
  generated/
    printer-adventurer3.cfg
    printer-adventurer3-pro.cfg
    printer-adventurer4-lite.cfg
    printer-adventurer4-pro.cfg
  parts/
    common/
      fluidd.cfg
      kamp.cfg
      mcu-uart-pi.cfg
    motion/
      cartesian-common.cfg
      adventurer3-geometry.cfg
      adventurer4-geometry.cfg
    boards/
      sz16-stm32f103.cfg
      adventurer4-stock-mainboard.cfg
    toolheads/
      adv3-stock-hotend.cfg
      adv4-stock-hotend.cfg
    calibration/
      adv3-bed-mesh.cfg
      adv4-lite-bed-mesh.cfg
      adv4-pro-bed-mesh.cfg
    macros/
      adv3-macros.cfg
      adventurer-common-macros.cfg
```

Generated Adventurer 4 Lite entrypoint should initially include:

```ini
[include ../parts/common/mcu-uart-pi.cfg]
[include ../parts/motion/cartesian-common.cfg]
[include ../parts/motion/adventurer4-geometry.cfg]
[include ../parts/boards/adventurer4-stock-mainboard.cfg]
[include ../parts/toolheads/adv4-stock-hotend.cfg]
[include ../parts/calibration/adv4-lite-bed-mesh.cfg]
[include ../parts/common/fluidd.cfg]
[include ../parts/common/kamp.cfg]
```

Generated Adventurer 4 Pro entrypoint should use the same base fragments but swap in `adv4-pro-bed-mesh.cfg` and any Pro-specific board/toolhead deltas discovered from firmware or hardware inspection.

Generator input should be data-first, not template-only. Suggested schema fields:

- `machine_id`: `adventurer3`, `adventurer3_pro`, `adventurer4_lite`, `adventurer4_pro`
- `board`: MCU, serial, bootloader, pin map fragment
- `geometry`: X/Y/Z min/max, endstop positions, homing speeds
- `motion`: max velocity, acceleration, Z velocity/accel, square-corner velocity
- `toolhead`: extruder pins, heater pin, thermistor, nozzle diameter, rotation distance
- `bed`: heater pin, thermistor, max temp, mesh extents, probe count
- `drivers`: TMC type, UART topology, currents, sense resistor
- `features`: camera, chamber fan, LED, filament runout, display/touch, buzzer

## Next Static-Analysis Targets

- Import `adventurer4-arm` into Ghidra via the project-local MCP config in `opencode.json`.
- Decompile `Config::readFromConfig()`, `Config::initRootConfig()`, and `Config::writeToConfig()` to recover the expected `config-AD4.json` schema and default values.
- Decompile `SerialObject::serialSendBase()`, `SerialObject::serialCheckCode()`, and `ServerListener::cmd_G28()` to identify stock homing and motion commands.
- Convert the Intel HEX MCU images to raw binaries and inspect/decompile for pin setup, thermistor tables, motion constants, and driver configuration.
- Obtain a live eMMC `/opt/finder_rush/config/config-AD4.json` if possible; this is likely the fastest path to factory calibration defaults.
