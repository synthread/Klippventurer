<div align="center">
  <img src="images/klippventurer-3.svg" alt="Klippventurer logo" height="185">
  <h1>Klippventurer</h1>
  <h3>Upgrade your FlashForge with Klipper.</h3>
  <a href="https://discord.gg/ns2pFdhdMW">
    <img src="https://dcbadge.limes.pink/api/server/ns2pFdhdMW" alt="Discord Server">
    
  </a>
</div>

<div align="center">
  <a href="#start-here">Start Here</a> •
  <a href="#compatibility">Compatibility</a> •
  <a href="#docs">Docs</a> •
  <a href="#known-issues">Known Issues</a> •
  <a href="#special-thanks">Special Thanks</a>
</div>

> [!IMPORTANT]
> This branch is the **USB Install Branch** for `v0.3.1-preview` planning and implementation.
> The current live manual/soldered install flow remains on `main` / `v0.2.x`.
> See [docs/preview-v0.3.md](docs/preview-v0.3.md) for preview scope and version planning.

Klippventurer is a work-in-progress effort to make Klipper installs on FlashForge hardware simpler, safer, and more repeatable.

## Start Here

- Want the current manual install flow? Use `main` / `v0.2.x` and read [docs/installation.md](docs/installation.md).
- Want to follow the next-generation installer work? Stay on `next/firmware-package` and read [docs/preview-v0.3.md](docs/preview-v0.3.md).
- Want the deeper support and platform docs? Start in [docs/specs/](docs/specs/README.md).

## Compatibility

| Printer family | Current state | Notes |
|---|---|---|
| Adventurer 3 family* | `working` manual path, `preview` USB installer path | Main focus of `v0.3.x`, including N32 enablement work. |
| Adventurer 3 Pro 2 | `experimental` | Treat separately until board and installer behavior are better pinned down. |
| Adventurer 4 family | `not yet working` | Active research, not ready for install claims. |
| Adventurer 5M / 5M Pro | `working` in current repo context | Separate platform family from ADV3 installer work. |
| Creator Pro 2 / Creator 3 family | `not yet working` | Not part of the current USB installer focus. |

For the detailed user-facing matrix, see [docs/specs/support-matrix.md](docs/specs/support-matrix.md).

> [!NOTE]
> "Adventurer 3" Includes the Adventurer 3C, Lite, and Pro, as well as rebrands such as the Bresser Rex, Arçelik PT1000, MonoPrice Voxel, and likely any other printer based on the SZ16 mainboard.

## Docs

- [Preview plan](docs/preview-v0.3.md)
- [Installation guide](docs/installation.md)
- [Specs and support docs](docs/specs/README.md)

> [!WARNING]
> Always calibrate your Z offset and mesh bed leveling after installing Klipper!

This repo, supported features, and guides change often, join our [Discord](https://discord.gg/ns2pFdhdMW) or watch the repo for updates.
Please open an issue or pull request if you encounter any problems with installation.

## Known Issues 
### Adventurer 3 Models
- Nation N32G MCU support is still under active flashing/build-target work.
- Screen, buzzer, USB, filament runout, and camera support are part of the longer-term easy-installer direction, not the current manual path.
- Adventurer 3 Pro works, but you need to modify printer.cfg to use TMC2209 drivers instead of TMC2208.

## Special Thanks
[@hw-lunemann](https://github.com/hw-lunemann) for fixing UART muxing and tuning input shaper on Adventurer 3

[@kyleisah](https://github.com/kyleisah) and everyone who contributed to [KAMP](https://github.com/kyleisah/Klipper-Adaptive-Meshing-Purging)

[@KevinOConnor](https://github.com/KevinOConnor) and everyone who contributed to [Klipper](https://github.com/Klipper3d/klipper)

[@FlashforgeOfficial](https://github.com/FlashforgeOfficial) for good hardware at a fair price
