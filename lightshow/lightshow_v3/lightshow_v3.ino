#include <DmxSimple.h>
#include <Wire.h>
#include <Adafruit_PWMServoDriver.h>
#include <Adafruit_GFX.h>    // Core graphics library
#include <Adafruit_ST7735.h> // Hardware-specific library for ST7735
#include <SPI.h>

#define MAX_PWM 4095
#define MAX_DMX 255
#define MUSIC_PIN      7
#define BRIGHTNESS_PIN A6
#define TFT_CS         10
#define TFT_RST        3
#define TFT_DC         9

enum DataField {
  DISPLAY_STATUS_CODE = 0x01,
  LIGHTING_APPLY = 0x02,
  LIGHTING_MODE = 0x03,
  LIGHTING_COLOR_SHIFT = 0x04,
  LIGHTING_SPEED = 0x05,
  LIGHTING_GENERAL_PURPOSE = 0x06,
  LIGHTING_PALLETTE_SIZE = 0x07,
  LIGHTING_PALLETTE_A = 0x08,
  LIGHTING_PALLETTE_B = 0x09,
  LIGHTING_PALLETTE_C = 0x0A,
  LIGHTING_PALLETTE_D = 0x0B,
  LIGHTING_PALLETTE_E = 0x0C,
  LIGHTING_PALLETTE_F = 0x0D,
  LIGHTING_PALLETTE_G = 0x0E,
  LIGHTING_PALLETTE_H = 0x0F,
}; 

enum LightingMode {
  OFF = 0x00,
  SET_COLOR = 0x01,
  FADE_TO_COLOR = 0x02,
  PALLETTE_FADE = 0x03,
  PALLETTE_SWITCH = 0x04,
  PALLETTE_FLASH = 0x05,
  BI_COLOR_FLASH = 0x06,
  STROBE = 0x07,
}; 

enum DisplayStatusCode {
  BOOTING = 0x00,
  HOST_AWAKE = 0x01,
  HOST_NO_ACTIVITY = 0x02,
  DSP_SERVER_FAILED = 0x03,
  HOST_CONNECTION_FAILED = 0x04,
  ACTIVE = 0x05
}; 

struct Color {
  uint8_t Red;
  uint8_t Green;
  uint8_t Blue;
  uint8_t White;
  uint8_t Amber;
  uint8_t UltraViolet;
}; 

struct FloatColor {
  float Red;
  float Green;
  float Blue;
  float White;
  float Amber;
  float UltraViolet;
}; 

struct DualColor {
  FloatColor Boxi1;
  FloatColor Boxi2;
}; 

struct DataFieldSet {
  Color Pallette[8];
  uint8_t PalletteSize;
  enum LightingMode Mode;
  uint8_t ColorShift;
  uint16_t Speed;
  uint8_t GeneralPurpose;
}; 

const int STROBE_SPEED = 7;
const float BYTE_TO_FLOAT = 1.0 / 255;
const float FADE_COLOR_BRIGHTNESS = 0.3; //The brightness of the color LEDs in fade mode
const int BEAT_SHORTEST_SWITCH_TIME = 100; //The holding time between to music peaks to prevent too fast switching
const int BEAT_MIN_DURATION = 1; //The number if cycles in a row that the beat line has to be pulled high to count as a beat
const int POWER_THRESHOLD_OFF = 20;
const int POWER_THRESHOLD_MAX = 1000;
const uint16_t HOST_CONNECTION_TIMEOUT = 1000;

//Variables
int beatCheck = 0;
int beatPassed = -1;
int timeSinceLastBeat = 0;
int referenceIndex = 0;
int referenceCounter = 0;
uint16_t hostConnectionCounter = HOST_CONNECTION_TIMEOUT;
DualColor referenceColor;
DualColor lastOutputColor;
bool updateReferenceColor = false;
bool applyLightingOnNextBeat = false;
DataFieldSet fieldSetA;
DataFieldSet fieldSetB;
uint8_t activeField = 0;
DisplayStatusCode dspStatusCode = BOOTING;

