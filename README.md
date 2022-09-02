# fv1utils
Misc utils for compiling banks and uploading programs to EEPROM.

## Rom2Header
Converts a rom/bank file to a proper C header file

## ArduinoEEPROMWriter
An Arduino program/sketch for writing a ROM file to en EEPROM. 

# Usage
Use `cat file1.bin file2.bin > bank.bin` to compile a bank of
programs.  The max amount of programs in a bank is 8.

Execute the `upload-rom-to-arduino.sh` script for preparing an arduino
program with the bank embedded and uploading it to the device.

# Dependencies
- Arduino CLI
    https://lindevs.com/install-arduino-cli-on-ubuntu
    $ arduino-cli core install arduino:avr
