# Adventurer-Voxel-Klipper
#### This repository aims to provide a way to flash Klipper to a MonoPrice Voxel or Adventurer 3 3D printer. This can be done by disabling and bypassing the MediaTek SBC with a Raspberry Pi.

While I am using Fluidd here instead of OctoPrint, if using a suitably powerful host OctoPrint is an option. There is further info at the end of the guide if you choose to attempt this with the original processor. You are solely responsible for any damage caused by following any part of this guide.

# Part 1 - Modify Your Board

## You will need:

+ very basic soldering skills
+ insulated wire, preferably silicone
+ a soldering iron with temperature control
+ solder (preferably SnPb or SnBi solder, SnAg is harder to work with)

+ recommended - paste or liquid flux
+ optional - cotton swabs and or cotton pads (something to clean the board with)
+ optional - alcohol or PCB cleaner

Step 0: Remove the board from the printer and place it on your work surface.

Step 1: Apply flux (if using) and then tin the pads shown in the diagram. If you will be using tweezers for the DFU pads, don't tin those ones.

Step 2: locate the pads in the diagram (ADV-VOX-Pin-Mod.png) and solder wires bridging the pads as indicated. You should now have the following:
+ one wire bridging the two pads with the green jumper in the diagram
+ one wire bridging the two with the white jumper indicator
+ one wire bridging the RST points indicated in the diagram
+ If you are using tweezers, leave off the DFU jumper indicated in magenta. If not, solder a jumper

![Pin Mod Diagram](https://github.com/VioSynthax/Adventurer-Voxel-Klipper/raw/main/ADV-VOX-Pin-Mod.png)

Congratulations, if you've successfully followed the instructions up to this point, you have a 3D printer board ready to flash with the lastest version of Klipper (or the firmware of your choosing) Follow to Part 2 for flashing Klipper and making your bridge cable

# Part 2 - Installing Klipper and establishing a link to your host
Note - FLASH YOUR HOST BOARD BEFORE DOING THIS. I recommend Fluidd Pi. https://github.com/fluidd-core/FluiddPI

## You will need: 

+ A flashed host board (probably a Raspberry Pi Model 2, 3, or 4, but anything you can run Klipper on will work)
+ some wire
+ female Dupont pins, 2.54mm pitch, and 2x JST XH 4-pin female connectors
+ 2x single pin female Dupont header connectors for powering the Pi
+ a crimping tool for said pins, here's a cheap one: https://www.amazon.com/gp/product/B088NQV8Z3
+ wire strippers or something with which to strip your wires
+ a 24 volt input capable buck converter set to 5.1 volts output, here's a good option: https://www.amazon.com/LM2596-Converter-4-0-40V-1-25-37V-Voltmeter%EF%BC%882pcs%EF%BC%89/dp/B085WC5G8N/
+ an ST-Link V2 programmer: https://www.amazon.com/gp/product/B07SQV6VLZ
+ STM32CubeProgrammer installed on your computer (or an alternative if you prefer)

Step 0: Locate the BOOT0 DFU pins on your mainboard as indicated in the diagram (ADV-VOX-Pin-Mod.png) solder a wire across these two pins or carefully bridge the two with tweezers when the guide says to

Step 1: Connect your ST-Link programmer to the ST-Link port on your mainboard (labeled V C D G on the silk screen in the top right corner)

Step 2: Open STM32CubeProgrammer "CubeProg" and select at the top right ST-LINK. You may need to update the ST-LINK if it's brand new. Port - SWD, Mode - Hot Plug. Reset Mode - Software. If using tweezers, tweeze now. Hit Connect. Remove tweezers.

Step 3a: Select the "Open file" tab at the top of CubeProg. Select your bootloader, STM32Duino.bin and under the Download dropdown, make sure Address is set to 0x08000000 then click Download.

Step 3b: Select the "Open file" tab once more. Select your Klipper.bin file. (I have provided one, but you should really compile your own on the Raspberry Pi or other host or you might get a version mismatch.) Select the Download dropdown and set the Address to 0x8002000 this time. Hit download again, and if everything is successful you can remove your BOOT0 DFU jumper (if you soldered one) and reinstall the mainboard in your printer. Don't connect everything yet, you should connect power and test the serial connection with your host first. If successful, reconnect everything.

Step 4: Measure 3 wires of about 4-6 inches in length, strip both ends of each just enough for the end to be exposed inside the Dupont connector (about 1mm). Crimp a pin onto both ends of each wire.
Insert these 3 wires into the GND, RX, and TX pins on JST connectors as shown in the diagram. Pay close attention to serial TX on the Pi being connected to RX on the main board. TX stands for "transmit" and RX for "receive" so these need to be crossed TX>RX and RX>TX.

Note: Don't connect the 5v line from the serial header to the Pi, it should not be connected.

## If you want a neat and tidy internal power connection to the Pi, complete Step 5.

Step 5: Remove the bottom plastic panel of the printer to expose the screw terminals on the power supply.
Unscrew the positive and negative 24 volt output screws.
Prepare two wires for powering your buck converter of about 4 inches. Strip one end about 2mm and the other 1/4 inch.
Solder the shorter ends to your buck converter's input pads or insert them into its screw terminals, depending on which it has.
Twist the ends of 1/4 inch of stripped wire to keep them together. Insert under your PSU's output terminals, on top of the spade connectors. Tighten town your terminal screws and tug on the wires and spade connectors to ensure a tight fit. Do NOT tin the wires, this can cause a poor connection and result in arc flash. If you want an even more secure connection, splice these wires into the mainboard supply wires instead.
Solder another pair of wires about 6 inches in length or more to the output on your buck converter, or insert them into the output screw terminals if your converter has them. Make sure your buck converter is not outputting more than 5.2 volts at this time and no less than 5.1 volts. Pi 1/2 likely only need 5.0 volts, Pi 3/4 may need 5.1 volts. Stay clear of the exposed mains voltage when you test this.
Strip 1mm off the ends of these output wires, Crimp a single pin Dupont female pin onto each wire's end. Insert each into a female Dupont header connector.
Connect these wires to the 5V and GND connections on your Pi's GPIO header that aren't taken up by our serial cable we made.

Step 6: INSULATE YOUR PI AND BUCK CONVERTER. You can 3D print cases for these to mount them on the top side of the access plate to sort of hinge open if you'd like, make sure you've cut wires long enough for however you are mounting your host board. If not, cover them in some tape and make sure nothing shorts when you close your access panel. You might even be able to squeeze the Pi somewhere near where the screen of the printer is, you'll need to remove the side to access that spot. I haven't tried, so not sure if there's room

Optional Step: Add an antenna connector to your Pi if it doesn't already have one. Attach the printer's built in WiFi antenna to your Pi.
More info here- https://hackaday.io/project/10091-raspberry-pi-3-external-antenna


# Part 3 - The Easy Part

Step 0: Just rename Adventurer-3.cfg to printer.cfg and place it inside Fluidd Pi's config directory, or your home directory if using OctoPi + Klipper.

Step 1: Live your hobbyist 3D printing dreams and tune your printer forever instead of actually printing with it. Or print with it. Your choice.

Step 2: ???

Step 3: Profit.
