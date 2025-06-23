#include <Wire.h>
#include <DmxSimple.h>

void setup() {
  pwm.begin();
  pwm.setOscillatorFrequency(27000000);
  pwm.setPWMFreq(1600);  // This is the maximum PWM frequency

  // if you want to really speed stuff up, you can go into 'fast 400khz I2C' mode
  // some i2c devices dont like this so much so if you're sharing the bus, watch
  // out for this!
  Wire.setClock(400000);

  //Set up DMX
  DmxSimple.usePin(6);
  DmxSimple.maxChannel(1);
}

void loop() {
  DmxSimple.write(1, 100);
  DmxSimple.write(2, 200);
  DmxSimple.write(3, 254);
}