Adafruit_ST7735 tft = Adafruit_ST7735(TFT_CS, TFT_DC, TFT_RST);
Adafruit_PWMServoDriver pwm = Adafruit_PWMServoDriver();

float clampFloat(float value) {
  float minRes = value < 1 ? value : 1;
  return minRes > 0 ? minRes : 0;
}

float lerp(float a, float b, float f) 
{
    return (a * (1.0 - f)) + (b * f);
}

FloatColor convertColor(Color color) {
    FloatColor ret;
    ret.Red = color.Red * BYTE_TO_FLOAT;
    ret.Green = color.Green * BYTE_TO_FLOAT;
    ret.Blue = color.Blue * BYTE_TO_FLOAT;
    ret.White = color.White * BYTE_TO_FLOAT;
    ret.Amber = color.Amber * BYTE_TO_FLOAT;
    ret.UltraViolet = color.UltraViolet * BYTE_TO_FLOAT;
    return ret;
}

FloatColor lerpColor(FloatColor a, Color b, float f) 
{
    FloatColor ret;
    ret.Red = lerp(a.Red, (float)b.Red * BYTE_TO_FLOAT, f);
    ret.Green = lerp(a.Green, (float)b.Green * BYTE_TO_FLOAT, f);
    ret.Blue = lerp(a.Blue, (float)b.Blue * BYTE_TO_FLOAT, f);
    ret.White = lerp(a.White, (float)b.White * BYTE_TO_FLOAT, f);
    ret.Amber = lerp(a.Amber, (float)b.Amber * BYTE_TO_FLOAT, f);
    ret.UltraViolet = lerp(a.UltraViolet, (float)b.UltraViolet * BYTE_TO_FLOAT, f);
    return ret;
}

FloatColor lerpColor(Color a, Color b, float f) 
{
    FloatColor ret;
    ret.Red = lerp(a.Red * BYTE_TO_FLOAT, b.Red * BYTE_TO_FLOAT, f);
    ret.Green = lerp(a.Green * BYTE_TO_FLOAT, b.Green * BYTE_TO_FLOAT, f);
    ret.Blue = lerp(a.Blue * BYTE_TO_FLOAT, b.Blue * BYTE_TO_FLOAT, f);
    ret.White = lerp(a.White * BYTE_TO_FLOAT, b.White * BYTE_TO_FLOAT, f);
    ret.Amber = lerp(a.Amber * BYTE_TO_FLOAT, b.Amber * BYTE_TO_FLOAT, f);
    ret.UltraViolet = lerp(a.UltraViolet * BYTE_TO_FLOAT, b.UltraViolet * BYTE_TO_FLOAT, f);
    return ret;
}

DualColor multiplyDualColor(DualColor a, float f) 
{
    DualColor ret;
    ret.Boxi1.Red = a.Boxi1.Red * f;
    ret.Boxi1.Green = a.Boxi1.Green * f;
    ret.Boxi1.Blue = a.Boxi1.Blue * f;
    ret.Boxi1.White = a.Boxi1.White * f;
    ret.Boxi1.Amber = a.Boxi1.Amber * f;
    ret.Boxi1.UltraViolet = a.Boxi1.UltraViolet * f;
    ret.Boxi2.Red = a.Boxi2.Red * f;
    ret.Boxi2.Green = a.Boxi2.Green * f;
    ret.Boxi2.Blue = a.Boxi2.Blue * f;
    ret.Boxi2.White = a.Boxi2.White * f;
    ret.Boxi2.Amber = a.Boxi2.Amber * f;
    ret.Boxi2.UltraViolet = a.Boxi2.UltraViolet * f;
    return ret;
}

FloatColor multiplyColor(Color a, float f) 
{
    FloatColor ret;
    ret.Red = (float)a.Red * BYTE_TO_FLOAT * f;
    ret.Green = (float)a.Green * BYTE_TO_FLOAT * f;
    ret.Blue = (float)a.Blue * BYTE_TO_FLOAT * f;
    ret.White = (float)a.White * BYTE_TO_FLOAT * f;
    ret.Amber = (float)a.Amber * BYTE_TO_FLOAT * f;
    ret.UltraViolet = (float)a.UltraViolet * BYTE_TO_FLOAT * f;
    return ret;
}

