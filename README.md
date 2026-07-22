<div align="center">
  <img src="images/klippventurer-3.svg" alt="Klippventurer logo" height="185">
  <h1>Klippventurer</h1>
  <h3>Klipperize your FlashForge printer!</h3>
  <a href="https://discord.gg/ns2pFdhdMW">
    <img src="https://dcbadge.limes.pink/api/server/ns2pFdhdMW" alt="Discord Server">
    
  </a>
</div>

<div align="center">
  <a href="#compatibility">Compatibility</a> •
  <a href="#installation">Installation</a> •
  <a href="#known-issues">Known Issues</a> •
  <a href="docs/changelog.md">Changelog</a> •
  <a href="#special-thanks">Special Thanks</a>
</div>

> [!IMPORTANT]
> `main` documents the current live Klippventurer setup: a **v0.2.x manual/soldered host install** workflow (ongoing state, not a tagged release line).
> Looking for the next-generation stock-like USB installer work? See the preview **USB Install Branch**: [`next/firmware-package`](https://gitlab.com/synthread/proj/Klippventurer/-/tree/next/firmware-package).
> That branch is preview-only and is not merged to `main` until a verified USB installer exists.


## Compatibility

|Printer|Printing|Easy Installer|Pin table|Printer.cfg|Thermistor Calibration|Display|Touch|Filament Runout|Camera|
|---|---|---|---|---|---|---|---|---|---|
|Adventurer 3/Pro*|✅|⏰|✅|✅|⏰|⏰|⏰|⏰|⏰|
|Adventurer 3 Pro 2|❓|⏰|❓|❓|⏰|⚠️|⚠️|⚠️|⚠️|
|Adventurer 4|⚠️|⏰|⚠️|⚠️|⏰|⚠️|⚠️|⚠️|⚠️|
|Adventurer 5M/Pro|✅|⚠️|✅|✅|✅|✅|✅|✅|✅|
|Creator Pro 2|⚠️|⚠️|❓|⚠️|⚠️|⚠️|⚠️|❌|❌|
|Creator 3/Pro |⚠️|⏰|⚠️|⚠️|⚠️|⚠️|⚠️|❌|❌|

✅Working ⠀⠀ ⏰In Progress ⠀⠀ ⚠️Not yet Working ⠀⠀ ❌Not planned or lacks hardware feature ⠀⠀ ❓Untested, might work

> [!NOTE]
> "Adventurer 3" Includes the Adventurer 3C, Lite, and Pro, as well as rebrands such as the Bresser Rex, Arçelik PT1000, MonoPrice Voxel, and potentially other printers based on the SZ16 mainboard.

## Installation
For installation instructions, please see [the docs](docs/installation.md)

> [!WARNING]
> Always calibrate your Z offset and mesh bed leveling after installing Klipper!

This repo, supported features, and guides change often, join our [Discord](https://discord.gg/ns2pFdhdMW) or watch the repo for updates.
Please open an issue or pull request if you encounter any problems with installation.

## Known Issues 
### Adventurer 3 Models
- Nation N32G MCU doesn't work with current .config
- Can't currently support screen, buzzer, USB, filament runout sensor, or camera. Support for these is planned as part of future USB installer work.
- Adventurer 3 Pro works, but you need to modify printer.cfg to use TMC2209 drivers instead of TMC2208.

## Special Thanks
[@hw-lunemann](https://github.com/hw-lunemann) for fixing UART muxing and tuning input shaper on Adventurer 3

[@kyleisah](https://github.com/kyleisah) and everyone who contributed to [KAMP](https://github.com/kyleisah/Klipper-Adaptive-Meshing-Purging)

[@KevinOConnor](https://github.com/KevinOConnor) and everyone who contributed to [Klipper](https://github.com/Klipper3d/klipper)

[@FlashforgeOfficial](https://github.com/FlashforgeOfficial) for good hardware at a fair price
