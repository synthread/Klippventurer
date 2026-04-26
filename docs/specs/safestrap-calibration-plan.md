# Safestrap And Stock Calibration Import Plan

## Purpose

This document compares the current install/runtime approaches present in the workspace and proposes a safer Adventurer 3 USB installer path built from the best parts of each.

The target outcome is a stock-like `safestrap` flow for Adventurer 3-family printers that:

- installs from USB without permanent helper hardware for supported variants,
- fails closed when model or MCU evidence is weak,
- preserves a reliable path back to stock boot,
- captures stock calibration state before switching runtime,
- guides the user through a touchscreen OOBE before first Klipper print.

## Evidence Summary

### `deps/sl-print`

Current approach:

- USB entrypoint is `deps/sl-print/dist/flashforge_init.sh`.
- It unpacks `slp-update.tar.gz` into `/opt/sl-print`.
- It backs up `/etc/rc.local` to `/etc/rc.local-slbkp`.
- It replaces `/etc/rc.local` with `deps/sl-print/dist/rc.local`.
- Replacement `rc.local` mounts USB and `/data`, restores MAC from `/etc/MAC`, supports `slp-uninstall` and `slp_run.sh`, and launches `bootstrap.sh`.

Strengths:

- Very small and easy to reason about.
- Uses stock USB update behavior directly.
- Has a simple uninstall marker.
- Already fits Adventurer 3-era `rc.local` boot style.

Weaknesses:

- Replaces the primary stock startup path early and globally.
- Rollback is narrow: restore `rc.local` and power off.
- No staged health-check before taking over future boots.
- No strong model gating or artifact compatibility checks.
- No persistent state migration beyond preserving backed-up `rc.local`.
- No formal distinction between installer, boot shim, runtime, and recovery state.

### X-Forge / `vendor/ff5m`

Current approach:

- Uses init script insertion under `/etc/init.d/` instead of replacing one monolithic startup file.
- Maintains explicit boot flags such as `SKIP_MOD`, `REMOVE_MOD`, `FIRMWARE_SCRIPT`, and detects competing installers.
- Brings up runtime conditionally, including network, kernel modules, `devpts`, `configfs`, and `debugfs` when needed.
- Uses a chrooted runtime and keeps explicit control/marker files under persistent storage.

Strengths:

- Better boot policy separation.
- Much better operator escape hatches.
- Compatible-installer detection is explicit.
- Runtime prerequisites are initialized deliberately.
- More robust logging and state markers.

Weaknesses relative to Adventurer 3 target:

- Heavier than necessary for a first ADV3 stock-like installer.
- Assumes a richer init layout and more mod-owned services.
- Some logic is 5M-specific, including Wi-Fi module handling and display switching assumptions.
- The full chroot/service overlay model may be unnecessary on smaller ADV3 targets.

### 5M Klipper Mod / `vendor/flashforge_ad5m_klipper_mod`

Current approach:

- Installs from USB with a guarded `flashforge_init.sh`.
- Verifies machine, PID, architecture, stock version, and free space.
- Installs to `/data/.klipper_mod`, not directly over stock root runtime assets.
- Uses `/etc/init.d/S00klipper_mod` as a boot interception point.
- Supports skip/remove/debug markers.
- Prepares a chroot with bind mounts for `/dev`, `/proc`, `/sys`, `/run`, `/tmp`, `/data`, and original root mounted read-only.
- Keeps important mutable data persistent via symlinked static storage.
- Explicitly preserves stock MCU state or restores stock MCU firmware when required.

Strengths:

- Best safety model of the three.
- Clean split between stock host system and mod runtime.
- Strong gating before install.
- Controlled persistence model.
- Strong recovery and coexistence behavior.

Weaknesses relative to Adventurer 3 target:

- Built around an ARMv7 + chroot + richer Linux host environment.
- Likely too large as-is for older ADV3-class host constraints.
- Assumes a serviceable `/etc/init.d` interception model rather than the simpler ADV3 `rc.local` path.

## Kernel, Init, And Runtime Implications

### What the evidence suggests about ADV3-like targets

