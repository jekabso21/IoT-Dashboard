#ifndef DISPLAYMANAGER_H
#define DISPLAYMANAGER_H
#include <U8g2lib.h>

class DisplayManager {
public:
  void showData(U8G2_SSD1306_128X64_NONAME_F_HW_I2C& u8g2, int co2, float temp, float humidity, bool serverConnected, uint8_t statusIndex);
  void showError(U8G2_SSD1306_128X64_NONAME_F_HW_I2C& u8g2, const char* message);

private:
  const char* getCO2Status(int co2);
  const char* getTempStatus(float temp);
  const char* getHumidityStatus(float humidity);
};

#endif
