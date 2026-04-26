# N32 Flashing Implementation Notes

## Purpose

This document consolidates the evidence currently available in the Klippventurer workspace for implementing actual Nation N32 flashing support.

It is intended to answer three questions:

1. What do we already know from prior art and local repo history?
2. What needs to change in Klippventurer and `go-flash` to support Adventurer 3 N32 boards?
3. What uncertainties still require validation on hardware or against protocol documentation?

## Scope

This document covers:

- Adventurer 3-family boards using Nation N32G45x MCUs
- the existing soldered-host / UART flashing path
- the eventual reuse of that path inside the USB installer work

This document does not claim that N32 flashing is already working.

## Evidence Sources In This Workspace

### Current Klippventurer docs

- `README.md`
- `docs/installation.md`
- `docs/preview-v0.3.md`
- `docs/n32-adventurer3-roadmap.md`
- `docs/specs/safestrap-calibration-plan.md`

### Local repo history

- `277960b` `Fixed some mistakes in flashing STM32`
- `a4fbf24` `N32G455 testing`

### Vendored `go-flash`

- `deps/go-flash/flash/mcu.go`
- `deps/go-flash/flash/serial.go`
- `deps/go-flash/flash/stm32.go`
- `deps/go-flash/flash/stm32_cmd.go`
- `deps/go-flash/flash/flash.go`

### Vendored Klipper

- upstream commit `23e82d37f` `stm32: Add support for Nation N32G45x mcus (#6116)`
- upstream commit `b81784856` `stm32: enable 64KiB bootloader for n32g45x, clarify Makefile output`
- `vendor/klipper/src/stm32/Kconfig`
- `vendor/klipper/src/stm32/Makefile`
- `vendor/klipper/docs/Bootloaders.md`
- `vendor/klipper/config/printer-voxelab-aquila-2021.cfg`

## What We Know Today

### 1. The current manual install path excludes N32

`docs/installation.md` explicitly says Nation N32 MCUs are not supported by the current guide and require a slightly different programming method.

That means the existing soldered-host workflow already established a split between:

- STM32/HK32/GD32 devices using the current `stm32flash`-style path
- Nation N32 devices needing a different implementation

### 2. Earlier project history already expected N32 validation work

Local commit `a4fbf24` updated the older README with the note:

- `Needs testers with the N32G455 MCU`

This is useful because it confirms N32 support was viewed as a board-validation problem early, not just a docs omission.

### 3. Earlier flashing guidance assumed STM32 ROM bootloader behavior

Local commit `277960b` described STM32 flashing as relying on the built-in ROM bootloader over serial rather than a custom bootloader installed on flash.

That assumption still shows up in current workflows:

- wire TX/RX + boot-entry pads to a Pi
- enter the ROM bootloader by board-level strap control
- flash `klipper.bin` over serial

This matters because the present `go-flash` implementation hard-codes those assumptions.

### 4. Upstream Klipper already supports N32G452/N32G455 as MCU targets

Upstream Klipper commit `23e82d37f` added Nation N32G45x support.

Important details from current vendored Klipper:

- `vendor/klipper/src/stm32/Kconfig` exposes:
  - `Nation N32G452`
  - `Nation N32G455`
- `MACH_N32G45x` reuses parts of the STM32F1 family model in Kconfig.
- `MCU` defaults to `stm32f103xe` for `MACH_N32G45x`.
- `CLOCK_FREQ` for `MACH_N32G45x` is:
  - `64000000` with internal reference
  - `128000000` otherwise
- `FLASH_SIZE` for `MACH_N32G45x` is `0x20000`.
- `RAM_SIZE` for `MACH_N32G45x` is `0x10000`.
- `vendor/klipper/src/stm32/Makefile` uses:
  - `-mcpu=cortex-m4`
  - `lib/n32g45x/include`
  - `lib/n32g45x/n32g45x_adc.c`

Interpretation:

- The firmware build target is not the main blocker.
- The transport and flashing path are the more likely blockers.

### 5. Upstream Klipper has explicit N32 bootloader-offset support

Upstream commit `b81784856` later enabled a `64KiB bootloader` option for `MACH_N32G45x`.

Current `vendor/klipper/src/stm32/Kconfig` bootloader menu shows:

- `8KiB bootloader`
- `28KiB bootloader` for `MACH_STM32F1`
- `32KiB bootloader`
- `36KiB bootloader`
- `48KiB bootloader`
- `64KiB bootloader` for `MACH_STM32F103 || MACH_STM32F4 || MACH_N32G45x`

This tells us:

- Upstream Klipper expects some N32 boards to need non-zero application offsets.
- We should not assume `0x08000000` is always correct for Adventurer 3 N32 boards.

### 6. Another upstream config references N32 + 28KB boot + PA9/PA10