- Existing `sl-print` integration for older FlashForge printers is based on `/etc/rc.local`, not `S00klipper_mod`-style init scripts.
- The replacement `rc.local` mounts `/dev/mmcblk0p2` as `/data`, which implies a simpler older storage layout than the 5M mod's `/dev/mmcblk0p7` handling.
- The runtime expects framebuffer access at `/dev/fb0` and NIC identity restoration from `/etc/MAC`.
- USB execution appears to happen before normal stock app startup and can install by copying a script plus tarball.
- Current known ADV3 install history still involves external-host or helper-driven MCU flashing, especially for N32 targets.

### What should carry forward

From `sl-print`:

- The lightweight USB entrypoint model.
- The use of a small first-stage boot shim.
- Support for one-off USB actions and uninstall markers.

From X-Forge:

- Explicit boot markers and mode selection.
- Competing-installer detection.
- Better logging and operator-visible recovery states.
- Runtime prerequisite checks before taking over the boot path.

From 5M Klipper Mod:

- Install gating by machine/profile/version/free space.
- Staged ownership under a dedicated directory instead of scattering mutable runtime files.
- Read-only view of stock assets where practical.
- Persistent state migration with named files and explicit schema.
- OOBE as a first-class post-install phase rather than a README-only step.

## Proposed Adventurer 3 Safestrap Architecture

### High-level design

Use a two-stage model:

1. USB installer stage.
2. Minimal persistent boot shim stage.

The installer should not immediately replace stock behavior with a fully custom runtime. It should instead install a dedicated `safestrap` directory and switch stock startup to a very small dispatcher that can still fall through to stock.

Suggested layout:

```text
/opt/safestrap/
  manifest.json
  boot/dispatcher.sh
  boot/stock-fallback.sh
  boot/runtime-check.sh
  runtime/
  capture/
  logs/
  assets/
  state/
    install-state.json
    health-ok
    oobe-required
    boot-once-stock
    remove-safestrap
```

### Install phase

USB `flashforge_init.sh` should:

- identify model family and refuse unknown targets by default,
- identify MCU family where possible and refuse unsupported flash paths,
- verify required free space,
- capture stock startup files before changes,
- copy `safestrap` payload to a dedicated directory,
- capture stock calibration/config files before any runtime takeover,
- install a minimal dispatcher into the stock startup hook,
- leave a boot-state record that requires first successful runtime health-check before full takeover.

For ADV3-class systems, this dispatcher will likely still be `/etc/rc.local`-based because that is the clearest in-repo evidence for the older platform.

### Boot dispatcher behavior

The replacement startup hook should be tiny and policy-driven:

1. Mount USB if present.
2. Mount `/data` if available.
3. Restore networking identity from `/etc/MAC` if needed.
4. Check markers in this order:
5. stock firmware update media present -> boot stock updater path.
6. `remove-safestrap` -> restore stock startup and power off.
7. `boot-once-stock` -> clear marker and continue stock boot.
8. runtime health failed too many times -> continue stock boot.
9. `oobe-required` and runtime available -> launch safestrap UI/runtime.
10. healthy runtime -> launch safestrap UI/runtime.
11. otherwise fall through to stock boot.

This is the biggest improvement over `sl-print`: stock remains reachable by policy, not only by lucky recovery.

### Runtime ownership

For Adventurer 3 phase 1, avoid a full 5M-style chroot unless real hardware proves it is needed.

Prefer this order of complexity:

1. minimal stock-host runtime using shipped binaries/scripts,
2. partial bind-mount sandbox only if required,
3. full chroot only if library/runtime conflicts force it.

Reason:

- ADV3 evidence points to an older and simpler boot environment.
- The immediate requirement is safe install, calibration capture, and OOBE, not a general-purpose host distro overlay.
- Chroot adds maintenance cost, storage pressure, and more failure modes.

### Library and runtime policy

- Ship runtime dependencies inside `safestrap` where possible.
- Avoid replacing stock shared libraries in place.
- Prefer self-contained binaries or private library paths under `/opt/safestrap`.
- Treat any need to overwrite system libraries as a phase-gate requiring dedicated validation.

### MCU flashing policy

