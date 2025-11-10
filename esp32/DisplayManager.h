#ifndef DISPLAYMANAGER_H
#define DISPLAYMANAGER_H
#include <Adafruit_SSD1306.h>

class DisplayManager {
public:
  void showWaitingForConnection(Adafruit_SSD1306& display, const char* ssid);
  void showAccessInfo(Adafruit_SSD1306& display, const char* ip, const char* pin);
  void showConnecting(Adafruit_SSD1306& display, const char* ssid);
  void showData(Adafruit_SSD1306& display, int co2, float temp, float humidity, bool wifiConnected);
  void showError(Adafruit_SSD1306& display, const char* message);
};

#endif