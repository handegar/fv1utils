
#include <Wire.h> //I2C library
#include <avr/pgmspace.h>
#include "programs.h"

/*
 * How to hook up the Arduino:
 * - Arduino Diecimila:
 *     3V3: V3.3d  
 *     GND: Gnd 
 *     Analog In 4: SADR (aka 'SDA') 
 *     Analog In 5: SCLK (aka 'SCL')
 * - Arduino Nano:
 *     (To be written, most likely same as above)
 */


#define VERIFY_EEPROM 1

const int EEPROM_ADDR_CMD = 0x50;
const int  EEPROM_CLOCK_SPEED = 400000;
const int PROGRAM_SIZE = 128*4;

void setup() {  
  delay(500);

  const int pinSDA = 4;  // SDA = Pin A4.  On the Uno Revision 3 board, this is the same as pin 16.
  const int pinSCL = 5;  // SCL = Pin A5.  On the Uno Revision 3 board, this is the same as pin 17.
  const int pinLED = 13; // The builtin LED

  pinMode(pinSDA, OUTPUT);
  pinMode(pinSCL, OUTPUT);
  digitalWrite(pinSDA, HIGH);
  digitalWrite(pinSCL, LOW);

  Serial.begin(115200);
    
  Wire.begin(); // Initialise Wire library
  Wire.setClock(EEPROM_CLOCK_SPEED);   
  
  Serial.println("== EEPROM Writer startup ==");
  Serial.print("Number of programs: ");
  Serial.println(NUM_PROGRAMS);

  // 
  // Writing the data to the EEPROM
  //
  digitalWrite(pinLED, HIGH);  

  int writtenTotal = 0;
  for (int i=0; i<NUM_PROGRAMS; i++) {
    Serial.print("Writing program ");  
    Serial.print(i);
    Serial.print(" (");
    Serial.print(PROGRAM_SIZE);
    Serial.println(" bytes).");

    const unsigned char * data = PROGRAMS[i]; 
    int written = 0;
    const int chunksize = 16;
    for (int j=0; j<PROGRAM_SIZE; j+= chunksize){
      const unsigned int eeaddress = written + writtenTotal;
      Wire.beginTransmission(EEPROM_ADDR_CMD);
      Wire.write((int)(eeaddress >> 8));    // MSB
      Wire.write((int)(eeaddress & 0xFF));  // LSB
      for (unsigned int k = 0; k < chunksize; k++ ) {
        Wire.write(pgm_read_byte_near(data + written));  // Write to EEPROM
        written += 1;
      }
      Wire.endTransmission();
      delay(10); // Small delay
    }
    writtenTotal += written;
  }

  Serial.print("Wrote ");
  Serial.print(writtenTotal);
  Serial.println(" bytes in total.");
   
  //
  // Verify the program
  //
#if VERIFY_EEPROM
  int readTotal = 0;
  for (int i=0; i<NUM_PROGRAMS; i++) {
    Serial.print("Verifying program ");  
    Serial.print(i);

    const unsigned char * data = PROGRAMS[i];  // Pointer to data to verify against EEPROM

    int nBytesOK = 0;  
    int nBytesFail = 0;
    int readBytes = 0;
    
    for (int j = 0; j<PROGRAM_SIZE; j++) {
      const unsigned int eeaddress = readBytes + readTotal;

      Wire.beginTransmission(EEPROM_ADDR_CMD);
      Wire.write((int)(eeaddress >> 8));    // MSB
      Wire.write((int)(eeaddress & 0xFF));  // LSB
      Wire.endTransmission();
      Wire.requestFrom(EEPROM_ADDR_CMD, 1);
      //delay(15);
      if (Wire.available()) {
        const byte read = Wire.read();  // Read byte from EEPROM
        if (read == pgm_read_byte_near(data + readBytes)) { // Check it is as expected
          nBytesOK++;
        }
        else {
          nBytesFail++;         
        }
        readBytes += 1;
      }
      else {
        // A "read" is not available. This means trouble...
      }
      
      if (readBytes % 16) {
        delay(10);      // Small delay every 16 bytes
      }
    }

    readTotal += readBytes;
    
    Serial.print(": "); 
    Serial.print(nBytesOK); 
    Serial.print(" bytes written correctly. ");
    if (nBytesFail > 0) {
      Serial.print(nBytesFail); 
      Serial.println(" bytes NOT correct.");
    }
    else {
      Serial.println("");   
    }
  }
#endif

  digitalWrite(pinLED, LOW);

  Serial.print("Read and verified ");
  Serial.print(readTotal);
  Serial.println(" bytes.");
  Serial.println("EEPROM written and verified...\n");
}

void loop() {
  // Do nothing
}
