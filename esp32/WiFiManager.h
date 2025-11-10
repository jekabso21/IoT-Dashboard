#ifndef WIFIMANAGER_H
#define WIFIMANAGER_H
#include <WiFi.h>

enum WiFiState {
  WIFI_AP_MODE,
  WIFI_CONNECTING,
  WIFI_CONNECTED,
  WIFI_FAILED
};

class WiFiManager {
public:
  void begin();
  void startAPMode();
  bool connectToWiFi(const char* ssid, const char* password);
  void stopAP();
  IPAddress getAPIP() const;
  bool isConnected() const;
  WiFiState getState() const;
  int getStationCount() const;
  
private:
  const char* apSSID = "ESP32_Setup";
  const char* apPassword = "";  // Open network for easier setup
  IPAddress apIP = IPAddress(192, 168, 4, 1);
  WiFiState state = WIFI_AP_MODE;
  unsigned long connectStartTime = 0;
};

#endif