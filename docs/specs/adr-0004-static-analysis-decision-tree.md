# ADR 0004: Static Analysis Decision Tree

## Status

Accepted for `v0.3.1-preview` planning.

## Decision

Use static analysis first for firmware, updater, kernel, and runtime questions. Promote findings into implementation only after they are documented and mapped to an explicit decision.

## Why

- safer than executing vendor artifacts
- gives repeatable evidence for support claims
- keeps implementation split decisions reviewable

## Decision tree

1. Extract package metadata and updater structure.
2. Record model, MCU, kernel, and rootfs evidence.
3. Decide whether the finding changes:
   - support matrix
   - platform group classification
   - `go-flash` responsibilities
   - Klipper runtime/build assumptions
4. If a finding would change runtime behavior, write or update a spec/ADR before implementation.
5. If evidence is weak, keep the status conservative.

## Current rule

Do not run vendor update payloads or helper binaries directly as part of normal analysis.
