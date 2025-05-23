#include <Wire.h>
#include <Adafruit_PWMServoDriver.h>
#include <Adafruit_GFX.h>    // Core graphics library
#include <Adafruit_ST7735.h> // Hardware-specific library for ST7735
#include <SPI.h>

#define TFT_CS         10
#define TFT_RST       -1
#define TFT_DC         9

enum DisplayStatusCode {
  BOOTING = 0x00,
  HOST_AWAKE = 0x01,
  HOST_NO_ACTIVITY = 0x02,
  DSP_SERVER_FAILED = 0x03,
  HOST_CONNECTION_FAILED = 0x04,
  ACTIVE = 0x05
}; 

struct Color {
  uint16_t Red;
  uint16_t Green;
  uint16_t Blue;
  uint16_t White;
  uint16_t Amber;
  uint16_t UltraViolet;
}; 

const uint16_t HOST_CONNECTION_TIMEOUT = 1000;

//Variables
uint16_t hostConnectionCounter = HOST_CONNECTION_TIMEOUT;
Color currentOutput;
DisplayStatusCode dspStatusCode = BOOTING;
Adafruit_ST7735 tft = Adafruit_ST7735(TFT_CS, TFT_DC, TFT_RST);
Adafruit_PWMServoDriver pwm = Adafruit_PWMServoDriver();


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

void checkHostActivity() {
  if (hostConnectionCounter > 0) {
    hostConnectionCounter--;
    
  } else if (hostConnectionCounter == 0 && dspStatusCode == BOOTING) {
    handleDisplayStatusCode(HOST_NO_ACTIVITY);
  }
}

//Processes the incoming data from the RPI at the UART port
bool processUart() {
  //Make sure header is fine before handling data
  if (Serial.available() < 14 ||
      Serial.read() != 0xe6 ||
      Serial.read() != 0x21) return false;

  //Read payload
  uint8_t receivedData[12];
  if (Serial.readBytes(receivedData, 12) != 12) {
    return false;
  }

  //Convert to uint16
  uint16_t convertedData[6];
  for(int i = 0; i < 6; i++) {
    convertedData[i] = ((uint16_t)receivedData[i*2] << 8) + receivedData[i*2 + 1];
  }

  currentOutput.Red = convertedData[0];
  currentOutput.Green = convertedData[1];
  currentOutput.Blue = convertedData[2];
  currentOutput.White = convertedData[3];
  currentOutput.Amber = convertedData[4];
  currentOutput.UltraViolet = convertedData[5];
  return true;
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

void setup() {
  pwm.begin();
  pwm.setOscillatorFrequency(27000000);
  pwm.setPWMFreq(1600);  // This is the maximum PWM frequency

  // if you want to really speed stuff up, you can go into 'fast 400khz I2C' mode
  // some i2c devices dont like this so much so if you're sharing the bus, watch
  // out for this!
  Wire.setClock(400000);

  // Initialize the screens
  tft.initR(INITR_BLACKTAB);
  // tft.initR(INITR_GREENTAB);      // Init ST7735S chip, green tab
  tft.setSPISpeed(2000000);
  printSplashScreen();

  Serial.begin(19200);
}

void loop() {
  checkHostActivity();

  processUart();

  //Update internal PWM signals accordingly
  pwm.setPWM(0, 0, currentOutput.Red);
  pwm.setPWM(1, 0, currentOutput.Green);
  pwm.setPWM(2, 0, currentOutput.Blue);
  pwm.setPWM(3, 0, currentOutput.White);
  pwm.setPWM(7, 0, currentOutput.Amber);
  pwm.setPWM(5, 0, currentOutput.UltraViolet);
}
