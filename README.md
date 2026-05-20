<div align="center">
  <img src="images/klippventurer-3.svg" alt="Klippventurer logo" height="185">
  <h1>Klippventurer</h1>
  <h3>Klipper-powered upgrades for FlashForge Adventurer printers.</h3>
  <p>
    Keep good hardware useful longer with safer firmware packaging, guided setup, and a path away from vendor software that has stopped keeping up.
  </p>
  <a href="https://discord.gg/ns2pFdhdMW">
    <img src="https://dcbadge.limes.pink/api/server/ns2pFdhdMW" alt="Discord Server">
  </a>
</div>

<div align="center">
  <a href="#start-here">Start Here</a> •
  <a href="#compatibility">Compatibility</a> •
  <a href="#how-support-grows">How Support Grows</a> •
  <a href="#docs">Docs</a> •
  <a href="#special-thanks">Special Thanks</a>
</div>

> [!IMPORTANT]
> This branch is the **USB Install Branch** for the `v0.3.x` preview line. It is where Klippventurer is becoming a packaged, stock-like upgrade path for FlashForge Adventurer printers.
>
> The stable **v0.2.x manual/soldered host install** line lives on [`main`](https://gitlab.com/synthread/proj/Klippventurer/-/tree/main).

Klippventurer is an open-source effort to make Klipper installs on FlashForge Adventurer hardware simpler, safer, and more repeatable. The project is shifting toward a USB-based installer, model-aware compatibility checks, stock calibration import, and recovery tooling that feels approachable for everyday printer owners.

## Start Here

| If you are... | Start with... |
|---|---|
| New to Klippventurer | Read this README, then join Discord before flashing anything. |
| Following the upcoming USB installer | [Preview plan](docs/preview-v0.3.md) |
| Checking whether your printer is a good target | [Compatibility](#compatibility), then the detailed [support matrix](docs/specs/support-matrix.md) |
| Bringing up new hardware | Collect evidence first; a community porting guide is planned for `v0.3.3+`. |
| Looking for the older hardware-assisted guide | It is kept as legacy documentation under [docs/installation.md](docs/installation.md). |

The short version: **Adventurer 3-family support is the main focus now**, with Nation N32 support and a safer USB install flow as the next big unlocks.

## Compatibility

This table is a friendly summary, not the source of truth. Installer decisions should come from the deeper specs and compatibility metadata as the `v0.3.x` line matures.

| Printer family | Status | What that means today |
|---|---|---|
| Adventurer 3 family | **Partial** | The legacy hardware-assisted path has known working coverage on supported boards. The USB installer path is the active `v0.3.x` focus. |
| Adventurer 3 rebrands | **Experimental** | Printers such as Monoprice Voxel, Bresser Rex, Arçelik PT1000, Robo E3, and Sharebot One may share the same platform, but need evidence and testing before broad claims. |
| Adventurer 3 Pro 2 | **Experimental** | Appears closer to the Adventurer 4 architecture while still being an Adventurer 3-class product; treat separately until validated. |
| Adventurer 4 family | **Planned** | Active research target. Not ready for normal users yet. |
| Adventurer 5M / 5M Pro | **Planned** | Later validation and polish target after the Adventurer 3 USB installer groundwork is stable. |

Status values:

- **Supported**: expected to work for the documented install path on known-good hardware and firmware.
- **Partial**: important pieces work, but not everything a normal user would expect is packaged or automated yet.
- **Experimental**: plausible and actively investigated, but users should expect testing and debugging.
- **Planned**: intended future support, not ready for users yet.

For detailed status, evidence, and engineering platform groups, see [docs/specs/support-matrix.md](docs/specs/support-matrix.md).

> [!NOTE]
> “Adventurer 3 family” includes the Adventurer 3, Adventurer 3C, Adventurer 3 Lite, and Adventurer 3 Pro. Rebrands need confirmation because firmware branding, board revisions, peripherals, and MCU variants can differ.

## What We Are Building

The `v0.3.x` work is focused on turning Klippventurer into a safer packaged upgrade flow:

- stock-like USB install/update behavior,
- model and firmware compatibility checks before risky actions,
- Nation N32 support for newer Adventurer 3-family boards,
- stock calibration capture before switching runtimes,
- fallback/recovery paths when something is not ready,
- a bridge-first design for constrained printers,
- future support for more Adventurer-family machines as evidence and testers become available.

Offline USB install remains a core requirement. Some older FlashForge networking stacks may not handle modern Wi-Fi setups such as WPA3 or band-steered 2.4+5 GHz networks, so online downloads should be treated as a convenience rather than a requirement until validated.

## How Support Grows

Klippventurer support is evidence-first. A printer is not considered supported just because it looks similar from the outside.

Before a new printer can be promoted, we need evidence such as:

- exact printer model and rebrand information,
- board photos and revision markings,
- MCU vendor and part number,
- stock firmware package URLs and checksums,
- updater script/package layout,
- CPU architecture, kernel, and storage layout,
- calibration/config locations,
- safe flash and recovery behavior,
- real hardware validation from someone willing to test.

A formal porting guide is planned for `v0.3.3+`. Until then, bring research and test results to Discord so we can compare them against the Adventurer platform work already underway.

## Docs

- [Preview plan](docs/preview-v0.3.md)
- [Specs and support docs](docs/specs/README.md)
- [Support matrix](docs/specs/support-matrix.md)
- [Platform groups](docs/specs/platform-groups.md)
- [Firmware policy](docs/specs/firmware-policy.md)
- [Legacy hardware-assisted installation guide](docs/installation.md)

> [!WARNING]
> Klipper conversions can move motors and heat components in ways stock firmware did not. Always validate the target printer, firmware, Z offset, bed mesh, and thermal behavior before trusting a converted machine.

## Known issues (legacy v0.2.x manual install)

These apply to the hardware-assisted path on `main`, not the in-development USB installer:

- Nation N32G MCU is not supported by the legacy manual guide; N32 work is tracked under the USB preview line.
- Screen, buzzer, USB, filament runout sensor, and camera are not supported on the legacy path; USB installer work may revisit peripherals.
- Adventurer 3 Pro may require TMC2209 driver settings in `printer.cfg` instead of TMC2208.

## Special Thanks

[@hw-lunemann](https://github.com/hw-lunemann) for fixing UART muxing and tuning input shaper on Adventurer 3.

[@kyleisah](https://github.com/kyleisah) and everyone who contributed to [KAMP](https://github.com/kyleisah/Klipper-Adaptive-Meshing-Purging).

[@KevinOConnor](https://github.com/KevinOConnor) and everyone who contributed to [Klipper](https://github.com/Klipper3d/klipper).

[@FlashforgeOfficial](https://github.com/FlashforgeOfficial) for good hardware at a fair price.

<div align="center">
  <hr>
  <p>Maintained by the Klippventurer community with Synthread Labs.</p>
  <img src="images/branding/synthread-wordmark.svg" alt="Synthread Labs wordmark" height="36">
</div>
