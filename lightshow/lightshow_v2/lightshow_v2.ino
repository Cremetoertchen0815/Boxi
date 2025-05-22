#include <DmxSimple.h>
#include <Wire.h>
#include <Adafruit_PWMServoDriver.h>

#define COLOR_A 0b0000000011
#define COLOR_B 0b0000001100
#define COLOR_C 0b0000110000
#define COLOR_D 0b0000001111
#define COLOR_E 0b0000111100
#define COLOR_F 0b0000110011
#define COLOR_AMBER 0b0011000000
#define COLOR_UV 0b1100000000
#define MAX_PWM 4095

#define MUSIC_PIN 9
#define THRESH_PIN A1
#define MODE_PIN A2

//Modes
const int MODE_FAST_COLOR_SWITCH = 0;
const int MODE_FAST_COLOR_FLASH = 1;
const int MODE_FAST_STROBE = 2;
const int MODE_SLOW_FADE = 3;
const int MODE_SLOW_UV = 4;
const int MODE_SLOW_AMBER = 5;
//Power multipliers
const float RED_HARDWARE_MULTIPLIER = 1;
const float GREEN_HARDWARE_MULTIPLIER = 1;
const float BLUE_HARDWARE_MULTIPLIER = 1;
const float WHITE_HARDWARE_MULTIPLIER = 1;
const float AMBER_HARDWARE_MULTIPLIER = 1;
const float UV_HARDWARE_MULTIPLIER = 1;

const int FADE_SPEED = 1500;
const int STROBE_SPEED = 7;
const float Third = (float)1 / 3;
const float DIMMED_COLOR_BRIGHTNESS = 0.15; //The brightness of the color LEDs in pulse mode, when dimmed down(the peaks go up to max power)
const float STANDARD_COLOR_BRIGHTNESS = 0.6; //The brightness of the color LEDs in switch mode
const float FADE_COLOR_BRIGHTNESS = 0.3; //The brightness of the color LEDs in fade mode
const int BEAT_SHORTEST_SWITCH_TIME = 100; //The holding time between to music peaks to prevent too fast switching
const int BEAT_MIN_DURATION = 1; //The number if cycles in a row that the beat line has to be pulled high to count as a beat
const int BEAT_THRESHOLD_MIN = 25; //The beat threshold value when the potentiometer is turned all the way down
const int TIME_BEFORE_FADE_MODE = 1500;
const int POWER_THRESHOLD_OFF = 20;
const int POWER_THRESHOLD_FLOOD = 1000;
const int SENS_THRESHOLD_EAST_MODE = 28;
const int EAST_COUNTER_MAX = 10000;

Adafruit_PWMServoDriver pwm = Adafruit_PWMServoDriver();
int beatCheck = 0;
int beatPassed = -1;
int timeSinceLastBeat = 0;
int isOn = -1;
int lastWasOn = 1;
float currentColorBrightness = -1;
int currentColor = COLOR_A;
int nextColor = COLOR_B;
int colorTimer = 0;
int brightTimer = 0;
int mode = MODE_FAST_COLOR_SWITCH;
int modeSwitchTimer = 0;
int maxDiff = 0;
int resetTimer = 100;
int eastCounter = 0;
bool enableColorFade = false;

float lerp(float a, float b, float f) 
{
    return (a * (1.0 - f)) + (b * f);
}

void setup() {
  //Set up DMX
  DmxSimple.usePin(11);
  DmxSimple.maxChannel(30);
  
  pinMode(MUSIC_PIN, INPUT);

  pwm.begin();
  pwm.setOscillatorFrequency(27000000);
  pwm.setPWMFreq(1600);  // This is the maximum PWM frequency

  // if you want to really speed stuff up, you can go into 'fast 400khz I2C' mode
  // some i2c devices dont like this so much so if you're sharing the bus, watch
  // out for this!
  Wire.setClock(400000);
  
  //Get random seed
  randomSeed(analogRead(6));

  Serial.begin(9600);
}

