# Adventurer-Voxel-Klipper
This repository aims to provide a way to flash Klipper to a MonoPrice Voxel or Adventurer 3 3D printer. This is achieved by replacing the MediaTek SBC with a Raspberry Pi SBC.

It is possible that the onboard processor can run Klipper, but it is significantly less powerful than even the original Raspberry Pi, and will not currently be covered under this guide.
While I am using Fluidd here instead of OctoPrint, if using a suitably powerful host OctoPrint is an option. You are solely responsible for any damage caused by following any part of this guide.

Part 1 - Modifying the Mainboard

You will need:

basic to intermediate soldering skills
flux
wire
alcohol (isopropanol, ethanol, methanol) or PCB cleaner
cotton swabs and or cotton pads (something to clean the board with)
solder wick braid
a soldering iron with temperature control
a heat gun or heat wand
A tool to lift the SBC while desoldering it (tweezers, knife, thin screwdriver, etc.)
aluminum tape (if you're worried about damaging surrounding components)
LEADED solder (highly recommended)
CHIPQUIK or some bismuth solder (optional, but is extremely helpful if you want to make your life easier or don't have great soldering skills or a heat wand)

Step 0: Remove the board from the printer and place it on a nice heat resistant surface.

Step 1:  Apply your tape if you're using it around the edges of the SBC at this time. Begin by heating your soldering iron and applying your bismuth solder or CHIPQUIK to the two rows of pads on either side of the main processor sub-board.
Make sure to add your low-temp solder to all 62 pads on the edges of the SBC and ensure the original solder melts and mixes with it. Theoretically you could remove this with only a soldering iron, but a heat gun or wand is easier and less likely to damage your pads.

Step 2: Use your heat gun on the low setting or heat wand on high (400-500C) with a large nozzle to pre-heat the SBC and surrounding area from about 6 inches away for about one minute (this is better for the board to prevent heat shock). Move your tool closer to the SBC moving it in a circular pattern slowly accross the pins, moving between sides every few seconds. Within a minute or so you should be able to bump the edge of the SBC gently and it should move a bit. Use your tool to gently lift the SBC off the board and place it somewhere heat resistant. You may now turn off your heat gun/wand and allow the board to cool.

Step 3: Apply flux to your pads and clean them off using your soldering iron and your solder wick. Clean flux residue off of your board with your alcohol or cleaner.

Step 4: Apply fresh solder to the pads in the diagram (ADV-VOX-Pin-Mod.png) and solder wires bridging the pads indicated. Congratulations, if you've successfully followed the instructions up to this point, you have a 3D printer board ready to flash with the lastest version of Klipper (or the firmware of your choosing) Follow to Part 2 for flashing Klipper and making your bridge cable

Part 2 - Installing Klipper and establishing a link to your host
Note - FLASH YOUR HOST BOARD BEFORE DOING THIS. I recommend Fluidd Pi. https://github.com/fluidd-core/FluiddPI

You will need: 

A flashed host board (probably a Raspberry Pi Model 2, 3, or 4) (If you go the route of attempting to use an unmodified board, open an issue if you need assistance, I'll update the guide if you're successful)
some wire
female Dupont pins, 2.54mm pitch, and 2 JST XH 4 pin female connectors
Optional - 2 single pin female Dupont header connectors for powering the Pi internally
a crimping tool for said pins
wire strippers will be helpful but you could use scissors, pliers, nail clippers, a knife, a lighter, etc.
a 24 volt input capable buck converter set to 5.1 volts output
an ST-Link V2 programmer (optional, can be helpful)
STM32CubeProgrammer installed on your computer (or an alternative if you prefer)

Step 0: Locate the BOOT0 and 3.3v pins on your mainboard as indicated in the diagram (ADV-VOX-Pin-Mod.png) solder a wire across these two pins or carefully bridge the two with tweezers when the guide says to

Step 1: Connect your ST-Link programmer to the ST-Link port on your mainboard (labeled V C D G on the silk screen in the top right corner)
Alternate Step 1: Connect your serial adapter to the serial port on your mainboard (labeled VCC RX TX GND on the silk screen in the bottom center)

Step 2: Open STM32CubeProgrammer "CubeProg".
    If ST-LINK: Select at the top right ST-LINK. You may need to update the ST-LINK if it's brand new. Port - SWD, Mode - Hot Plug. Reset Mode - Software. If using tweezers, tweeze here. Hit Connect.
    If Serial: Select at the top right UART. Find your serial port name e.g. "COM1" and set Baudrate to 115200. Tweezer time is now. Hit Connect.

Step 3a: Select the "Open file" tab at the top of CubeProg. Select your bootloader, STM32Duino.bin and under the Download dropdown, make sure Address is set to 0x08000000 then click Download.
Step 3b: Select the "Open file" tab once more. Select your Klipper.bin file. (I have provided one, but you should really compile your own on the Raspberry Pi or other host or you'll probably get a version mismatch.) Select the Download dropdown and set the Address to 0x8002000 this time. Hit download again, and if everything is successful you can remove your BOOT0 DFU jumper and reinstall the mainboard in your printer.

Step 4: Measure 4 wires of about 4-6 inches in length, strip both ends of each just enough for the end to be exposed inside the Dupont connector (about 1mm). Crimp a pin onto both ends of each wire.
Insert these wires into the JST connectors as shown in the diagram. Pay close attention to serial TX on the Pi being connected to RX on the main board. TX stands for "transmit" and RX for "receive" so these need to be crossed TX>RX and RX>TX.

If you want a neat and tidy internal power connection to the Pi, complete Step 5.

Step 5: Remove the bottom plastic panel of the printer to expose the screw terminals on the power supply.
Unscrew the positive and negative 24 volt output screws.
Prepare two wires for powering your buck converter of about 4 inches. Strip one end about 2mm and the other 1/4 inch.
Solder the shorter ends to your buck converter's input pads.
Twist the ends of 1/4 inch of stripped wire to keep them together. Insert under your PSU's output terminals, on top of the spade connectors. Tighten town your terminal screws and tug on the wires and spade connectors to ensure a tight fit. Do NOT tin the wires, this can cause a poor connection and result in arc flash. If you want an even more secure connection, splice these wires into the mainboard supply wires instead.
Solder another pair of wires about 6 inches in length or more to the output on your buck converter. Make sure your buck converter is not outputting more than 5.2 volts at this time and no less than 5.05 volts. Stay clear of the exposed mains voltage when you test this.
Strip 1mm off the ends of these output wires, Crimp a single pin Dupont female pin onto each wire's end. Insert each into a female Dupont header connector.
Connect these wires to the 5V input and Ground connections on your Pi's GPIO header that aren't taken up by our serial cable we made.

Step 6: INSULATE YOUR PI AND BUCK CONVERTER. You can 3D print cases for these to mount them on the top side of the access plate to sort of hinge open if you'd like, make sure you've cut wires long enough for however you are mounting your host board. If not, cover them in some tape and make sure nothing shorts when you close your access panel.

Optional Step: Add an antenna connector to your Pi if it doesn't already have one. Attach the Adventurer or Voxel's built in WiFi antenna to your Pi.
More info here- https://hackaday.io/project/10091-raspberry-pi-3-external-antenna


Part 3 - The Easy Part

Step 0: Just rename Adventurer-3.cfg to printer.cfg and place it inside Fluidd Pi's config directory, or your home directory if using OctoPi + Klipper.

Step 1: Live your hobbyist 3D printing dreams and tune your printer forever instead of actually printing with it. Or print with it. Your choice.

Step 2: ???

Step 3: Profit.




Potential future additions:

Replacing the MediaTek SBC with a Raspberry Pi Compute Module 4 interface board to make this entire process way easier and super quick. (quite possibly)

Make an update file that removes the "Finder" firmware and replaces it with Klipper, sets up Fluidd, and flashes the MCU all in one go, that's also reversible with a regular FlashForge update USB stick. (probably won't do this but who knows, it is definitely possible and really not that much further work, the thing is already running Linux [OpenWrt] See: https://github.com/ihrapsa/KlipperWrt) Completely bypasses the need for soldering. It can also be installed manually using something like CrashForge which for some reason seems to only work on some people's ADV 3s (USB stick formatting problem? Different board revision?)