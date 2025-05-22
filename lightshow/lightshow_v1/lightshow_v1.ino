#include <RunningMedian.h>
#include <DmxSimple.h>

#define COLOR_A 0b00000011
#define COLOR_B 0b00001100
#define COLOR_C 0b00110000
#define COLOR_D 0b00001111
#define COLOR_E 0b00111100
#define COLOR_F 0b00110011
#define RED_PIN 6
#define GREEN_PIN 5
#define BLUE_PIN 9
#define WHITE_PIN 10
#define MUSIC_PIN A5
#define THRESH_PIN A3
#define MODE_PIN A1

#define CORRECT_CLOCK 64
void fixDelay(uint32_t ms) {
  delay(ms << CORRECT_CLOCK);
}

const int FADE_SPEED = 1500;
const int FAST_SWITCH_SPEED = 35;
const int STROBE_SPEED = 7;
const float Third = (float)1 / 3;
const float DIMMED_COLOR_BRIGHTNESS = 0.15; //The brightness of the color LEDs in pulse mode, when dimmed down(the peaks go up to max power)
const float STANDARD_COLOR_BRIGHTNESS = 0.5; //The brightness of the color LEDs in switch mode
const float WHITE_LED_FLOOD_BRIGHTNESS = 1; //The brightness of the white LEDs in flood mode(shouldn't be close to 1, to prevent heat issues in long term usage)
const float COLOR_LED_BRIGHTNESS_LIMIT = 1; //Absolute power multiplier for the color LEDs to prevent them overheating
const float WHITE_LED_BRIGHTNESS_LIMIT = 1; //Absolute power multiplier for the white LEDs to prevent them overheating
const int BEAT_SHORTEST_SWITCH_TIME = 70; //The holding time between to music peaks to prevent too fast switching
const int BEAT_THRESHOLD_MIN = 25; //The beat threshold value when the potentiometer is turned all the way down
const int BEAT_THRESHOLD_MAX = 250; //The beat threshold value when the potentiometer is turned all the way up
const int POWER_THRESHOLD_OFF = 20;
const int POWER_THRESHOLD_FLOOD = 1000;
RunningMedian samples = RunningMedian(64);
int isOn = -1;
int lastWasOn = 1;
int mean = 462;
int currentColor = COLOR_A;
int nextColor = COLOR_B;
int colorTimer = 0;
int brightTimer = 0;
int mode = 2;
int modeSwitchTimer = 0;
int maxDiff = 0;
int resetTimer = 100;

float lerp(float a, float b, float f) 
{
    return (a * (1.0 - f)) + (b * f);
}

void setup() {

  //Set up DMX
  DmxSimple.usePin(11);
  DmxSimple.maxChannel(7);
  
  //Prepare MOSFET outputs
  pinMode(RED_PIN, OUTPUT);
  pinMode(GREEN_PIN, OUTPUT);
  pinMode(BLUE_PIN, OUTPUT);
  pinMode(WHITE_PIN, OUTPUT);

  //Setup timers
  // Pins D9 and D10 - 31.4 kHz
  TCCR1B = 0b00000001; // x1  
  TCCR1A = 0b00000001; // phase correct
  // Pins D5 and D6 - 31.4 kHz  
  TCCR0B = 0b00000001; // x1
  TCCR0A = 0b00000001; // phase correct
  
  //Get random seed
  randomSeed(analogRead(0));
}

