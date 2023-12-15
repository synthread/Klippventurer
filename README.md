<p align="center">
  <a>
    <img src="https://raw.githubusercontent.com/VioSynthax/Adventurer-Voxel-Klipper/2.0-preview/images/klippventurer.svg" alt="Klippventurer logo" height="185">
2.0-preview
    <h1 align="center">Klipper for Adventurer 3</h1>
  </a>
</p>

#### You are solely responsible for any damage or injury caused by following any part of this guide.
### You will need:

+ very basic soldering skills
+ insulated wire, preferably silicone
+ a soldering iron with temperature control
+ solder (preferably SnPb or SnBi solder, SnAg is harder to work with)
+ A Raspberry Pi (a Pi with built-in WiFi is preferred)
+ a 24 volt input capable buck converter set to 5.1 volts output (unless you want to keep the Pi outside the printer or power it seperately)
+ + here's a good option: https://www.amazon.com/LM2596-Converter-4-0-40V-1-25-37V-Voltmeter%EF%BC%882pcs%EF%BC%89/dp/B085WC5G8N/
+ female Dupont pins, 2.54mm pitch (not needed if soldering directly to Pi Zero)
+ 1x Dupont 8-pin connector (Multiple Dupont connectors with fewer pins work too, or you could solder straight to GPIO pads on Pi Zero)
+ a crimping tool (not needed if soldering directly to Pi Zero)
+ some wire
+ + here's a kit with all the connectors, pins, crimping tool, and wire that you'll need
+ + https://www.amazon.com/Qibaok-Crimping-Ratcheting-Crimper-Connectors/dp/B07ZHB4BBY/
#### Recommended: 
+ paste or liquid flux
+ cotton swabs and or cotton pads (something to clean the board with)
+ alcohol or PCB cleaner

# Part 1 - Tapping in 🪛



