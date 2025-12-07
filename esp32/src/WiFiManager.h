#ifndef WIFIMANAGER_H
#define WIFIMANAGER_H

#include <WiFi.h>
#include <vector>

struct WiFiNetwork {
  String ssid;
  int32_t rssi;
  bool encrypted;
};

class WiFiManager {
public:
  // Initialize WiFi subsystem
  void begin();

  // AP Mode: Create access point for device setup
  bool startAccessPoint(const char* ssid = "ESP32_Setup", const char* password = nullptr);
  void stopAccessPoint();
  bool isAccessPointActive();
  String getAccessPointIP();

  // WiFi Scanning: Discover available networks
  std::vector<WiFiNetwork> scanNetworks();

  // Station Mode: Connect to WiFi network
  bool connect(const char* ssid, const char* password, uint32_t timeout_ms = 10000);
  void disconnect();
  bool isConnected();
  String getIPAddress();
  int32_t getRSSI();

  // Auto-reconnection
  void enableAutoReconnect(bool enable = true);
  void handleReconnection();

private:
  bool autoReconnectEnabled = false;
  unsigned long lastReconnectAttempt = 0;
  const unsigned long RECONNECT_INTERVAL = 30000; // 30 seconds between attempts
  String lastSSID = "";
  String lastPassword = "";
};

#endif
