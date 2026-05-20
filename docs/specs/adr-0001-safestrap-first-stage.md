# ADR 0001: Safestrap First Stage

## Status

Accepted for `v0.3.1-preview` planning.

## Decision

Use a tiny first-stage `safestrap` boot dispatcher instead of taking over the whole stock runtime immediately.

## Why

- Adventurer 3 evidence points to a simple older boot path.
- Recovery must stay easy and predictable.
- We want install, runtime, and recovery to be separate concerns.

## Result

- installer copies a dedicated payload into its own directory
- stock startup hook becomes a tiny dispatcher
- dispatcher can fall through to stock boot when health checks fail or removal is requested

## Non-goals for now

- full distro replacement
- broad service overlay model
- heavy chroot-by-default runtime
