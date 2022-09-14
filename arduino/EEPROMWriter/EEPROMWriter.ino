
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


const int EEPROM_ADDR_CMD = 0x50;
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
 
  Wire.begin(); // Initialise Wire library
  
  Serial.begin(115200);
  
  Serial.println("== EEPROM Writer startup ==");
  Serial.print("Number of programs: ");
  Serial.println(NUM_PROGRAMS);

  // 
  // Writing the data to the EEPROM
  //
  digitalWrite(pinLED, HIGH);  

  int counter = 0;
  for (int i=0; i<NUM_PROGRAMS; i++) {
    Serial.print("Writing program ");  
    Serial.print(i);
    Serial.print(" (");
    Serial.print(PROGRAM_SIZE);
    Serial.println(" bytes).");

    const unsigned char * data = PROGRAMS[i]; 

    const int chunksize = 1;
    for (int j=0; j<PROGRAM_SIZE; j+= chunksize){
      const unsigned int eeaddress = i*PROGRAM_SIZE + j;
      Wire.beginTransmission(EEPROM_ADDR_CMD);
      Wire.write((int)(eeaddress >> 8));    // MSB
      Wire.write((int)(eeaddress & 0xFF));  // LSB
      for (unsigned int k = 0; k < chunksize; k++ ) {
        Wire.write(pgm_read_byte_near(data + (counter++)));  // Write to EEPROM
      }
      Wire.endTransmission();
      delay(10); // Small delay
    }
  }
  
  //
  // Verify the program
  //
  counter = 0;
  for (int i=0; i<NUM_PROGRAMS; i++) {
    Serial.print("Verifying program ");  
    Serial.print(i);

    const unsigned char * data = PROGRAMS[i];  // Pointer to data to verify against EEPROM

    int nBytesOK = 0;  
    int nBytesFail = 0;
    
    for (int j = 0; j<PROGRAM_SIZE; j++) {
      const unsigned int eeaddress = i*PROGRAM_SIZE + j;  // EEPROM address

      Wire.beginTransmission(EEPROM_ADDR_CMD);
      Wire.write((int)(eeaddress >> 8));    // MSB
      Wire.write((int)(eeaddress & 0xFF));  // LSB
      Wire.endTransmission();
      Wire.requestFrom(EEPROM_ADDR_CMD, 1);
      //delay(15);
      if (Wire.available()) {
        const byte read = Wire.read();  // Read byte from EEPROM
        if (read == pgm_read_byte_near(data + (counter++))) { // Check it is as expected
          nBytesOK++;
        }
        else {
          nBytesFail++;         
        }
      }
      else {
        // A "read" is not available. This means trouble...
      }
      
      if ((j % 32) == 31) {
        delay(10);      // Small delay every 32 bytes
      }
    }
    
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
  
  digitalWrite(pinLED, LOW);

  
  Serial.println("EEPROM written and verified...\n");
}

void loop() {
  // Do nothing
}
