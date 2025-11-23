#include <Arduino.h>
#include <Wire.h>
#include <U8g2lib.h>
#include "SensorManager.h"
#include "DisplayManager.h"

SensorManager sensor;
DisplayManager displayManager;
U8G2_SSD1306_128X64_NONAME_F_HW_I2C u8g2(U8G2_R0, U8X8_PIN_NONE, 22, 21);

unsigned long lastSensorRead = 0;
unsigned long lastStatusChange = 0;
const unsigned long SENSOR_READ_INTERVAL = 5000;
const unsigned long STATUS_CHANGE_INTERVAL = 5000;

int lastCO2 = 0;
float lastTemp = 0;
float lastHumidity = 0;
bool hasValidData = false;
uint8_t currentStatusIndex = 0;

void setup() {
  Wire.begin(21, 22);
  u8g2.begin();

  u8g2.clearBuffer();
  u8g2.setFont(u8g2_font_6x10_tr);
  u8g2.drawStr(20, 32, "Initializing...");
  u8g2.sendBuffer();

  sensor.begin(Wire, 0x62);
  delay(1000);
  sensor.startMeasurement();

  delay(1000);
}

void loop() {
  if (millis() - lastSensorRead >= SENSOR_READ_INTERVAL) {
    if (sensor.read()) {
      lastCO2 = sensor.getCO2();
      lastTemp = sensor.getTemperature();
      lastHumidity = sensor.getHumidity();
      hasValidData = true;
    } else {
      if (!hasValidData) {
        u8g2.clearBuffer();
        u8g2.setFont(u8g2_font_6x10_tr);
        u8g2.drawStr(10, 32, "Waiting for");
        u8g2.drawStr(10, 44, "sensor data...");
        u8g2.sendBuffer();
      }
    }
    lastSensorRead = millis();
  }

  if (millis() - lastStatusChange >= STATUS_CHANGE_INTERVAL) {
    currentStatusIndex++;
    lastStatusChange = millis();
  }

  if (hasValidData) {
    displayManager.showData(u8g2, lastCO2, lastTemp, lastHumidity, false, currentStatusIndex);
  }

  delay(100);
}
