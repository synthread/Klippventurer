# Firmware Analysis Workflow

This document describes a repeatable, evidence-first process for analyzing factory firmware packages without committing proprietary payloads.

See policy context: [firmware-policy.md](./firmware-policy.md)

## Goals

- Determine whether a firmware package belongs to an existing compatibility class.
- Extract enough metadata to drive support-matrix decisions.
- Avoid unsafe execution and avoid storing proprietary firmware blobs in git.

## Storage rules

- Commit metadata, hashes, manifests, parsed notes, and analysis outputs.
- Do **not** commit proprietary firmware binaries by default.
- Record original source URL + retrieval timestamp for reproducibility.
- `.firmware-analysis/` is local scratch/artifact cache and should not be committed.

## Current observed package traits (analysis-derived)

- ADV3 classic/Lite stock updater lineage: MIPS, older monolithic package style.
- ADV3 Pro 2: ARMv7, ADV4-like modular packages, `MACHINE=Adventurer3-pro2`, `PID=001D`.
- ADV4 Lite: ARMv7, kernel `3.4.39+`, `MACHINE=Adventurer4-Lite`, `PID=0016`.
- ADV4 Pro: ARMv7, kernel `3.4.39+`, `MACHINE=Adventurer4-Pro`, `PID=001E`.
- ADV5M factory package: ARMv7, `MACHINE=Adventurer5M`, `PID=0023`, kernel-family reference `5.6.0-svn539`, control bundle uses `NationsCommand`.

> [!NOTE]
> These findings improve platform classification confidence but do not, by themselves, imply end-to-end install/runtime support status changes.

## Preferred acquisition approach

1. Attempt partial retrieval first (range-fetch / selective ZIP extraction).
2. Current preferred candidate tool: `marcograss/partialzip`.
3. Fall back to full download only when necessary for analysis.

## Safety requirements

- Never execute updater scripts/binaries from extracted packages.
- Perform static extraction/inspection only.
- Treat scripts as data to parse, not commands to run.

## Analysis checklist

Capture the following for each package:

1. **Package identity**
   - model label(s), version string(s), source URL, download date
   - file size and checksum set (at minimum SHA256)

2. **Package structure**
   - archive layout (top-level dirs/files)
   - updater entrypoints and invocation flow (static interpretation)

3. **Model/board gating**
   - explicit model checks
   - board/MCU identifiers and compatibility gates

4. **System/runtime surface**
   - kernel/rootfs identifiers
   - critical library/runtime versions relevant to Klippventurer tooling
   - partition targets and mount/update destinations

5. **Persistence behavior**
   - paths expected to survive updates
   - config/data locations used by vendor update scripts

6. **Integrity/signature mechanics**
   - signature/checksum verification fields and methods
   - any anti-tamper/update constraints that affect install strategy

7. **Compatibility conclusion**
   - equivalent to known class / new class / inconclusive
   - confidence (`confirmed` / `likely` / `untested` / `unknown`)
   - suggested support status change (`supported`, `experimental`, `untested-likely`, `unsupported`, `unknown`)

## Output format guidance

For each analyzed package, produce a concise record containing:

- package metadata + checksums
- extracted structural notes
- model/board/MCU findings
- compatibility decision + confidence
- links to affected docs (usually [platform-groups.md](./platform-groups.md) and [support-matrix.md](./support-matrix.md))

## Current focus areas

- Validate ADV3 Pro 2 as likely ADV4-family architecture with ADV3-like devices-config assumptions, or split if evidence contradicts.
- Verify ADV3 vs ADV4 memory constraints (RAM budget assumptions).
- Prioritize N32 MCU protocol/boot/update behavior characterization across ADV3 Pro 2 / ADV4 / ADV5-related packages.
- Feed implementation split decisions into Klipper patchset vs `go-flash` responsibilities.
