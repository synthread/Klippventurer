<p align="center">
    <img src="images/klippventurer.svg" alt="Klippventurer logo" height="185">
    <h1 align="center">
      Klippventurer<br>
</p>
        
[![](https://dcbadge.vercel.app/api/server/ns2pFdhdMW)](https://discord.gg/ns2pFdhdMW)
### Welcome to the Klippventurer project!
#### This is an unofficial port of Klipper for FlashForge printers without official Klipper support.

See the compatibility chart for printer and feature support.
If your printer is a rebrand, its compatibility is likely to be the same as its FlashForge counterpart.

***
|Printer|Printing|Easy Installer|Pin table|Thermistor Calibration|Display|Touch|Filament Runout|Camera|
|---|---|---|---|---|---|---|---|---|
|Adventurer 3*|✅|⏰|✅|⏰|⏰|⏰|⏰|⏰|
|Adventurer 4|⚠️|⏰|⏰|⚠️|⚠️|⚠️|⚠️|⚠️|

✅Completed ⠀⠀ ⏰In Progress ⠀⠀ ⚠️Not yet Working ⠀⠀ ❌Not planned or lacks hardware feature

>**Note:**
    "Adventurer 3" Includes the Adventurer 3C, Lite, and Pro, as well as rebrands such as the Bresser Rex, Arçelik PT1000, MonoPrice Voxel, and likely any other printer based on the SZ16 mainboard.

#### Always calibrate your Z offset and mesh bed leveling after installing Klipper!!!
Happy printing!

This repo, supported features, and guides change often, join our Discord or watch the repo for updates.
Please open an issue or pull request if you encounter any problems with any part of this guide.
Still a work in progress, but I need outside testers now to get feedback on the install process and print quality.

# Known Issues
- Nation N32G MCU doesn't work with current .config, please open an issue if you have one.
- Can't currently support screen, buzzer, USB, filament runout sensor, or camera, as these components are connected to the MediaTek chip.
- Adventurer 3 Pro works, but you need to switch the stepper driver types from TMC2208 to TMC2209.

# Special Thanks To
[@hw-lunemann](https://github.com/hw-lunemann) for fixing UART muxing and tuning input shaper on Adventurer 3