void loop() {
  int musicVal = analogRead(MUSIC_PIN);  // read the input pin
  int beatThreshold = map(1024 - analogRead(THRESH_PIN), 0, 1023, BEAT_THRESHOLD_MIN, BEAT_THRESHOLD_MAX);
  int power = analogRead(MODE_PIN);

  //Calculate music transient
  bool triggeredChange = false;
  int musicMean = samples.getAverage();
  int musicDiff = musicMean - musicVal;
  if (musicDiff < 0) musicDiff *= -1;
  samples.add(musicVal);

  //Detect music peaks
  if (musicDiff > beatThreshold && isOn < 0 && !lastWasOn ) {
    triggeredChange = true;
    isOn = BEAT_SHORTEST_SWITCH_TIME;
    colorTimer = FADE_SPEED;
    brightTimer = 0;
  }
    lastWasOn = musicDiff > beatThreshold;

  isOn--;

  fixDelay(1);

  float whiteBright = 0;
  float colorBright = 1;
  if (power <= POWER_THRESHOLD_OFF) {
    colorBright = 0;
  } else if (power >= POWER_THRESHOLD_FLOOD) {
    colorBright = 0;
    whiteBright = WHITE_LED_FLOOD_BRIGHTNESS;
  } else {
    colorBright = map(power, POWER_THRESHOLD_OFF, POWER_THRESHOLD_FLOOD, 1, 1000) * 0.001;
    if (mode == 0) colorBright *= STANDARD_COLOR_BRIGHTNESS;
    if (mode == 1) colorBright *= exp(-brightTimer * 0.1) * 5 + DIMMED_COLOR_BRIGHTNESS;
    if (mode == 2) {
      whiteBright = (int)((brightTimer / STROBE_SPEED) % 6) == 0;
      colorBright *= DIMMED_COLOR_BRIGHTNESS;
    }
  }

  //Make lights pulse in mode 1
  if (mode != 1) colorTimer++;

  //Enforce maximum power limits to protect the LEDs
  colorBright *= COLOR_LED_BRIGHTNESS_LIMIT;
  whiteBright *= WHITE_LED_BRIGHTNESS_LIMIT;

  //Generate new color & tick color timer
  if (colorTimer >= FADE_SPEED) {
    currentColor = nextColor;
    nextColor = getNextColor();
    colorTimer = 0;
  }
  if (brightTimer < 10000) brightTimer++;
  
  //Automatically switch modes after a while(only change from slow modes when beat hits, fast modes can leave any time)
  if (--modeSwitchTimer <= 0 && (triggeredChange || mode > 1)) {
    mode = random(0, 5) % 3; //Randomly pick new mode, while prefering slow modes by double
    //Determine time till switching to the next mode
    switch(mode) {
      case 2:
        modeSwitchTimer = random(70, 600);
        break;
      default:
        modeSwitchTimer = random(400, 7000);
        break;
    }
  }

  //Calculate current and future RGB values
  float tto = (float)colorTimer / FADE_SPEED;
  float currentR = (currentColor & 0b11) * Third;
  float currentG = ((currentColor >> 2) & 0b11) * Third;
  float currentB = ((currentColor >> 4) & 0b11) * Third;
  float nextR = (nextColor & 0b11) * Third;
  float nextG = ((nextColor >> 2) & 0b11) * Third;
  float nextB = ((nextColor >> 4) & 0b11) * Third;

  //Calculate final RGBW values
  int currR = lerp(currentR, nextR, tto) * 255 * colorBright;
  int currG = lerp(currentG, nextG, tto) * 255 * colorBright;
  int currB = lerp(currentB, nextB, tto) * 255 * colorBright;
  int currW = 255 * whiteBright;

  //Update internal PWM signals accordingly
  analogWrite(RED_PIN, currR);
  analogWrite(GREEN_PIN, currG);
  analogWrite(BLUE_PIN, currB);
  analogWrite(WHITE_PIN, currW);

  //Update external lighting via DMX
  DmxSimple.write(1, currR);
  DmxSimple.write(2, currG);
  DmxSimple.write(3, currB);
  DmxSimple.write(4, currR);
  DmxSimple.write(5, currG);
  DmxSimple.write(6, currB);
}

int getNextColor() {
  switch(currentColor) {
    case COLOR_A:
      return COLOR_B;
    case COLOR_B:
      return COLOR_C;
    case COLOR_C:
      return COLOR_D;
    case COLOR_D:
      return COLOR_E;
    case COLOR_E:
      return COLOR_F;
    default:
      return COLOR_A;
  }
}
