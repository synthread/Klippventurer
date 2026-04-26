# ADR 0005: Runtime Overlay And Recovery

## Status

Accepted for `v0.3.1-preview` planning.

## Decision

Runtime changes should be read-only by default once booted, with a separate safe writable overlay that is only mounted after runtime health checks succeed.

## Why

- recovery should not fail because a writable runtime overlay corrupted boot state
- first boot must be able to fall back to stock with minimal assumptions
- keeping mutable state separate makes uninstall and debugging easier

## Runtime model

- stock rootfs remains authoritative
- Klippventurer runtime assets are mounted or referenced deliberately
- writable overlay is not mounted on failed or unverified boot paths
- recovery path should work unless stock rootfs itself is remounted and changed

## Practical early stance

- keep most payload on USB for now
- avoid broad eMMC mutation until required
- separate boot-state markers from larger mutable runtime data

## Recovery decision tree

1. Boot dispatcher starts.
2. Check uninstall or stock-once markers.
3. If runtime health is unknown or failing, do not mount writable overlay.
4. Fall back to stock or recovery-safe path.
5. Only mount writable overlay after runtime success criteria are met.