- Do not promise one-path flashing across all ADV3 variants until N32 support is validated.
- Installer should separate host install success from MCU flash success.
- If a safe on-device flash path exists for a detected board, run it only after preflight validation.
- Otherwise install the safestrap environment and OOBE first, then expose a guided flash or assisted path.

This matches current evidence that host takeover and MCU-family support maturity are not the same problem.

## OOBE Proposal

### Goals

- confirm imported stock calibration data,
- avoid bed damage during first-run checks,
- produce a known-good initial Klipper config state,
- keep all destructive or heat-requiring actions explicit.

### First-boot OOBE phases

#### Phase 1: identity and recovery

- Show detected printer model, board family, MCU family, and chosen config profile.
- Show whether stock calibration data was found.
- Offer `Continue`, `Boot stock once`, and `Uninstall safestrap`.

#### Phase 2: imported calibration review

- Display imported fields and confidence level for each.
- Distinguish raw captured stock values from translated Klipper values.
- Mark uncertain mappings as `captured only` until validated.

#### Phase 3: cold Z safety check

Policy:

- heaters disabled,
- nozzle and bed cold,
- very low-speed movement,
- only move to a conservative above-bed offset first,
- no automatic downward motion beyond a small bounded range without user confirmation.

Suggested behavior:

- home using the chosen safe sequence,
- move to center,
- apply imported Z offset only partially at first,
- stop with nozzle a few tenths of a millimeter above expected contact plane,
- ask the user whether clearance looks reasonable,
- allow only small step jogs from touchscreen before confirming.

The first test should validate sign and rough magnitude, not perfect first-layer behavior.

#### Phase 4: seeded mesh confirmation

- If stock point calibration was captured, translate it into a provisional Klipper mesh profile.
- Mark it as `seeded_from_stock`, not `verified`.
- Probe a reduced set of validation points first.
- If deviation exceeds tolerance, discard provisional mesh and fall back to guided fresh calibration.

#### Phase 5: thermal calibration policy gate

- Do not auto-import thermistor compensation into active heater config without schema confirmation.
- If temperature-related stock values were captured, store them in metadata and show that they are not yet applied automatically.

#### Phase 6: completion

- Save verified values into `printer_data` or equivalent persistent config area.
- Mark OOBE complete.
- Keep the raw stock capture for support and rollback analysis.

## Stock Calibration Import Policy

### General policy terms

Each stock calibration entry should be classified as one of:

- `import-and-apply`
- `import-and-confirm`
- `capture-only`
- `ignore`

And each should also have a confidence level:

- `confirmed`
- `likely`
- `experimental`
- `unknown`

### Group A: Safe geometry and offset entries

These are the highest-value phase-1 targets.

Entries:

- global Z offset,
- home offset / build plate offset,
- probe/nozzle relationship if explicitly stored,
- per-point bed calibration offsets or mesh samples.

Policy:

- Global Z offset: `import-and-confirm`
- Mesh point offsets: `import-and-confirm`
- Probe/nozzle relationship: `import-and-confirm` only when source meaning is clear

Reason:

- These materially improve OOBE.
- They can be validated cold and slowly.
- Wrong values are dangerous, but the failure mode is observable with conservative movement bounds.

### Group B: Thermal compensation entries

Entries:

- `extruderCalibrationValue`
- `nozzleTempDifferentValue`
- any bed temperature correction or sensor calibration values discovered later.

Policy:

- default to `capture-only`
- promote to `import-and-confirm` only after schema meaning and unit/sign semantics are proven on hardware.

Reason:

- A wrong thermal import can cause silent overheating or underheating.
- This is not acceptable for unattended automatic migration.

### Group C: Motion and dimensional compensation

Entries:

- axis skew or geometry trim,
- flow/extrusion calibration,
- dimensional compensation,
- pressure/advance-like values if any exist in newer stock firmware.

Policy:

- default to `capture-only` or `ignore`
- do not apply automatically in phase 1.

Reason:

- Meanings often do not map 1:1 to Klipper concepts.
- These are lower-value than safe first motion and bed contact validation.

### Group D: Device and UI calibration

Entries:

- touchscreen calibration,
- Wi-Fi/network identity,
- display preferences,
- camera settings,
- language or region.

