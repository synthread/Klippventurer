# Platform Groups (Engineering View)

This document groups machines by hardware/spec platform for implementation work.

> [!IMPORTANT]
> Platform groups are engineering abstractions. Keep vendor/manufacturer marketing names in the user-facing [support-matrix.md](./support-matrix.md).

## Group: `adv3-sz16-family`

### Summary

First target for the shared bridge/shim package.

### Scope

- Legacy Adventurer-3-class SZ16-style hardware and close derivatives.
- Excludes ADV3 Pro 2 as a direct member until stronger evidence is collected.

### Evidence

- Project compatibility note groups ADV3 variants and SZ16-based rebrands: [../../README.md](../../README.md)
- Current known issues and config assumptions (TMC2208/TMC2209, N32 caveat): [../../README.md#known-issues](../../README.md#known-issues)
- Recent package analysis indicates ADV3 classic/Lite stock updater lineage is MIPS-based and uses an older monolithic package style.

### Known deltas

- Some units appear to require TMC2209 assumptions instead of TMC2208 defaults.
- Nation N32 MCU variants are not yet broadly validated in current config paths.
- Rebrand units may differ in peripherals and defaults (branding, camera, runout, nozzle/bed presets).

### Support stance

- Core profile: `supported`/`experimental` depending on exact board + MCU evidence.
- Rebrands: typically `untested-likely` unless confirmed by board/firmware evidence.

---

## Group: `adv4-family`

### Summary

Likely ARMv7 Linux-generation platform, distinct from ADV3 SZ16 MIPS/monolithic updater lineage.

### Scope

- Adventurer-4-class architecture and close derivatives.
- ADV3 Pro 2 is currently treated as **adjacent** and likely closer to this architecture, with package-level evidence now available.

### Evidence

- Current project compatibility marks ADV4 as not yet working: [../../README.md](../../README.md)
- Vendor firmware tracks for Lite/Pro indicate meaningful differentiation (to verify in firmware analysis).
- Package-analysis findings so far:
  - ADV3 Pro 2: ARMv7, ADV4-like modular packages, `MACHINE=Adventurer3-pro2`, `PID=001D`
  - ADV4 Lite: ARMv7, kernel `3.4.39+`, `MACHINE=Adventurer4-Lite`, `PID=0016`
  - ADV4 Pro: ARMv7, kernel `3.4.39+`, `MACHINE=Adventurer4-Pro`, `PID=001E`

### Known/likely constraints

- RAM constraints may be similar to ADV3 in practice; treat this as a hypothesis and verify before committing feature assumptions.
- ARM/Linux expectations are likely, but exact SoC/kernel/package behavior must be confirmed per firmware set.
- N32 MCU support is a primary compatibility goal for this family and related families.

### Support stance

- Default `experimental` or `unknown` until board + firmware internals are confirmed.

---

## Group: `adv5m-t113-family`

### Summary

Newer Linux platform associated with existing AD5M mod ecosystem.

### Scope

- Adventurer 5M and close derivatives likely sharing T113-class architecture.
- Do not auto-merge with other AD5x models without direct evidence.

### Evidence

- Community/vendor mod ecosystem reports T113-S3 + 128MB-class environment (needs package-level verification per model).
- Project compatibility marks 5M/Pro as working in current repo context: [../../README.md](../../README.md)
- Factory ADV5M package analysis indicates ARMv7 packaging with `MACHINE=Adventurer5M`, `PID=0023`, kernel-family reference `5.6.0-svn539`, and control bundle usage of `NationsCommand`.

### Known deltas

- Pro variants may add enclosure/camera/filtration/noise-control features and package naming differences.
- Reported N32 MCU presence should be treated as a key cross-family compatibility factor.

### Support stance

- Bridge-oriented path is likely practical first; onboard strategy remains separate work.

---

## Cross-family decisions (current)

- **ADV3 Pro 2 handling:** may appear as its own category in user docs; package analysis supports ADV4-like modular ARMv7 architecture, while devices-config behavior still requires conservative validation before broader status promotion.
- **N32 MCU support:** primary compatibility objective for ADV3 Pro 2 / ADV4 / ADV5-related coverage.
- **Implementation split:** decide per finding whether changes belong in Klipper patchset vs `go-flash` transport/installer logic.
- **Reference implementations:** N32 support design should draw from lessons in X-Forge / 5M Klipper Mod work where technically applicable, then be validated against Klippventurer constraints.

## Evidence quality guidance

When uncertain, record assumptions as `likely` or `unknown`, and defer status promotion until verified by firmware package analysis and board-level evidence.
