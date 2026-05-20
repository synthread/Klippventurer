# Firmware Compatibility Policy

This policy defines how Klippventurer evaluates stock firmware compatibility and installer gating.

## Core principles

1. **Do not blindly prefer latest vendor firmware.**
   - Newer vendor release date does not automatically imply safer/better compatibility for Klippventurer.

2. **Track two layers of versioning.**
   - Vendor-visible version strings (what users see in stock UI/download pages).
   - Klippventurer normalized compatibility class (what actually matters for our tooling/runtime assumptions).

3. **Equivalence is interface-based, not marketing-based.**
   - Firmware builds can be treated as equivalent when kernel/libs/interfaces/components we depend on are unchanged, even if unrelated stock binaries changed.

4. **Evidence-first compatibility promotion.**
   - Promote matrix status only after package-level analysis and reproducible evidence.

## Installer gating rules

- `supported`: allow normal install path.
- `experimental`: allow with explicit warning and confirmation.
- `untested-likely`: allow only with stronger warning + explicit confirmation.
- `unsupported`: block by default; advanced override only.
- `unknown`: block by default; advanced override only.

Unknown or unsupported firmware should stop install by default unless user explicitly chooses advanced override.

## Artifact handling policy

- Synthread/Klippventurer should **not** host proprietary stock firmware blobs by default.
- Maintain references to official vendor sources (e.g., FlashForge CDN URLs) plus checksums and metadata.
- Keep only analysis metadata in git where possible (hashes, manifests, parsed fields, notes).

## Provenance and signatures

- A machine-readable compatibility matrix may be signed by Synthread CI for provenance/integrity.
- Signature/provenance data is an integrity aid, not a replacement for user risk acknowledgement.

## Platform-policy implications (current)

- ADV3 Pro 2 may be user-visible as a separate matrix row, but engineering assumptions should remain conservative: likely ADV4-family architecture with ADV3-like devices-config behavior until verified.
- ADV3 and ADV4 RAM constraints may be similar; treat this as unverified until firmware/system evidence confirms memory budgets.
- N32 MCU support is a primary compatibility goal for ADV3 Pro 2 / ADV4 / ADV5-related coverage; implementation decisions should be split intentionally between:
  - Klipper patchset changes (protocol/MCU runtime support)
  - `go-flash` changes (transport, detection, installer behavior)

## Decision record requirements

When changing support status or policy assumptions, include:

- source package identifiers and checksums
- analyzed model/board/MCU evidence
- impact on platform group classification
- rationale for any change in installer gating behavior

Cross-reference findings in [firmware-analysis.md](./firmware-analysis.md) and user impact in [support-matrix.md](./support-matrix.md).
