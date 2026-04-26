# Klippventurer v0.3 Preview Plan

This document tracks the first preview push for the `next/firmware-package` branch, marketed as the **USB Install Branch**.

## Release intent

`v0.3.0` is the next-generation installer preview line. It is not a replacement for the current v0.2.x manual/soldered install flow yet.

The v0.3 preview targets:

- a stock-like USB install/update experience,
- Adventurer 3 as the first active target,
- Nation N32 support as the next major Adventurer 3 enablement goal,
- migration of currently working machines toward the new install flow,
- evidence-driven support decisions based on firmware and hardware platform groups.

## Current branch status

- Branch: `next/firmware-package`
- Version line: `v0.3.0-preview`
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
- `v0.3.0-preview`: USB installer branch preview with Adventurer 3/N32 as the first new target.
- `v0.3.0`: first usable beta of the USB installer flow once tested enough to tag.
- `v0.4.0`: planned major machine-support expansion, especially tested Adventurer 5M support with tentative 5M Pro support.
- `v1.0.0`: future stable release once the beta installer flow is ready to exit beta.

## Next documentation tasks

- Keep `docs/specs/` as the evidence-first platform/support source of truth.
- Convert research notes into reviewed docs before relying on them for installer behavior.
- Add `.gitignore` entries for local analysis/vendor/dependency caches.
- Add a vendor/dependency manifest policy before committing vendor-derived work.
- Add architecture decision records for safestrap, package layout, and N32 flashing support.
