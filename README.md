# fv1utils
Misc utils for compiling banks and uploading programs to EEPROM (24LCxx)

## Rom2Header
Converts a rom/bank file to a proper C header file

## ArduinoEEPROMWriter
An Arduino program/sketch for writing a ROM file to en EEPROM. 

# Usage
Use `cat file1.bin file2.bin > bank.bin` to compile a bank of
programs.  The max amount of programs in a bank is 8.

Execute the `upload-rom-to-arduino.sh` script for preparing an arduino
program with the bank embedded and uploading it to the device.

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
