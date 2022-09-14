# fv1utils
Misc utils for compiling banks and uploading programs to EEPROM (24LCxx)

## ROM2Header
Converts a ROM/Bank file to a proper C header file.

## ROMSplit
Reads a multiprogram ROM/Bank file and splits it into separate binary
files.

## ROMBuilder
Reads multiple program binaries and compiles them together in to
multi-program ROM/Bank file.

## ArduinoEEPROMWriter
An Arduino program/sketch for writing a binary file to an EEPROM. 

# Usage
Use the `ROMBuilder` utility to compile a bank of programs. The max
amount of programs in a bank is 8. (An alternative is to use `cat
file1.bin file2.bin > bank.bin` command)

Execute the `upload-rom-to-arduino.sh` script for preparing an arduino
program with the bank embedded and uploading it to the device.

Use the following command to monitor the output from the upload process:

  $ arduino-cli monitor -p /dev/ttyUSB0  -c baudrate=115200


## Arduino programmer setup
- Connect the `3v3` pin to `pin8` of the 24LCxx chip. 
  - We can also use the `AREF` pin which is 5V, but this might be too
    much for other ICs which might be connected to the same line if
    the EEPROM is to be programed _in-situ_.
- Connect the `GND` pin to `pin4` of the 24LCxx chip.
- Connect `Analog in 4` pin to `pin5` of the 24LCxx chip (`SDA`).
- Connect `Analog in 5` pin to `pin6` of the 24LCxx chip (`SCL`).

Ensure that `pin7` of the 24LCxx chip is tied to ground or else it
will be in write-protect mode.

# Dependencies
- Arduino CLI
    https://lindevs.com/install-arduino-cli-on-ubuntu
    $ arduino-cli core install arduino:avr
