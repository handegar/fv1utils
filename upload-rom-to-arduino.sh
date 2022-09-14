#!/bin/bash -e

if [ "$#" != "1" ]; then
    echo "* Upload ROM file to Arduino"
    echo "  - Usage: $0 <rom-binary-file>"
    echo ""
    exit
fi

if ! command -v arduino-cli &> /dev/null; then
    echo "The 'arduino-cli' program could not be found"
    exit
fi

if ! command -v ./rom2header/rom2header &> /dev/null; then
    echo "The 'rom2header' executable could not be found"
    exit
fi

# Alternatives are uni, nano etc.
ARDUINO="diecimila"
# Alternatives are atmega328, atmega328old
CPU="atmega168"
PORT="/dev/ttyUSB0"
ROMFILE=$1

# Convert the bin-file to a c-header file
./rom2header/rom2header -in $ROMFILE -out ./programs.h
mv ./programs.h arduino/EEPROMWriter/programs.h

# Compile the writer with the embedded binary file
echo ""
echo "== Compiling ($ARDUINO)"
arduino-cli compile -b arduino:avr:$ARDUINO:cpu=$CPU ./arduino/EEPROMWriter

# Upload the result to an Arduino
echo ""
echo "== Uploading to device ($PORT)"
arduino-cli upload -b arduino:avr:$ARDUINO:cpu=$CPU -p $PORT ./arduino/EEPROMWriter

echo ""
echo "== Done!"

echo ""
echo "Connect your Arduino to the EEPROM as described in README.md"
echo "file. Press the RESET button on the Arduino to start programming"
echo "the chip. The Pin13 LED on the Arduino turns off when the write"
echo "process has finished."
echo ""
