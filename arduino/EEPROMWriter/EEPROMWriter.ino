
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


const int EEPROM_ADDR = 0x50;

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
  
  Serial.begin(9200);
  Serial.print("EEPROM Writer startup");
  Serial.print("Number of programs: ");
  Serial.println(NUM_PROGRAMS);

  // 
  // Writing the data to the EEPROM
  //
  digitalWrite(pinLED, HIGH);  
  for (int i=0; i<NUM_PROGRAMS; i++) {
    Serial.print("Writing program ");  
    Serial.println(i);

    const unsigned char * data = PROGRAMS[i]; 

    for (int j=0; j<512; j+= 16){
      const unsigned int eeaddress = i*512 + j;
      Wire.beginTransmission(EEPROM_ADDR);
      Wire.write((int)(eeaddress >> 8));    // MSB
      Wire.write((int)(eeaddress & 0xFF));  // LSB
      for (unsigned int k = 0; k < 16; k++ ) {
        Wire.write(pgm_read_byte_near(data + j + k));  // Write to EEPROM
      }
      Wire.endTransmission();
      delay(10); // Small delay
    }
  }
  
  /*
  // Verify ...
  for (int i=0; i<NUM_PROGRAMS; i++) {
    Serial.print("Verifying program ");  
    Serial.print(i);

    const unsigned char * data = PROGRAMS[i];  // Pointer to data to verify against EEPROM

    int nBytesOK(0);  // Number of bytes that were written OK
    for (int j = 0; j<512; j++) {
      const unsigned int eeaddress = i*512 + j;  // EEPROM address

      Wire.beginTransmission(EEPROM_ADDR);
      Wire.write((int)(eeaddress >> 8));    // MSB
      Wire.write((int)(eeaddress & 0xFF));  // LSB
      Wire.endTransmission();
      Wire.requestFrom(EEPROM_ADDR, 1);
      if (Wire.available()) {
        const byte read = Wire.read();  // Read byte from EEPROM
        if (read == pgm_read_byte_near(data + i)) { // Check it is as expected
          nBytesOK++;
        }
      }
      if ((j % 32) == 31) {
        delay(10);      // Small delay every 32 bytes
      }
    }
    Serial.print(": "); 
    Serial.print(nBytesOK); 
    Serial.println(" Bytes written correctly");
  }
  */
  digitalWrite(pinLED, LOW);

  
  Serial.println("EEPROM written and verified :)");
}

void loop() {
  //while(true) {};
}