void createDefaultMode() {
  fieldSetA.Pallette[0] = {255, 0, 0, 0, 0, 0};
  fieldSetA.Pallette[1] = {255, 255, 0, 0, 0, 0};
  fieldSetA.Pallette[2] = {0, 255, 0, 0, 0, 0};
  fieldSetA.Pallette[3] = {0, 255, 255, 0, 0, 0};
  fieldSetA.Pallette[4] = {0, 0, 255, 0, 0, 0};
  fieldSetA.Pallette[5] = {255, 0, 255, 0, 0, 0};
  fieldSetA.PalletteSize = 6;
  fieldSetA.Mode = PALLETTE_FADE;
  fieldSetA.Speed = 200;
}

void handleDisplayStatusCode(DisplayStatusCode statusCode) {
  dspStatusCode = statusCode;
  char* textToPrint;

  switch(statusCode) {
    case BOOTING:
      textToPrint = "Booting...";
      break;
    case HOST_AWAKE:
      return;
    case HOST_NO_ACTIVITY:
      textToPrint = "ERR_0x02: HOST_NO_ACTIVITY";
      break;
    case DSP_SERVER_FAILED:
      textToPrint = "ERR_0x03: DSP_SERVER_FAILED";
      break;
    case HOST_CONNECTION_FAILED:
      textToPrint = "ERR_0x04: HOST_CONNECTION_FAILED";
      break;
    case ACTIVE:
      textToPrint = "Booting complete!";
      break;
    default:
      textToPrint = "ERR_0xFF: UNKNOWN_ERR";
      break;
  }


  tft.fillRect(0, 119, 160, 9, 0x0000);

  tft.setCursor(0, 119);
  tft.setTextColor(0XFFFF);
  tft.print(textToPrint);
}

void checkHostActivity() {
  if (hostConnectionCounter > 0) {
    hostConnectionCounter--;
    
  } else if (hostConnectionCounter == 0 && dspStatusCode == BOOTING) {
    handleDisplayStatusCode(HOST_NO_ACTIVITY);
  }
}

void applyLighting() {
  activeField = !activeField;
  updateReferenceColor = true;
  referenceCounter = 0;
  referenceIndex = 0;
}

