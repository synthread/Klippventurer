# Uninstall (soldered)

Reverting your printer to stock firmware is easy. Simply remove your wires from the mainboard, remove the Pi mount, and power on the printer. After a short wait on the splash screen, you should be greeted by a prompt asking if you want to retry the "update". This will reinstall the factory firmware on the microcontroller. Follow the on-screen prompts and your printer should be back to stock in no time. If this is unsuccessful, don't worry, and continue on to the troubleshooting section below.

## Troubleshooting
Sometimes, this built-in recovery mechanism can fail to reinstall the original firmware. In this case, you will need to reinstall the firmware manually. Thankfully, this is easy to do and only requires a USB flash drive. Steps are as follows: 

0) Download the latest firmware for your printer from the manufacturer's website. For the FlashForge branded Adventurer 3, it can be found [here.](https://www.flashforge.com/download-center/77) For rebrands, you can likely find your firmware on the seller's website by Googling "\<printer name> firmware" or you can try using the official firmware from FlashForge. Some models such as the MonoPrice Voxel can unlock additional features by using official firmware instead of firmware from the printer's distributor.

> [!CAUTION]
> Do not use a flash drive containing any important files, you will need to erase this drive.

1) Insert your USB flash drive into your PC.
2) Flash your drive to FAT32:
    - (Windows) Use [GUIFormat](http://ridgecrop.co.uk/index.htm?guiformat.htm) to format your flash drive as FAT32
    - (Linux) Use your tool of choice to format your flash drive as FAT32
    - (macOS) Open Disk Utility. Select your flash drive on the left, then "Erase". Choose "MS-DOS (FAT32)" from the dropdown, then click the "Erase" button.
3) Extract the firmware .zip file. Open the folder within to find your firmware name e.g. Adventurer3-*.tgz Move this file onto your empty flash drive, and eject it from the computer.
4) With the printer powered off, insert your USB flash drive. Power on the printer.
5) You should see a firmware update begin after the splash screen. Once this is completed, you may remove the flash drive and restart the printer. You should now be back to factory firmware!