void loop() {
  int musicVal = digitalRead(MUSIC_PIN);  // read the input pin
  int rawThreshold = analogRead(THRESH_PIN);
  int beatThreshold = map(1024 - rawThreshold, 0, 1023, BEAT_THRESHOLD_MIN, 1010);
  int power = analogRead(MODE_PIN);
  bool triggeredChange = false;

  //---BEAT DETECTION---0
  //-The filtering and thresholding gets done by the DSP, all we get is a digital signal, whether a beat is currently going on
  //-The beat is only considered if there were a certain number of HIGHs on the beat line consecutively. This is tracked in the beatCheck.
  //-The beat is only allowed to happen every x cycles. This is tracked in beatPassed. If it is smaller than 0, no beat happened yet
  timeSinceLastBeat++;
  if (musicVal == 1) beatCheck++; else beatCheck = 0;
  if (beatPassed >= 0 && beatPassed < BEAT_SHORTEST_SWITCH_TIME) beatPassed++;
  if (beatCheck >= BEAT_MIN_DURATION && (beatPassed < 0 || beatPassed >= BEAT_SHORTEST_SWITCH_TIME)) {
    beatPassed = 0;
    colorTimer = FADE_SPEED;
    brightTimer = 0;
    triggeredChange = true;
    timeSinceLastBeat = 0;
  }

  //---COLOR MIXING ALGORITH---
  //-Current running program calculates base colors(RGBAUV)
  //-Base color + overall color brightness + white brightness => aparrant RBGWAUV values

  //Calculate brightnesses
  float whiteBright = 0;
  float colorBrightMaster = 1;
  if (power <= POWER_THRESHOLD_OFF) {
    colorBrightMaster = 0;
  } else if (power >= POWER_THRESHOLD_FLOOD) {
    colorBrightMaster = 0;
    whiteBright = 1;
  } else {
    //The general color is determin
    colorBrightMaster = map(power, POWER_THRESHOLD_OFF, POWER_THRESHOLD_FLOOD, 1, 1000) * 0.001;
  }


  if (mode == MODE_FAST_COLOR_FLASH) 
      currentColorBrightness = exp(-brightTimer * 0.1) * 5 + DIMMED_COLOR_BRIGHTNESS;
  else if (mode == MODE_FAST_STROBE && power < POWER_THRESHOLD_FLOOD) {
      whiteBright = (int)((brightTimer / STROBE_SPEED) % 6) == 0;
      currentColorBrightness = 0;
  } else if (mode == MODE_FAST_COLOR_SWITCH) {
      currentColorBrightness = STANDARD_COLOR_BRIGHTNESS;
  }


  //Make lights pulse in mode 1
  if (mode != MODE_FAST_COLOR_FLASH) colorTimer++;

  //If there hasn't been an impulse in a while, move to a slow mode
  if (mode < MODE_SLOW_FADE && timeSinceLastBeat >= TIME_BEFORE_FADE_MODE) {
    mode = MODE_SLOW_FADE + random(0, 3);
    colorTimer = 0;
    nextColor = getNextColor();
  //Generate new color & tick color timer
  } else if (colorTimer >= FADE_SPEED) {
    currentColor = nextColor;
    nextColor = getNextColor();
    colorTimer = 0;
    
    //If in fade mode, fade to the common fade brightness
    if (mode >= MODE_SLOW_FADE) currentColorBrightness = FADE_COLOR_BRIGHTNESS;
  }
  if (brightTimer < 10000) brightTimer++;
  
  //Automatically switch hard modes after a while(only change when beat hits)
  //-Happens automatically after a while in fast mode to add variety
  //-If in slow mode any trigger will bring it back to a hard mode
  bool jumpingInFromSlow = mode >= MODE_SLOW_FADE;
  bool jumpingInFromStrobe = mode == MODE_FAST_STROBE;
  if ((--modeSwitchTimer <= 0 || jumpingInFromSlow) && triggeredChange) {
    mode = random(0, (jumpingInFromSlow ? 3 : 5)) % 3; //Randomly pick new mode, while prefering slow modes by double if not switching from a slow mode
    if (jumpingInFromStrobe && mode == MODE_FAST_STROBE) mode = MODE_FAST_COLOR_FLASH;
    //Determine time till switching to the next mode
    switch(mode) {
      case MODE_FAST_STROBE:
        modeSwitchTimer = random(50, 200);
        break;
      default:
        modeSwitchTimer = random(400, 7000);
        break;
    }

  //When changing back into fast mode, update the colors to something fresh
    if (jumpingInFromSlow) {
      currentColor = getColorByIndex(random(0, 6));
      nextColor = getNextColor();
    }
  }

  //Calculate current and future RGB values
  float tto = (mode < MODE_SLOW_FADE) ? 0 : ((float)colorTimer / FADE_SPEED);
  float colorBrightness = lerp(currentColorBrightness, FADE_COLOR_BRIGHTNESS, tto);
  colorBrightness *= colorBrightMaster;
  
  float currentR = (currentColor & 0b11) * Third;
  float currentG = ((currentColor >> 2) & 0b11) * Third;
  float currentB = ((currentColor >> 4) & 0b11) * Third;
  float currentA = ((currentColor >> 6) & 0b11) * Third;
  float currentU = ((currentColor >> 8) & 0b11) * Third;
  float nextR = (nextColor & 0b11) * Third;
  float nextG = ((nextColor >> 2) & 0b11) * Third;
  float nextB = ((nextColor >> 4) & 0b11) * Third;
  float nextA = ((nextColor >> 6) & 0b11) * Third;
  float nextU = ((nextColor >> 8) & 0b11) * Third;


  //If sensitivity is turned all the way down, activate EAST mode
  if (rawThreshold <= SENS_THRESHOLD_EAST_MODE) {
    currentR = 0.3;
    currentG = 0;
    currentB = 1;
    currentA = 0;
    currentU = 0;
    whiteBright = 0;
    nextR = 0;
    nextG = 1;
    nextB = 0.6;
    nextA = 0;
    nextU = 0;
    tto = (float)abs(eastCounter - EAST_COUNTER_MAX / 2) / EAST_COUNTER_MAX * 2;
    eastCounter = (eastCounter+1) % EAST_COUNTER_MAX;
    colorBrightness = colorBrightMaster;
    
    if (power >= POWER_THRESHOLD_FLOOD) whiteBright = 1;
  }

  //Prevent strobe oder flood from going off when the lighting is turned off
  if (power <= POWER_THRESHOLD_OFF) {
    whiteBright = 0;
  }

  //Calculate final RGBW values
  uint16_t currR = lerp(currentR, nextR, tto) * MAX_PWM * colorBrightness * RED_HARDWARE_MULTIPLIER;
  uint16_t currG = lerp(currentG, nextG, tto) * MAX_PWM * colorBrightness * GREEN_HARDWARE_MULTIPLIER;
  uint16_t currB = lerp(currentB, nextB, tto) * MAX_PWM * colorBrightness * BLUE_HARDWARE_MULTIPLIER;
  uint16_t currW = MAX_PWM * whiteBright * WHITE_HARDWARE_MULTIPLIER;
  uint16_t currA = lerp(currentA, nextA, tto) * MAX_PWM * colorBrightness * AMBER_HARDWARE_MULTIPLIER;
  uint16_t currU = lerp(currentU, nextU, tto) * MAX_PWM * colorBrightness * UV_HARDWARE_MULTIPLIER;

  //Update internal PWM signals accordingly
  pwm.setPWM(0, 0, currR);
  pwm.setPWM(1, 0, currG);
  pwm.setPWM(2, 0, currB);
  pwm.setPWM(3, 0, currW);
  pwm.setPWM(7, 0, currA);
  pwm.setPWM(5, 0, currU);

  //Update external lighting via DMX
  DmxSimple.write(1, currR / 16);
  DmxSimple.write(2, currG / 16);
  DmxSimple.write(3, currB / 16);
  DmxSimple.write(4, currW / 16);
  DmxSimple.write(5, currA / 16);
  DmxSimple.write(6, currU / 16);
  DmxSimple.write(7, currB / 16);
  DmxSimple.write(8, currR / 16);
  DmxSimple.write(9, currG / 16);
  DmxSimple.write(10, currW / 16);
  DmxSimple.write(11, currA / 16);
  DmxSimple.write(12, currU / 16);
  DmxSimple.write(13, currG / 16);
  DmxSimple.write(14, currB / 16);
  DmxSimple.write(15, currR / 16);
  DmxSimple.write(16, currW / 16);
  DmxSimple.write(17, currA / 16);
  DmxSimple.write(18, currU / 16);

  DmxSimple.write(19, currR / 16);
  DmxSimple.write(20, currG / 16);
  DmxSimple.write(21, currB / 16);
  DmxSimple.write(22, currU / 16);
  DmxSimple.write(23, currB / 16);
  DmxSimple.write(24, currR / 16);
  DmxSimple.write(25, currG / 16);
  DmxSimple.write(26, currU / 16);
  DmxSimple.write(27, currG / 16);
  DmxSimple.write(28, currB / 16);
  DmxSimple.write(29, currR / 16);
  DmxSimple.write(30, currU / 16);
}

int getNextColor() {
  //In those special fade modes only one color should be used
  if (mode == MODE_SLOW_UV) {
    return COLOR_UV;
  } else if (mode == MODE_SLOW_AMBER) {
    return COLOR_AMBER;
  }

  //In any other mode, just switch between all colors
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

int getColorByIndex(int idx) {
  switch (idx) {
    case 0:
      return COLOR_A;
    case 1:
      return COLOR_B;
    case 2:
      return COLOR_C;
    case 3:
      return COLOR_D;
    case 4:
      return COLOR_E;
    default:
      return COLOR_F;
  }
}