//Processes the incoming data from the RPI at the UART port
void processUart() {
  //Make sure header is fine before handling data
  if (Serial.available() < 4 || 
      Serial.read() != 0x24 ||  
      Serial.read() != 0x20 ||  
      Serial.read() != 0x1F) return;

  //Read field
  enum DataField field = Serial.read();

  //Determine how much data should be received
  uint8_t dataLen = 0;
  switch(field) {
    case DISPLAY_STATUS_CODE:
    case LIGHTING_APPLY:
    case LIGHTING_MODE:
    case LIGHTING_PALLETTE_SIZE:
    case LIGHTING_COLOR_SHIFT:
    case LIGHTING_GENERAL_PURPOSE:
      dataLen = 1;
      break;
    case LIGHTING_SPEED:
      dataLen = 2;
      break;
    case LIGHTING_PALLETTE_A:
    case LIGHTING_PALLETTE_B:
    case LIGHTING_PALLETTE_C:
    case LIGHTING_PALLETTE_D:
    case LIGHTING_PALLETTE_E:
    case LIGHTING_PALLETTE_F:
    case LIGHTING_PALLETTE_G:
    case LIGHTING_PALLETTE_H:
      dataLen = 6;
      break;
    default:
      return;
  }

  //Read payload
  uint8_t receivedData[dataLen];
  if (Serial.readBytes(receivedData, dataLen) != dataLen) {
    return;
  }

  //Do stuff with data
  DataFieldSet* fieldSet = activeField == 0 ? &fieldSetB : &fieldSetA;

  switch(field) {
    case DISPLAY_STATUS_CODE:
      handleDisplayStatusCode(receivedData[0]);
      break;
    case LIGHTING_APPLY:
      if (receivedData[0] != 0) {
        applyLightingOnNextBeat = true;
      } else {
        applyLighting();
      }
      break;
    case LIGHTING_MODE:
      fieldSet->Mode = receivedData[0];
      break;
    case LIGHTING_SPEED:
      fieldSet->Speed = ((uint16_t)receivedData[0] << 8) + receivedData[1];
      break;
    case LIGHTING_PALLETTE_SIZE:
      fieldSet->PalletteSize = receivedData[0] > 8 ? 8 : receivedData[0];
      break;
    case LIGHTING_COLOR_SHIFT:
      fieldSet->ColorShift = receivedData[0];
      break;
    case LIGHTING_GENERAL_PURPOSE:
      fieldSet->GeneralPurpose = receivedData[0];
      break;
    default:
      uint8_t index = field - LIGHTING_PALLETTE_A;
      if (index < 0 || index > 7) return;

      Color* color = &fieldSet->Pallette[index];
      color->Red = receivedData[0];
      color->Green = receivedData[1];
      color->Blue = receivedData[2];
      color->White = receivedData[3];
      color->Amber = receivedData[4];
      color->UltraViolet = receivedData[5];
      break;
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

//Reads the master brightness
float getBrightness() {
  int power = analogRead(BRIGHTNESS_PIN);
  if (power <= POWER_THRESHOLD_OFF) {
    return 0;
  } else if (power >= POWER_THRESHOLD_MAX) {
    return 1;
  } else {
    //The general color is determin
    return map(power, POWER_THRESHOLD_OFF, POWER_THRESHOLD_MAX, 1, 1000) * 0.001;
  }
}

//Mode SET_COLOR
DualColor handleModeA(DataFieldSet* settings) {
  DualColor ret;
  ret.Boxi1 = convertColor(settings->Pallette[0]);
  ret.Boxi2 = convertColor(settings->Pallette[1]);
  return ret;
}

//Mode FADE_TO_COLOR
DualColor handleModeB(DataFieldSet* settings) {
  DualColor ret;
  float fadeProgress = referenceCounter / (float)settings->Speed;
  ret.Boxi1 = lerpColor(referenceColor.Boxi1, settings->Pallette[0], fadeProgress);
  ret.Boxi2 = lerpColor(referenceColor.Boxi2, settings->Pallette[1], fadeProgress);
  referenceCounter++;
  return ret;
}

//Mode PALLETTE_FADE
DualColor handleModeC(DataFieldSet* settings) {
  DualColor ret;
  Color targetColorA = settings->Pallette[referenceIndex];
  Color targetColorB = settings->Pallette[(referenceIndex + settings->ColorShift) % settings->PalletteSize];

  float fadeProgress = (float)referenceCounter / settings->Speed;
  ret.Boxi1 = lerpColor(referenceColor.Boxi1, targetColorA, fadeProgress);
  ret.Boxi2 = lerpColor(referenceColor.Boxi2, targetColorB, fadeProgress);

  if (referenceCounter++ >= settings->Speed) {
    referenceCounter = 0;
    referenceIndex = (referenceIndex + 1) % settings->PalletteSize;
    updateReferenceColor = true;
  }

  return ret;
}

//Mode PALLETTE_SWITCH
DualColor handleModeD(DataFieldSet* settings, bool onBeat) {
  DualColor ret;

  if (onBeat) {
    referenceIndex = (referenceIndex + 1) % settings->PalletteSize;
  }

  ret.Boxi1 = convertColor(settings->Pallette[referenceIndex]);
  ret.Boxi2 = convertColor(settings->Pallette[(referenceIndex + settings->ColorShift) % settings->PalletteSize]);
  return ret;
}

//Mode PALLETTE_FLASH
DualColor handleModeE(DataFieldSet* settings, bool onBeat) {
  DualColor ret;

  if (onBeat) {
    referenceIndex = (referenceIndex + 1) % settings->PalletteSize;
    referenceCounter = 0;
  }

  float flashBrightness = clampFloat(exp(referenceCounter * settings->Speed * -0.05) * 5 + settings->GeneralPurpose / 255.0);
  ret.Boxi1 = convertColor(settings->Pallette[referenceIndex]);
  ret.Boxi2 = convertColor(settings->Pallette[(referenceIndex + settings->ColorShift) % settings->PalletteSize]);
  return multiplyDualColor(ret, flashBrightness);
}

//Mode BI_COLOR_FLASH
DualColor handleModeF(DataFieldSet* settings, bool onBeat) {
  DualColor ret;

  if (onBeat) {
    referenceIndex = (referenceIndex + 1) % settings->PalletteSize;
    referenceCounter = 0;
  }

  Color startColor = settings->Pallette[referenceIndex];
  Color endColor = settings->Pallette[(referenceIndex + settings->ColorShift) % settings->PalletteSize];

  float flashValue = clampFloat(exp(referenceCounter * settings->Speed * -0.05) * 5 + settings->GeneralPurpose / 255.0);
  FloatColor resultColor = lerpColor(startColor, endColor, flashValue);
  ret.Boxi1 = resultColor;
  ret.Boxi2 = resultColor;
  return ret;
}

//Mode STROBE
DualColor handleModeG(DataFieldSet* settings) {
  DualColor ret;
  int strobeMultiplier = (int)((referenceCounter++ / settings->Speed) % 6) == 0;
  FloatColor resultColor = multiplyColor(settings->Pallette[0], strobeMultiplier);
  ret.Boxi1 = resultColor;
  ret.Boxi2 = resultColor;
  return ret;

}

//Transmits the colors via PWM and DMX
void transmitColors(DualColor outputColor) {
  //Calculate final RGBW values
  uint16_t pwmValR = outputColor.Boxi1.Red * MAX_PWM;
  uint16_t pwmValG = outputColor.Boxi1.Green * MAX_PWM;
  uint16_t pwmValB = outputColor.Boxi1.Blue * MAX_PWM;
  uint16_t pwmValW = outputColor.Boxi1.White * MAX_PWM;
  uint16_t pwmValA = outputColor.Boxi1.Amber * MAX_PWM;
  uint16_t pwmValUV = outputColor.Boxi1.UltraViolet * MAX_PWM;

  //Update internal PWM signals accordingly
  pwm.setPWM(0, 0, pwmValR);
  pwm.setPWM(1, 0, pwmValG);
  pwm.setPWM(2, 0, pwmValB);
  pwm.setPWM(3, 0, pwmValW);
  pwm.setPWM(7, 0, pwmValA);
  pwm.setPWM(5, 0, pwmValUV);

  uint16_t sendValR = outputColor.Boxi1.Red * MAX_PWM;
  uint16_t sendValG = outputColor.Boxi1.Green * MAX_PWM;
  uint16_t sendValB = outputColor.Boxi1.Blue * MAX_PWM;
  uint16_t sendValW = outputColor.Boxi1.White * MAX_PWM;
  uint16_t sendValA = outputColor.Boxi1.Amber * MAX_PWM;
  uint16_t sendValUV = outputColor.Boxi1.UltraViolet * MAX_PWM;
  uint8_t bytesToSend[] = {
    0xe6, 0x21,
    sendValR>>8, sendValR, sendValG>>8, sendValG,
    sendValB>>8, sendValB, sendValW>>8, sendValW,
    sendValA>>8, sendValA, sendValUV>>8, sendValUV};
  Serial.write(bytesToSend, 14);

  //Update external lighting via DMX
  uint8_t dmxValR1 = outputColor.Boxi1.Red * MAX_DMX;
  uint8_t dmxValG1 = outputColor.Boxi1.Green * MAX_DMX;
  uint8_t dmxValB1 = outputColor.Boxi1.Blue * MAX_DMX;
  uint8_t dmxValW1 = outputColor.Boxi1.White * MAX_DMX;
  uint8_t dmxValA1 = outputColor.Boxi1.Amber * MAX_DMX;
  uint8_t dmxValUV1 = outputColor.Boxi1.UltraViolet * MAX_DMX;
  uint8_t dmxValR2 = outputColor.Boxi2.Red * MAX_DMX;
  uint8_t dmxValG2 = outputColor.Boxi2.Green * MAX_DMX;
  uint8_t dmxValB2 = outputColor.Boxi2.Blue * MAX_DMX;
  uint8_t dmxValW2 = outputColor.Boxi2.White * MAX_DMX;
  uint8_t dmxValA2 = outputColor.Boxi2.Amber * MAX_DMX;
  uint8_t dmxValUV2 = outputColor.Boxi2.UltraViolet * MAX_DMX;

  //RGBWAUV
  DmxSimple.write(1, dmxValR1);
  DmxSimple.write(2, dmxValG1);
  DmxSimple.write(3, dmxValB1);
  DmxSimple.write(4, dmxValW1);
  DmxSimple.write(5, dmxValA1);
  DmxSimple.write(6, dmxValUV1);
  DmxSimple.write(7, dmxValR2);
  DmxSimple.write(8, dmxValG2);
  DmxSimple.write(9, dmxValB2);
  DmxSimple.write(10, dmxValW2);
  DmxSimple.write(11, dmxValA2);
  DmxSimple.write(12, dmxValUV2);

  //RGBUV
  DmxSimple.write(13, dmxValR1);
  DmxSimple.write(14, dmxValG1);
  DmxSimple.write(15, dmxValB1);
  DmxSimple.write(16, dmxValUV1);
  DmxSimple.write(17, dmxValR2);
  DmxSimple.write(18, dmxValG2);
  DmxSimple.write(19, dmxValB2);
  DmxSimple.write(20, dmxValUV2);
}

void printSplashScreen() {
  tft.fillScreen(0x0000);
  uint16_t barColors[] = {0xFFFF, 0xC600, 0x0618, 0x0600, 0xC018, 0xC000, 0x0018};

  for (int i=6; i>=0; i-=1) {
    tft.fillRect(0, 160-(i+1)*20, 110, 20, barColors[i]);
  }


  tft.setRotation(3);
  tft.setCursor(0, 111);
  tft.setTextColor(0XFFFF);
  tft.print("Boxi V.3 by       &\nBooting...");

  tft.setCursor(72, 111);
  tft.setTextColor(0X03FE);
  tft.print("Creme");

  tft.setCursor(119, 111);
  tft.setTextColor(0XFD25);
  tft.print("Foxpaw");
}

void setup() {
  // Initialize the screens
  tft.initR(INITR_BLACKTAB);
  // tft.initR(INITR_GREENTAB);      // Init ST7735S chip, green tab
  tft.setSPISpeed(2000000);
  printSplashScreen();

  //Set up DMX
  DmxSimple.usePin(6);
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
  randomSeed(analogRead(7));

  //Send default colors
  createDefaultMode();
  float brightness = getBrightness();
  referenceColor.Boxi1 = {brightness, 0, 0, 0, 0, 0};
  referenceColor.Boxi2 = {brightness, 0, 0, 0, 0, 0};
  transmitColors(referenceColor);

  Serial.begin(9600);
}

void loop() {
  bool onBeat = checkForBeat();

  processUart();
  checkHostActivity();

  if (applyLightingOnNextBeat) {
    applyLighting();
  }

  if (updateReferenceColor) {
    referenceColor = lastOutputColor;
    updateReferenceColor = false;
  }

  DataFieldSet* fieldSet = activeField == 0 ? &fieldSetA : &fieldSetB;

  switch(fieldSet->Mode) {
    case SET_COLOR:
      lastOutputColor = handleModeA(fieldSet);
      break;
    case FADE_TO_COLOR:
      lastOutputColor = handleModeB(fieldSet);
      break;
    case PALLETTE_FADE:
      lastOutputColor = handleModeC(fieldSet);
      break;
    case PALLETTE_SWITCH:
      lastOutputColor = handleModeD(fieldSet, onBeat);
      break;
    case PALLETTE_FLASH:
      lastOutputColor = handleModeE(fieldSet, onBeat);
      break;
    case BI_COLOR_FLASH:
      lastOutputColor = handleModeF(fieldSet, onBeat);
      break;
    case STROBE:
      lastOutputColor = handleModeG(fieldSet);
      break;
    default:
      lastOutputColor = {};
      break;
  }

  DualColor outputColor = multiplyDualColor(lastOutputColor, getBrightness());
  transmitColors(outputColor);
}