`vendor/klipper/config/printer-voxelab-aquila-2021.cfg` contains:

- `Nation N32G452 for N32 version, 28KB boot, serial PA9/PA10`

Interpretation:

- At least one real-world N32 board is expected to use:
  - an N32G452 target
  - a 28KB bootloader/application offset
  - serial on `PA9/PA10`

This is not Adventurer 3-specific evidence, but it is strong prior art that:

- the N32 path is expected to be close to STM32 serial flashing,
- but offset and target selection matter,
- and board-specific bootloader layout must be treated as data, not hard-coded globally.

### 7. Current `go-flash` is structurally reusable but functionally STM-only

The current `deps/go-flash` code is organized as a reusable microcontroller flasher, but it is effectively STM32-only.

Important current behavior:

- `flash/serial.go`
  - always opens UART with `EvenParity`
  - always calls `mc.stmInit()` in `Open()`
- `flash/flash.go`
  - always calls `mc.stmFlash()` in `FlashPayload()`
- `flash/mcu.go`
  - `Identify()` always calls `stmCmdGetId()` and prefixes `STM_`
  - `Reset()` always calls `exitSTBL()`
- `flash/stm32.go`
  - hard-codes STM sync byte `0x7f`
  - hard-codes ACK `0x79` and NACK `0x1f`
  - assumes the current boot-entry GPIO sequence
- `flash/stm32_cmd.go`
  - assumes STM command map for `GET`, `GET ID`, `ERASE`, `WRITE MEMORY`, etc.

Interpretation:

- `go-flash` can likely host N32 support, but it currently has no family abstraction.
- The first implementation step is architectural separation, not only protocol tweaks.

## Most Likely Implementation Shape

## A. Treat N32 as a first-class target family in `go-flash`

Do not try to hide N32 support behind STM-only naming.

Recommended model:

- introduce a target-family concept such as:
  - `stm32`
  - `n32g45x`
- thread that through:
  - initialization
  - identify
  - flash
  - reset/exit behavior
  - command-code selection

Phase 1 should prefer explicit family selection over auto-detection.

Reason:

- current code enters bootloader and immediately speaks STM framing
- if N32 differs even slightly, optimistic auto-detection could turn a recoverable mismatch into a confusing flash failure

### B. Separate boot-entry GPIO control from transport protocol

Right now the code conflates:

- board-level boot-entry pin sequencing
- STM32 bootloader protocol behavior

That makes N32 support harder than it needs to be.

Recommended split:

- boot-entry strategy:
  - controls power / reset / boot pins
  - board-specific
- transport/protocol strategy:
  - sync byte
  - ACK/NACK behavior
  - command map
  - erase/write flow
  - family-specific

This lets us reuse Adventurer 3 board control even if the N32 serial protocol diverges from STM assumptions.

### C. Keep the serial defaults unless evidence disproves them

The strongest current prior art suggests starting with STM-compatible serial assumptions:

- UART serial flashing
- even parity
- 8 data bits
- one stop bit
- sync-first command flow

Reason:

- Klipper upstream treats N32G45x as mostly STM32F103-compatible for build/runtime purposes
- upstream configs and docs imply familiar serial bootloader usage patterns
- current repo docs say N32 uses a slightly different method, not a fundamentally different physical transport

However, these assumptions must be validated against:

- observed bootloader traffic
- official Nation documentation if available
- hardware tests on Adventurer 3 N32 boards

### D. Make application offset an input, not a hard-coded constant

The flashing implementation must accept an explicit application start address.

Why:

- upstream Klipper supports multiple bootloader offsets for N32G45x
- other N32 boards already reference a 28KB bootloader layout
- later upstream support added 64KiB bootloader option for N32G45x

Practical implication:

- the Adventurer 3 N32 flashing path should not assume one universal offset
- the installer or board-profile layer should select:
  - MCU family
  - firmware image
  - application address

## Unknowns That Must Be Resolved

### 1. Exact Adventurer 3 N32 part number(s)

We still need hard evidence for which exact Nation parts are present on supported Adventurer 3 non-Pro-V2 boards.

Known candidate family:

- `N32G45x`

Still needed:

- board photos
- chip markings
- mapping of model -> MCU variant

### 2. Exact Adventurer 3 bootloader/application offset

This is the highest-risk build/runtime unknown after transport compatibility.

Possible outcomes:

- no offset / direct application at `0x08000000`
- 28KB offset
- 64KB offset
- model-specific offsets

Evidence currently supports the need to investigate this, not the final answer.

### 3. Whether N32 ROM bootloader command framing is byte-for-byte STM compatible

We do not yet have proof for:

- sync byte compatibility
- ACK/NACK bytes
- `GET` response format
- chip ID format
- mass erase command behavior
- write unprotect behavior
- reset-after-unprotect behavior

