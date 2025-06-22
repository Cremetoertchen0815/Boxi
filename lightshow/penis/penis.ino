#include <Wire.h>


void setup() {
  Serial.begin(19200);
}

void loop() {
  if (Serial.available() < 1) {
    delay(10);
    return;
  }

  int lol = Serial.read();
  Serial.write(lol);
}
