# N32 Adventurer 3 Roadmap

## Goal

Build a new `go-flash` path for Nation N32-based Adventurer 3 mainboards, then use that as the firmware-delivery foundation for a USB install package covering non-Pro-V2 Adventurer 3 variants.

## Current State

- The current manual install guide explicitly excludes Nation N32 MCUs.
- `deps/go-flash` is STM32-only in practice even though it is already structured as a reusable flasher library.
- The vendored Klipper tree appears to already include at least partial `n32g45x` support, so the highest-risk gap is likely flashing behavior rather than basic MCU-family awareness.
- The repo already points at a future `next/firmware-package` direction for stock-like installation.

## Assumptions To Verify Early

- The affected Adventurer 3 family boards use an N32G45x part with a bootloader interface close enough to STM32 to reuse most packet handling.
- The Adventurer 3 board-level boot-entry wiring for N32 can still be driven through the same or similar GPIO-controlled sequence used by the current soldered-host workflow.
- The correct Klipper build target, boot offset, and UART pins can be standardized per board family without needing per-printer manual guessing.
- Non-Pro-V2 Adventurer 3 variants share enough install behavior that one USB package can branch on detected board/MCU details instead of shipping totally separate installers.

## Workstreams

### 1. Hardware And Variant Inventory

Deliverable: a compatibility matrix for Adventurer 3-family boards relevant to installation.

Need to capture:

- Board revision names and photos.
- MCU vendor and exact part number for each supported printer variant.
- Known differences between Adventurer 3, 3 Lite, 3C, 3 Pro, rebrands, and Pro V2.
- Boot strap, reset, and UART pin behavior as exposed on the board.
- Whether Pro V2 is excluded because of a different board, MCU, display stack, or installer path.

Exit criteria:

- We can state exactly which models are in-scope for the new N32 path.
- We can map each in-scope board to one flashing recipe and one firmware build recipe.

### 2. N32 Flash Protocol Support In `go-flash`

Deliverable: a `go-flash` implementation that can identify, erase, write, and reset the target N32 MCU reliably.

What needs investigation:

- Does N32 accept the same sync byte, parity, ACK/NACK, `GET`, `GET ID`, `ERASE`, and `WRITE MEMORY` framing as STM32?
- Are command codes identical, partially identical, or version-dependent?
- Does mass erase differ?
- Does write protection or readout protection behave differently enough to require special handling?
- Does the chip require a different post-write reset/exit sequence?

Likely code changes:

- Stop treating every target as STM unconditionally in `deps/go-flash/flash`.
- Introduce explicit target-family selection or protocol detection.
- Separate boot-entry GPIO sequencing from protocol framing where they are currently conflated.
- Add an N32 command/backend path that can reuse shared serial helpers.
- Make `Identify()` return the actual family so higher-level installer logic can branch safely.

Exit criteria:

- We can flash an N32 board repeatedly without needing `stm32flash` or manual recovery steps beyond the known hardware entry wiring.
- We have logs or captured transactions proving the implementation is stable across cold boot and retry cases.

### 3. Adventurer 3 N32 Klipper Build Target

Deliverable: a documented and reproducible firmware build for the N32 Adventurer 3 board.

Need to pin down:

- Exact Klipper menuconfig target for the Adventurer 3 N32 board.
- Bootloader offset.
- UART pins used by the installed workflow.
- Any board-specific differences between standard Adventurer 3 and Pro hardware.
- Whether the existing pin table and `printer.cfg` need N32-specific changes or only MCU build-target changes.

Artifacts to produce:

- A committed reference `.config` or generated config source for N32 Adventurer 3 builds.
- Build instructions or scripted build entrypoint.
- Notes on any variant-specific config deltas.

Exit criteria:

- A fresh checkout can build the correct N32 firmware artifact without manual menuconfig exploration.

### 4. Safe Flashing And Recovery Workflow

Deliverable: an operator-safe process for first flash, retry, and recovery.

Need to define:

- Pre-flight checks: detected serial port, target family, expected chip ID, expected image size, expected load address.
- Failure modes: no ACK, wrong chip ID, erase timeout, partial write, bad reboot.
- Recovery steps for a board that remains in bootloader mode or fails to start Klipper.
- Clear rules for when the installer should stop instead of trying another destructive action.

Exit criteria:

- We have a manual validation checklist and a scripted pre-flight sequence.
- A failed flash attempt does not leave the user with ambiguous next steps.

### 5. USB Installer Foundation

Deliverable: a package architecture that can install Klippventurer from USB on supported non-Pro-V2 Adventurer 3 models.

This should build on the N32 flasher work, not run in parallel with guesswork.

Core design questions:

- Where does the installer execute: stock Linux userspace, init hook, recovery environment, or external helper?
- How will the package identify printer model, board revision, and MCU family automatically?
- How will firmware artifacts, configs, and UI assets be packaged and versioned?
- What is the rollback story if package install succeeds but runtime startup fails?
- Which features are in phase 1 of the USB installer: only MCU flash plus host files, or also screen/buzzer/USB integration?

Exit criteria:

- There is a concrete package format, entrypoint, detection flow, and rollback strategy.

## Recommended Phase Order

1. Inventory the supported Adventurer 3 N32 variants and confirm exact MCU/board targets.
2. Prove a manual N32 flash path outside the final installer UX.
3. Refactor `go-flash` so N32 support is explicit and testable.
4. Lock the Adventurer 3 N32 Klipper build config and artifact naming.
5. Validate repeated flash/recovery cycles on real hardware.
6. Build the USB installer on top of the proven flasher and build artifacts.

## Immediate Next Tasks

1. Document every Adventurer 3-family board/MCU variant we currently know about and mark which are non-Pro-V2 targets.
2. Compare `go-flash` assumptions against the actual N32 bootloader protocol documentation or observed traffic.
3. Produce a first known-good N32 Klipper build recipe for the Adventurer 3 board.
4. Decide whether `go-flash` should auto-detect protocol family or require an explicit target in phase 1.

## Open Questions

- Which exact N32 part numbers have been observed on Adventurer 3 non-Pro-V2 boards in hand?
- Do we already have a known-good N32 Klipper `.config` from a test machine, or is that still to be derived?
- Is the future USB installer expected to run entirely from stock printer hardware, or may it rely on temporary helper hardware during beta?
- Are there model-specific bootloader offsets or UART routes within the Adventurer 3 family?