The implementation should therefore be staged to validate these one at a time.

### 4. Exact Adventurer 3 boot-entry control semantics for N32 boards

Current `go-flash` assumes the existing power/boot pin sequence is correct.

That may still be true for N32 boards, but it needs confirmation.

Specific things to validate:

- whether the same board pads are used on N32 variants
- whether boot-entry requires the same pin polarity
- whether reset/exit is the same after flashing
- whether the current soldered-host wiring guide needs N32-specific notes

## Recommended Implementation Plan

### Phase 1: Refactor `go-flash` for family separation

Deliverable:

- no functional N32 promise yet
- clean internal split that stops hard-coding STM behavior everywhere

Changes:

- add target-family config to `Microcontroller`
- route `Open()`, `Identify()`, `FlashPayload()`, and `Reset()` through family-specific handlers
- keep existing STM path behavior unchanged

Exit criteria:

- STM path still works
- N32 scaffolding can be added without copy-pasting the whole library

### Phase 2: Build an N32 probe path

Deliverable:

- an identify/probe step that can confirm whether the candidate N32 protocol works on hardware

Suggested steps:

1. enter bootloader using the current board-control sequence
2. attempt sync using STM-like serial settings
3. try `GET`
4. try `GET ID`
5. log raw responses before assuming semantic success

Exit criteria:

- we can answer whether STM framing is directly reusable or not

### Phase 3: Add N32 erase/write flow

Deliverable:

- a working flash path once probe compatibility is proven

Suggested safeguards:

- require explicit target family
- require explicit application address
- verify expected chip ID before erase
- fail closed on unknown response formats
- add a dry-run / identify-only mode for installer preflight

Exit criteria:

- repeated successful flashes on real hardware

### Phase 4: Lock the Adventurer 3 N32 build recipe

Deliverable:

- one known-good Klipper N32 firmware build recipe per supported board profile

Need to record:

- MCU target (`N32G452` vs `N32G455`)
- bootloader offset
- serial pins
- any Adventurer 3 Pro-specific deltas

Exit criteria:

- build artifact and flash target can be selected deterministically

### Phase 5: Add recovery and installer integration

Deliverable:

- a safe N32 flashing workflow suitable for later USB installer integration

Required behavior:

- identify-only preflight
- strong model/board gating
- explicit failure messages
- documented recovery path when a board remains in bootloader mode or fails to start Klipper

## Concrete Code-Level Recommendations

### `deps/go-flash/flash/mcu.go`

Needs:

- a target-family field in config/state
- `Identify()` returning real family-aware identity instead of always `STM_<id>`

### `deps/go-flash/flash/serial.go`

Needs:

- family-aware init path instead of unconditional `stmInit()`
- likely configurable serial mode if N32 differs after validation

### `deps/go-flash/flash/flash.go`

Needs:

- family-aware flash dispatch instead of unconditional `stmFlash()`

### `deps/go-flash/flash/stm32.go` and `stm32_cmd.go`

Needs:

- shared helpers extracted where protocol overlap exists
- family-specific command map/behavior where N32 differs

### Higher-level installer or wrapper layer

Needs:

- board-profile selection
- explicit application address selection
- preflight checks before erase/write
- separation between host install success and MCU flash success

## Safe Defaults For First Implementation

Unless contradicted by new evidence, the safest initial implementation assumptions are:

1. Start with explicit `n32g45x` family selection, not auto-detection.
2. Reuse the current board-level UART + boot-entry wiring path.
3. Reuse serial parity settings from the STM path for the first probe attempt.
4. Do not mass-erase until probe + chip-id checks succeed.
5. Treat bootloader offset as required board-profile data.
6. Keep installer integration out of scope until manual repeatable flashing works.

## Evidence Gaps To Close Next

1. Capture actual N32 Adventurer 3 chip markings and board revisions.
2. Determine the correct Klipper build target and flash offset for at least one real Adventurer 3 N32 board.
3. Validate whether STM sync / `GET` / `GET ID` works against the target bootloader.
4. Record raw successful and failing bootloader traffic.
5. Decide whether the board-control sequence in current hardware wiring is unchanged for N32 variants.

## Bottom Line

The evidence across this workspace supports a conservative conclusion:

- Upstream Klipper already knows how to build for N32G45x.
- Klippventurer already recognized N32 as a real target family needing special handling.
- The main missing work is not MCU-family compilation support; it is a correct, explicit, testable N32 flashing implementation in `go-flash` plus board-profile data for Adventurer 3.

The safest path is:

1. refactor `go-flash` to stop assuming every target is STM32,
2. probe N32 compatibility with the existing serial bootloader assumptions,
3. lock board-specific application offsets and target IDs,
4. only then wire the flow into the USB installer work.
