# ADR 0002: Package Layout And RootFS Reuse

## Status

Accepted for `v0.3.1-preview` planning.

## Decision

Keep the USB installer payload as minimal as possible, reuse the underlying stock root filesystem where it is safe, and overlay only what Klippventurer cannot safely borrow from stock.

## Why

- smaller payloads are easier to reason about and recover from
- ADV3-class systems are likely storage- and runtime-constrained
- shipping less reduces compatibility risk and maintenance burden

## Result

- prefer self-contained runtime assets under a dedicated Klippventurer directory
- prefer read-only use of stock assets where practical
- avoid replacing stock libraries in place
- keep as much as possible on USB during early preview work

## Open deployment stance

Do not decide broad eMMC deployment yet. First prove what must live on eMMC, what may remain USB-resident, and what must be persistent across boots.
