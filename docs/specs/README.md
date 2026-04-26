# Specs Documentation

This folder is the planning/source-of-truth area for Klippventurer hardware compatibility.

## Organization model

- Group printers by **hardware/spec platform first**, not by marketing name.
- Marketing names (FlashForge, Monoprice, Bresser, etc.) are aliases used for installer UX and user-facing support messaging.
- A single platform group can still require different `printer.cfg` deltas.
- Engineering docs in this folder should stay platform-centric; vendor/manufacturer naming belongs primarily in [support-matrix.md](./support-matrix.md).

## Support status language

Use only:

- `supported`
- `experimental`
- `untested-likely`
- `unsupported`
- `unknown`

## Confidence language (evidence quality)

Use these confidence labels in notes/evidence columns:

- `confirmed`
- `likely`
- `untested`
- `unknown`

Status and confidence are separate. Example: a model may be `experimental` status with `likely` confidence.

## Evidence-first rule

Every support claim should cite evidence where available, such as:

- local repo docs/paths
- firmware package metadata/checksums
- board markings/teardown notes
- user reports with sufficient technical detail

If evidence is weak or conflicting, avoid overclaiming and prefer `unknown` or `untested-likely` with explicit warnings.

## Docs in this folder

- [platform-groups.md](./platform-groups.md): engineering-facing platform abstractions
- [support-matrix.md](./support-matrix.md): user-facing model/vendor matrix
- [firmware-policy.md](./firmware-policy.md): compatibility policy and installer gating rules
- [firmware-analysis.md](./firmware-analysis.md): repeatable firmware-package analysis workflow

Related project docs:

- [../../README.md](../../README.md)
- [../installation.md](../installation.md)

## Architecture decision records

- [adr-0001-safestrap-first-stage.md](./adr-0001-safestrap-first-stage.md)
- [adr-0002-package-layout-and-rootfs-reuse.md](./adr-0002-package-layout-and-rootfs-reuse.md)
- [adr-0003-kernel-strategy.md](./adr-0003-kernel-strategy.md)
- [adr-0004-static-analysis-decision-tree.md](./adr-0004-static-analysis-decision-tree.md)
- [adr-0005-runtime-overlay-and-recovery.md](./adr-0005-runtime-overlay-and-recovery.md)
