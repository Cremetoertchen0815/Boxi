#include <Wire.h>
#include <Adafruit_PWMServoDriver.h>

Adafruit_PWMServoDriver pwm = Adafruit_PWMServoDriver();

void setup() {
  pwm.begin();
  pwm.setOscillatorFrequency(27000000);
  pwm.setPWMFreq(1600);  // This is the maximum PWM frequency

  // if you want to really speed stuff up, you can go into 'fast 400khz I2C' mode
  // some i2c devices dont like this so much so if you're sharing the bus, watch
  // out for this!
  Wire.setClock(400000);

  Serial.begin(19200);
}

void loop() {
  for(int i = 0; i < 6; i++) {
    
  pwm.setPin(i, 0, 50);
  delay(500);
  pwm.setPin(i, 0, 0);
  delay(500);
  }
}