Step 0: Unplug your 3D printer. Remove the bottom plastic cover from the printer (this is so we can access the power supply's DC output to power our Pi). You may choose to remove the board from the printer and place it on your work surface to make soldering easier.
![Pin Mod Diagram](https://github.com/VioSynthax/Adventurer-Voxel-Klipper/blob/a8485fdaa321842ca7af45ca6d088fc077095493/images/Wiring-diagram.png?raw=true)

Step 1: Locate the pads in the diagram (Wiring-diagram.png), apply flux, then tin the pads. Don't use flux on RST>GND or you'll have difficulty making a blob.

Step 2: Solder your wires and solder blob to the pads as indicated. Make sure all 5 wires here are long enough to connect to your Pi, wherever you've chosen to mount it. Strip both ends before soldering to make your life easier later (about 2mm on each end) You should connect the following to your printer's mainboard:
+ one wire to TX (white)
+ one wire to RX (green)
+ one wire to MCU RST (orange)
+ one wire to DFU (purple)
+ one wire to GND (black)
+ one solder blob bridging the RST pad to the GND pad (blue/black) see magnified view
    OR use a very short piece of wire to create the bridge if you prefer, or have difficulty bridging the pads with a blob

# Part 2 - Connecting to the Pi 🔌

Step 0: Download the latest Raspberry Pi Imager from https://www.raspberrypi.com/software/ choose Raspberry Pi OS (32-bit), then click the gear icon and enable SSH, use password authentication (unless you know what you're doing), set your username and password, configure wireless LAN (enter your WiFi SSID and password, choose your wireless LAN country) then save, choose your SD card or USB drive, and write your image. Once done, remove the boot drive and insert it into the Pi.
    
Step 1: Choose where you would like to mount your buck converter, if you're using one. I recommend using foam tape to mount it right between the large capacitors / power connector and the MediaTek board. This will power the Pi when it's installed into the printer's electronics box. Cut a pair of wires (one red, one black) long enough to reach from your buck converter to the printer's power supply. Strip the ends by about 3-5mm. Loosen the screws of the two unused terminals next to your printer board's power cable. Twist the ends of your wires, and bend them into a hook shape. Wrap the hooked end of your black wire around the "V-" screw in a clockwise direction (insert from the left side of the screw). Do the same for your red wire, this time with the "V+" screw terminal. Route these power wires toward the same area as your printer's mainboard.

Step 2: Insert the V+ and V- wires into your buck converter's "IN+" and "IN-" terminals respectively, and tighten down the screw terminals (or solder them on if your buck converter doesn't have screw terminals)

Step 3: Reinstall the bottom cover of the printer, leaving the metal access panel open. Sit your buck converter on the plastic outer shell of the printer with the unit upside-down. (careful with your printer in this orientation, don't go slamming it around. The bed is only mounted on one linear rod) Plug your AC power cable into the printer with the power switch turned off. Make sure your buck converter isn't sitting on a metal surface, then go ahead and flip the power switch to the "on" position. Turn the small potentiometer dial on the potentiometer until the voltage reads 5.1 volts. If your converter doesn't have an LCD voltage readout, use a multimeter. Once the voltage is set, turn the switch off and unplug the printer.

###### If soldering wires directly to Pi Zero, skip step 4 and solder your GPIO connections now (plus a 5.1v wire to your buck converter's OUT+ terminal)
Step 4: Back over to our mainboard. Prepare your crimp connectors and make sure the ends of your wires soldered to the board are stripped. Crimp a female pin onto each of the 5 wires. Cut one additional wire long enough to reach from near the mainboard's power input plug to the Pi's GPIO header. Strip the ends and crimp a pin onto one of them. Insert all of these wires into an 8-pin Dupont connector EXACTLY in the order shown in the diagram. (Wiring-diagram.png) You should have, in this order: 
##### Empty, red (5.1v), black (GND), white (TX), green (RX), orange (MCU RST), Empty, purple (DFU)
You will have an additional loose wire coming from the Dupont connector.

Step 5: Reinstall the mainboard into the printer, and reconnect all of the cables. Connect the loose 5.1v wire from the Dupont connector to the "OUT+" terminal on the buck converter. Wrap the buck converter in kapton tape. Apply double-sided foam tape to the underside, and stick it on the printer's mainboard. Leave the 8-pin Dupont connector loose for the moment. Time to test! Plug the printer back into power, and flip the switch on. The printer shouldn't boot up. Screen should remain black, LEDs should come on, and nothing should happen. If it does boot normally, make sure RST is bridged to GND and try again. Don't work inside the printer with the power cable connected.

Step 6: Apply foam tape to the bottom of the Pi. Make sure the entire bottom is covered in foam tape or you can use Kapton as well if you don't want to use a ton of foam tape if using a full sized Pi. Don't remove the backing from the tape just yet. With the power cable disconnected, connect the 8-pin Dupont connector to the Pi's GPIO header as shown in the diagram. Test fit the Pi by placing it inside the printer's electronics box, next to the buck converter. Make sure both fit inside the case, between the mainboard and the access cover. You can now remove the backing from the foam tape, and stick your Pi to the inside of the metal cover in roughly the same position as it was in while test fitting. Sit the cover loosely onto the bottom of the printer, and do one more test. The LED on the Pi should light up, and the LEDs on the mainboard should as well. Screw on the bottom cover, and that concludes the hardware portion of this guide! Onto part 3 to setup the firmware and software.


# Part 3 - The Home Stretch? 🍰

# This part of the guide is unfinished. Watch this repo to receive a notification when it's ready.
This will involve a script that configures the UART port, builds Klipper.bin for the STM32 inside the printer, and installs it over serial- all automatically. This is a huge departure from the earlier versions of this guide, where an ST-Link programmer was used to manually flash the microcontroller with a bootloader and Klipper firmware, then the Pi needed manual configuration to make serial work. This was janky, overly complicated, and just outright didn't work for most people. I don't want to call the new method foolproof just yet, but I do think it will massively simplify most of the firmware setup.
=======

