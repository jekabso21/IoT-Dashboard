#include "DisplayManager.h"
#include <stdio.h>

const char* DisplayManager::getCO2Status(int co2) {
  if (co2 < 400) return "CO2: Low";
  if (co2 <= 800) return "CO2: Good";
  if (co2 <= 1000) return "CO2: Fair";
  if (co2 <= 1500) return "CO2: Poor";
  return "CO2: Bad";
}

const char* DisplayManager::getTempStatus(float temp) {
  if (temp < 18) return "Temp: Cold";
  if (temp <= 22) return "Temp: Cool";
  if (temp <= 26) return "Temp: Good";
  if (temp <= 30) return "Temp: Warm";
  return "Temp: Hot";
}

const char* DisplayManager::getHumidityStatus(float humidity) {
  if (humidity < 30) return "Humid: Dry";
  if (humidity <= 50) return "Humid: Good";
  if (humidity <= 70) return "Humid: High";
  return "Humid: VeryHigh";
}

void DisplayManager::showData(U8G2_SSD1306_128X64_NONAME_F_HW_I2C& u8g2, int co2, float temp, float humidity, bool serverConnected, uint8_t statusIndex) {
  char buffer[20];

  u8g2.clearBuffer();
  u8g2.setFont(u8g2_font_6x10_tr);

  const char* statusMsg;
  switch(statusIndex % 3) {
    case 0:
      statusMsg = getCO2Status(co2);
      break;
    case 1:
      statusMsg = getTempStatus(temp);
      break;
    case 2:
      statusMsg = getHumidityStatus(humidity);
      break;
  }

  u8g2.drawFrame(0, 0, 128, 12);
  int msgWidth = u8g2.getStrWidth(statusMsg);
  u8g2.drawStr((128 - msgWidth) / 2, 9, statusMsg);

  u8g2.drawFrame(0, 14, 128, 16);
  snprintf(buffer, sizeof(buffer), "CO2: %d ppm", co2);
  u8g2.drawStr(4, 25, buffer);

  u8g2.drawFrame(0, 32, 64, 16);
  snprintf(buffer, sizeof(buffer), "T: %.1fC", temp);
  u8g2.drawStr(4, 43, buffer);

  u8g2.drawFrame(64, 32, 64, 16);
  snprintf(buffer, sizeof(buffer), "H: %.1f%%", humidity);
  u8g2.drawStr(68, 43, buffer);

  u8g2.drawFrame(0, 50, 128, 14);
  if (serverConnected) {
    u8g2.drawStr(4, 60, "Server: Online");
  } else {
    u8g2.drawStr(4, 60, "Server: Offline");
  }

  u8g2.sendBuffer();
}

void DisplayManager::showError(U8G2_SSD1306_128X64_NONAME_F_HW_I2C& u8g2, const char* message) {
  u8g2.clearBuffer();
  u8g2.setFont(u8g2_font_6x10_tr);

  u8g2.drawFrame(5, 10, 118, 44);
  u8g2.drawFrame(6, 11, 116, 42);

  u8g2.setFont(u8g2_font_10x20_tr);
  u8g2.drawStr(30, 28, "ERROR");

  u8g2.setFont(u8g2_font_6x10_tr);
  u8g2.drawStr(12, 45, message);

  u8g2.sendBuffer();
}
