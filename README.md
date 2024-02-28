<p align="center">
    <img src="https://raw.githubusercontent.com/VioSynthax/Adventurer-Voxel-Klipper/2.0-preview/images/klippventurer.svg" alt="Klippventurer logo" height="185">
    <h1 align="center">
      Klipper for FlashForge Adventurer<br>
</p>

### Welcome to the Klippventurer project!
#### This is an unofficial port of Klipper for FlashForge printers without official Klipper support.
#### Currently supported machines include Adventurer 3, MonoPrice Voxel, and Arçelik PT1000. See the known issues section at the bottom for exceptions to this.
```TODO: link to "how to find mainboard model"```
### You will need:  

+ very basic soldering skills
+ insulated wire, multiple colors preferred
+ a soldering iron with temperature control
+ solder (preferably SnPb or SnBi solder, SnAg is harder to work with)
+ A Raspberry Pi Zero 2 W with SD card (at least 4GB, 16GB or larger recommended)
+ some wire (silicone preferred for soldered connections as the insulation won't melt)
+ 4x M2.5 screws
+ 4x M2.5 threaded heat-set inserts
+ Raspberry Pi Zero mount from the [STLs](/STLs/) folder
#### Recommended: 
+ Extruder tensioning spacer from the [STLs](/STLs/) folder
+ paste or liquid flux
+ cotton swabs and or cotton pads (something to clean the board with)
+ alcohol or PCB cleaner

# Part 1 - Get Connected 🔌


Step 0: Print the Pi mount and (optional) extruder spacer from the [STLs](/STLs/) folder. The Pi mount is designed to be printable on Adventurer 3 without support material. PLA may not work due to the temperature of the electronics enclosure, HTPLA or higher glass transition filament (PETG, ABS, etc.) recommended.
Insert the 4 threaded inserts into the Pi mount using your soldering iron with a conical tip or threaded insert tip.

Optional Step: With the printer unplugged, use a hex 2.5 bit or Allen key to remove the tensioning arm from the extruder assembly behind the filament hatch. Be careful not to drop the spring into the housing of the printer, it's much less work if you don't have to fish the spring out. Install the spacer onto the spring retention peg. Reassemble the extruder. Video demonstration here: https://www.youtube.com/watch?v=EVRpLRKePUk
This will improve the maximum flow rate significantly, and coupled with the increased extruder motor torque in Klipper, massively reduces or eliminates the clicking extruder issue.

Step 1: Unplug your printer. Carefully turn the printer upside down on a soft surface so as to not scratch the top acrylic panel. Remove the 4 screws securing the metal electronics cover. You may choose to remove the board from the printer and place it on your work surface to make soldering easier. If you do, you'll need to remove the glue covering the MHF micro-coax Wi-Fi antenna cable, unplug all connectors, and remove all 4 screws securing the board.

Step 2: Locate the pads in the diagram (Wiring-diagram.png), apply flux, then tin the pads.
![Pin Mod Diagram](/images/Wiring-diagram.png?raw=true)

Step 2: Solder your wires to the pads as indicated. Make sure all wires are long enough to connect to your Pi, about 10cm or 4in. Strip both ends before soldering to make your life easier later (about 2mm on each end) You should connect the following to your printer's mainboard:
+ one wire to TX (white)
+ one wire to RX (green)
+ one wire to MCU RST (orange)
+ one wire to DFU (purple)
+ one wire to VCC (red)
+ one wire to GND (black)
+ one solder blob bridging the RST pad to the GND pad next to it (blue)
    OR use a very short piece of wire to create the bridge if you prefer, or have difficulty bridging the pads with a blob

# Part 2 - Prepare The Pi 🪛

Step 0: Download the latest Raspberry Pi Imager from https://www.raspberrypi.com/software/ Under "Raspberry Pi Device" select "Raspberry Pi Zero 2 W". For Operating System, go to "Raspberry Pi OS (other)" and select "Raspberry Pi OS (32-bit) Lite", choose your SD card, then hit next. Click "Edit Settings" Set your username and password, select your Wireless LAN country, and set your locale. Under "Services" enable SSH and select "Use password authentication" (unless you know what you're doing). Save, then click "Yes". Once done writing the SD, remove it and insert it into the Pi.
    
Step 1: If you removed the mainboard to solder your wires, reinstall it now, but leave out the 2 screws closest to the power socket. If you left the board in, go ahead and remove those 2 screws now. Place your Pi mounting bracket on the mainboard with the two legs aligned with the holes in the mainboard. Insert the two mainboard screws through the bracket and mainboard to secure them both. Using the 4 M2.5 screws, secure the Pi to the mounting bracket.

Step 2: Time to test! Plug the printer back into power, and flip the switch on. The printer **should not** boot up. Screen should remain black, mainboard LEDs should come on, the Pi should power up, and nothing else should happen. If the printer does boot normally into the stock firmware, make sure RST is bridged to GND and try again. Don't work inside the printer with the power cable connected.

# Part 3 - Klipper, Mainsail, Fluidd, Orca Slicer! ⛵💧🐋

Step 1: SSH into your Pi. You'll need to find the IP in your router's config page or app.

Run the following commands: 
```
echo -e "enable_uart=1\ndtoverlay=miniuart-bt\ndtoverlay=disable-bt" | sudo tee -a /boot/config.txt
```
```
sudo apt update && sudo apt install stm32flash git -y
```
```
cd ~ && git clone https://github.com/dw-0/kiauh.git
```
```
./kiauh/kiauh.sh
```
Once at the KIAUH menu, install 1) Klipper, 2) Moonraker, and 4) Fluidd. You can also install 11) Mobileraker for mobile push notifications via the Mobileraker app. 
Exit KIAUH with B then Q and do
```
git clone https://github.com/synthread/Klippventurer.git ~
```
```
cp -r ~/Klippventurer/config/ ~/printer_data/
```
```
wget -O ~/klipper/.config https://raw.githubusercontent.com/VioSynthax/Klippventurer-Installer/main/configs/adventurer3.config
```
```
cd ~/klipper && make menuconfig
```
press Q followed by Y. Now
```
make
```
and if the build completes without errors,
```
sudo reboot
```
Wait for the Pi to boot back up, SSH back in and do
```
sudo stm32flash -w ~/klipper/out/klipper.bin -R -i -18,23,18:-18,-23,18 /dev/ttyAMA0
```
You should now be able to access the Fluidd interface. If you don't get an error, go ahead and unplug the printer and reinstall the bottom cover. You may need to try restarting the firmware.
Set it back right-side-up and power it on. Reconnect to Fluidd, cross your fingers, and hit the home button! 
#### Calibrate your Z offset and mesh bed leveling before printing!!!
Happy printing!

Part 3 is planned to be a simple script soon.
Please open an issue or pull request if you encounter any problems with any part of this guide.
Still a work in progress, but I need outside testers now to get feedback on the install process and print quality.

# Known Issues
- N32G455 MCU doesn't work with current .config, please open an issue if you have one.
- Can't currently support screen, buzzer, USB, filament runout sensor, or camera, as these components are connected to the MediaTek chip.
- Adventurer 3 Pro works, but you need to switch the stepper drive types from TMC2208 to TMC2209.