Policy:

- network identity: `import-and-apply` if already handled safely by existing platform logic,
- touchscreen calibration: `capture-only` unless the runtime reuses the stock stack,
- UI preferences: `ignore` for phase 1 unless directly beneficial.

Reason:

- Helpful, but not on the critical path for safe printing.

## Development Target Groups For Documentation

To keep implementation docs readable, split by stability and platform similarity.

### Target Group 1: ADV3 SZ16 baseline printers

Includes:

- Adventurer 3
- Adventurer 3 Lite
- Adventurer 3C
- likely close rebrands when board evidence matches

Documentation focus:

- stock boot path and USB updater behavior,
- startup hook format,
- `/data` mount behavior,
- stock calibration file locations,
- supported MCU families and flash paths.

### Target Group 2: ADV3 Pro deltas

Includes:

- Adventurer 3 Pro

Documentation focus:

- driver differences such as TMC2209 handling,
- any mesh/Z semantics that differ from baseline,
- config deltas that must be surfaced in OOBE.

### Target Group 3: ADV3 N32 path

Includes:

- all in-scope ADV3-family boards using Nation N32 parts

Documentation focus:

- flasher protocol support,
- image build target,
- boot-entry wiring assumptions,
- fail-closed policy when N32 identification is uncertain.

### Target Group 4: ADV4-family adjacent research

Includes:

- Adventurer 3 Pro 2
- Adventurer 4
- Adventurer 4 Lite
- Adventurer 4 Pro

Documentation focus:

- stock config schema evidence such as `/opt/finder_rush/config/config-AD4.json`,
- named stock calibration fields such as `buildZOffset`, `pointN_offset`, `extruderCalibrationValue`,
- translation experiments that may later inform broader policy.

This group should not define ADV3 installer behavior by default, but it is valuable as the clearest calibration-schema evidence currently in-repo.

### Target Group 5: ADV5M reference patterns

Includes:

- Adventurer 5M
- Adventurer 5M Pro

Documentation focus:

- installer gating patterns,
- marker-file policy,
- persistent data layout,
- safe boot interception design.

This group is primarily a design reference for `safestrap`, not a direct runtime template for ADV3.

## Per-Entry Documentation Template

Every stock calibration entry documented in the project should use the same short schema.

Suggested fields:

- Stock name
- Platform group(s)
- Source file or source path
- Observed type and units
- Meaning hypothesis
- Klipper target mapping
- Import policy
- Confidence
- Validation method
- Failure risk if wrong
- Notes

Example:

```text
Stock name: buildZOffset
Platform group(s): adv4-family
Source: /opt/finder_rush/config/config-AD4.json
Observed type and units: unknown numeric, likely mm
Meaning hypothesis: global nozzle-to-bed Z calibration offset
Klipper target mapping: probe.z_offset or saved OOBE z offset
Import policy: import-and-confirm
Confidence: likely
Validation method: cold center-bed clearance test with bounded jog UI
Failure risk if wrong: nozzle crash or false-high first layer
Notes: do not auto-apply without sign validation
```

## Recommended Phase Order

1. Implement `safestrap` as a tiny dispatcher plus dedicated persistent state directory.
2. Add explicit marker-file and stock-fallback behavior before adding more runtime features.
3. Capture raw stock calibration data and write it to a stable schema in persistent storage.
4. Build touchscreen OOBE around `import-and-confirm` geometry entries only.
5. Add seeded mesh validation.
6. Leave thermal entries at `capture-only` until confirmed.
7. Reassess whether a chroot or broader runtime overlay is actually needed for ADV3.

## Bottom Line

For Adventurer 3, the best `safestrap` approach is:

- keep the old-platform-friendly `rc.local` style first-stage entry from `sl-print`,
- add the boot markers, fallback rules, and operator control philosophy from X-Forge,
- add the strong install gating, persistent-state discipline, and OOBE mindset from 5M Klipper Mod,
- avoid the full 5M chroot/runtime complexity unless ADV3 hardware proves it is necessary.

That gives a safer stock-like installer without overfitting Adventurer 3 to the newer 5M architecture.
