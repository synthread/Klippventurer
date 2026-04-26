# Klippventurer v0.3 Preview Plan

This document tracks the current preview push for the `next/firmware-package` branch, marketed as the **USB Install Branch**.

## Release intent

`v0.3.x` is the next-generation installer preview line. It is not a replacement for the current v0.2.x manual/soldered install flow yet.

The v0.3 preview targets:

- a stock-like USB install/update experience,
- Adventurer 3 as the first active target,
- Nation N32 support as the next major Adventurer 3 enablement goal,
- migration of currently working machines toward the new install flow,
- evidence-driven support decisions based on firmware and hardware platform groups.

## Current branch status

- Branch: `next/firmware-package`
- Version line: `v0.3.1-preview`
- Product name: USB Install Branch
- Stability: planning and implementation preview, not release-ready

## Scope for early v0.3 work

### In scope

- Specs and support-matrix documentation.
- Firmware compatibility policy and analysis workflow.
- Adventurer 3 N32 support planning.
- Safestrap/runtime installer architecture.
- Stock calibration capture/import planning.
- Shared USB package layout planning.
- Host-side setup/wrapper planning.

### Out of scope for first preview push

- Definitive v0.3.0 release tag.
- Replacing the live v0.2.x manual install guide.
- Shipping a final USB installer image.
- Claiming untested rebrands or firmware versions are supported.

## Version path

- `v0.2.x`: live manual/soldered install project line.
- `v0.3.0`: first documented USB installer branch baseline.
- `v0.3.1-preview`: active preview line after the initial branch-definition pass.
- `v0.3.0`: first usable beta of the USB installer flow once tested enough to tag.
- `v0.4.0`: planned major machine-support expansion, especially tested Adventurer 5M support with tentative 5M Pro support.
- `v1.0.0`: future stable release once the beta installer flow is ready to exit beta.

## Planned patch roadmap

The branch is currently in `v0.3.1-preview`, not in a released `0.3.1` patch yet.

The patch roadmap below is the current planning breakdown for the preview line and early beta hardening.

### `0.3.1`

Documentation and planning maturity pass focused on the Adventurer 3 USB installer direction:

- safestrap architecture and stock calibration import plan,
- clearer supported/unsupported Adventurer 3 model framing,
- Adventurer 3 N32 flash/build roadmap,
- package-layout and installer-behavior documentation needed before implementation claims become broader.

Current status:

- active preview line on `next/firmware-package`
- documentation-heavy and architecture-heavy
- not yet a tagged patch release

### `0.3.2`

Expected first N32 implementation milestone:

- initial `go-flash` support for Nation N32 targets,
- explicit target-family handling instead of STM-only assumptions,
- repeatable Adventurer 3 N32 firmware build target,
- documented flash preflight and recovery workflow.

### `0.3.3`

Expected first end-to-end USB installer prototype milestone for supported non-Pro-V2 Adventurer 3 variants:

- model and MCU detection wired into the installer path,
- stock calibration capture plus first OOBE path,
- fallback and rollback behavior for failed runtime takeover,
- limited real-hardware beta support statement once the prototype is validated enough to document.

These patch entries are planning targets, not release promises. They should be revised when implementation evidence changes.

## Next documentation tasks

- Keep `docs/specs/` as the evidence-first platform/support source of truth.
- Convert research notes into reviewed docs before relying on them for installer behavior.
- Add `.gitignore` entries for local analysis/vendor/dependency caches.
- Add a vendor/dependency manifest policy before committing vendor-derived work.
- Add architecture decision records for safestrap, package layout, and N32 flashing support.
