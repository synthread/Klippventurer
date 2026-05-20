# Support Matrix (User-Facing)

This matrix is organized by marketing model/vendor names for user guidance, while mapping each entry to an engineering platform group.

See also:

- Platform abstraction: [platform-groups.md](./platform-groups.md)
- Policy for firmware gating: [firmware-policy.md](./firmware-policy.md)

> [!WARNING]
> `untested-likely` means there is some compatibility evidence, but no guarantee. Installer flows should require explicit user confirmation when evidence is weak.

## Status values

- `supported`
- `experimental`
- `untested-likely`
- `unsupported`
- `unknown`

## Matrix

| Brand / model / SKU | Platform group | Status | Config profile attempted | Known warnings | Evidence / confidence |
|---|---|---|---|---|---|
| FlashForge Adventurer 3 | `adv3-sz16-family` | `supported` | `adv3-sz16-baseline` | Manual path works; USB installer path is still preview work. | Repo compatibility shows ADV3 working ([../../README.md](../../README.md)); confidence: `confirmed` |
| FlashForge Adventurer 3 Pro | `adv3-sz16-family` | `supported` | `adv3-sz16-pro` (TMC2209 delta) | May require TMC2209-specific config handling; N32 board variants still need explicit flashing support. | Known issues note TMC2209 adjustment ([../../README.md#known-issues](../../README.md#known-issues)); confidence: `confirmed` |
| FlashForge Adventurer 3C | `adv3-sz16-family` | `untested-likely` | `adv3-sz16-baseline` | Validate board/MCU variant before assuming parity with ADV3. | Included in ADV3 family note ([../../README.md](../../README.md)); confidence: `likely` |
| FlashForge Adventurer 3 Lite | `adv3-sz16-family` | `untested-likely` | `adv3-sz16-baseline` | Validate bed/nozzle defaults and peripheral differences. | Included in ADV3 family note; package analysis indicates MIPS + older monolithic updater style for classic/Lite stock lineage; confidence: `confirmed` |
| Monoprice Voxel | `adv3-sz16-family` | `untested-likely` | `adv3-sz16-baseline` | Rebrand differences may affect assets/peripherals. | Included in ADV3 rebrand note ([../../README.md](../../README.md)); confidence: `likely` |
| Bresser Rex / Rex WiFi | `adv3-sz16-family` | `untested-likely` | `adv3-sz16-baseline` | Confirm firmware package model checks before install. | Included in ADV3 rebrand note ([../../README.md](../../README.md)); confidence: `likely` |
| Arçelik PT1000 | `adv3-sz16-family` | `untested-likely` | `adv3-sz16-baseline` | Confirm board markings and updater model checks first. | Included in ADV3 rebrand note ([../../README.md](../../README.md)); confidence: `likely` |
| Robo E3 | `adv3-sz16-family` | `unknown` | `adv3-sz16-baseline` (candidate only) | No maintained in-repo evidence yet; do not assume compatibility. | Mentioned in platform planning only; confidence: `unknown` |
| Sharebot One | `adv3-sz16-family` | `unknown` | `adv3-sz16-baseline` (candidate only) | No maintained in-repo evidence yet; do not assume compatibility. | Mentioned in platform planning only; confidence: `unknown` |
| FlashForge Adventurer 3 Pro 2 | `adv4-family` (adjacent) | `experimental` | `adv3-like-devices-config` on likely ADV4-like architecture | Treat as its own category in UX; package internals are clearer, but runtime/config compatibility still needs validation before status promotion. | Package analysis shows ARMv7 + ADV4-like modular packages with `MACHINE=Adventurer3-pro2`, `PID=001D`; confidence: `confirmed` |
| FlashForge Adventurer 4 | `adv4-family` | `unsupported` | `adv4-bridge-candidate` | Currently not yet working in repo; fail closed by default. | Repo compatibility marks ADV4 not yet working ([../../README.md](../../README.md)); confidence: `confirmed` |
| FlashForge Adventurer 4 Lite | `adv4-family` | `unknown` | `adv4-bridge-candidate` | Separate vendor firmware track suggests meaningful deltas from base ADV4. | Package analysis shows ARMv7, kernel `3.4.39+`, `MACHINE=Adventurer4-Lite`, `PID=0016`; confidence: `confirmed` |
| FlashForge Adventurer 4 Pro | `adv4-family` | `unknown` | `adv4-bridge-candidate` | Separate vendor firmware track suggests meaningful deltas from base ADV4. | Package analysis shows ARMv7, kernel `3.4.39+`, `MACHINE=Adventurer4-Pro`, `PID=001E`; confidence: `confirmed` |
| FlashForge Adventurer 5M | `adv5m-t113-family` | `unsupported` | `ad5m-bridge` | Not yet working in repo; keep fail-closed until polish, macro review, and end-to-end validation are complete. | Package analysis shows ARMv7 `MACHINE=Adventurer5M`, `PID=0023`, kernel-family `5.6.0-svn539`, and `NationsCommand`; architecture evidence only, not current repo support; confidence: `confirmed` |
| FlashForge Adventurer 5M Pro | `adv5m-t113-family` | `unsupported` | `ad5m-pro-bridge` | Not yet working in repo; Pro-specific peripherals and macros still require full validation. | Package analysis places Pro in the same ADV5M/T113 family; architecture evidence only, not current repo support; confidence: `confirmed` |
| Afinia QD330 (likely rebrand) | `adv5m-t113-family` | `untested-likely` | `ad5m-pro-bridge` (candidate) | Rebrand package naming and model checks may differ. | Planning hypothesis only pending firmware comparison; confidence: `likely` |

## Installer behavior expectations

- Installer should display detected/selected platform group and chosen config profile before applying changes.
- `unsupported` and `unknown` should fail closed by default.
- Advanced override can exist, but must require explicit user acknowledgement of risks.
- Where board/MCU evidence suggests probable compatibility but no direct test exists, keep status at `untested-likely`.

## Developer reading notes

- Treat this file as the user-facing naming layer.
- For implementation grouping, always cross-check [platform-groups.md](./platform-groups.md).
- For status promotion, follow [firmware-policy.md](./firmware-policy.md).
- For MCU-family questions, especially N32, cross-check [n32-flashing-implementation.md](./n32-flashing-implementation.md).

## Current priority notes

- N32 MCU support is a primary compatibility goal, especially for ADV3 Pro 2 / ADV4 / ADV5-related support paths.
- Rebrands stay conservative in status until hardware and firmware evidence is captured in [firmware-analysis.md](./firmware-analysis.md).
