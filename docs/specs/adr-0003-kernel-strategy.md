# ADR 0003: Kernel Strategy

## Status

Accepted for `v0.3.1-preview` planning.

## Decision

Assume the stock kernel should remain in place unless static analysis or runtime validation proves that Klippventurer requires kernel changes.

## Why

- kernel replacement is one of the riskiest recovery surfaces
- the current goal is a safe installer path, not a custom distro fork
- many needed changes may be solvable in userspace or by shipping private runtime assets

## Decision tree

1. Can the feature run with the stock kernel as-is?
2. If not, can it run with existing modules or configuration already present?
3. If not, can the feature be deferred without blocking the installer milestone?
4. Only then consider kernel patching or recompilation.

## Current default

- leave kernel alone
- document the blocker
- justify any future kernel change with evidence
