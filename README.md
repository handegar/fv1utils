# fv1utils
Misc utils for compiling banks and uploading programs to EEPROM.

## Rom2Header
Converts a rom/bank file to a proper C header file

## ArduinoEEPROMWriter
An Arduino program/sketch for writing a ROM file to en EEPROM. 

## Scripts
The scripts uses the srecord tools

### RomConcat
Concats several compiled files to a single bank-file (binary)

### Rom2Hex
Converts a rom file to a hex file (script using on the "srecord" tools)


# Dependencies
- Arduino CLI
  https://lindevs.com/install-arduino-cli-on-ubuntu
  $ arduino-cli core install arduino:avr
- SRecord
  Installed via "apt"
