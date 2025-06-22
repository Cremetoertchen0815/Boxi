#include <Wire.h>

#define MUSIC_PIN      7
const int BEAT_SHORTEST_SWITCH_TIME = 20; //The holding time between to music peaks to prevent too fast switching
const int BEAT_MIN_DURATION = 1; //The number if cycles in a row that the beat line has to be pulled high to count as a beat
int beatCheck = 0;
int beatPassed = -1;
int timeSinceLastBeat = 0;

void setup() {
  Serial.begin(19200);
}

void loop() {
  bool beat = checkForBeat();

  if (beat) {
    Serial.write("Beat.\n");
  }
}


//Checks for a beat
bool checkForBeat() {
  //The filtering and thresholding gets done by the DSP, all we get is a digital signal, whether a beat is currently going on
  //The beat is only considered if there were a certain number of HIGHs on the beat line consecutively. This is tracked in the beatCheck.
  //The beat is only allowed to happen every x cycles. This is tracked in beatPassed. If it is smaller than 0, no beat happened yet
  int musicVal = digitalRead(MUSIC_PIN);  // read the input pin
  timeSinceLastBeat++;
  if (musicVal == 1) beatCheck++; else beatCheck = 0;
  if (beatPassed >= 0 && beatPassed < BEAT_SHORTEST_SWITCH_TIME) beatPassed++;
  if (beatCheck >= BEAT_MIN_DURATION && (beatPassed < 0 || beatPassed >= BEAT_SHORTEST_SWITCH_TIME)) {
    beatPassed = 0;
    timeSinceLastBeat = 0;
    return true;
  }

  return false;
}

void checkSerial() {
  if (Serial.available() < 1) {
    delay(10);
    return;
  }

  int lol = Serial.read();
  Serial.write(lol);
}