#ifndef CONFIGMANAGER_H
#define CONFIGMANAGER_H

#include <Preferences.h>
#include <Arduino.h>

class ConfigManager {
public:
  // Initialize NVS
  bool begin();

  // WiFi Credentials
  bool saveWiFiCredentials(const String& ssid, const String& password);
  bool loadWiFiCredentials(String& ssid, String& password);
  bool hasWiFiCredentials();

  // Device Authentication Tokens
  bool saveDeviceTokens(const String& deviceId, const String& authToken,
                        const String& refreshToken, uint64_t tokenExpiry);
  bool loadDeviceTokens(String& deviceId, String& authToken,
                        String& refreshToken, uint64_t& tokenExpiry);
  bool hasValidToken();

  // Server Configuration
  bool saveServerURL(const String& url);
  bool loadServerURL(String& url);
  bool hasServerURL();

  // Clear stored data
  void clearWiFiCredentials();
  void clearDeviceTokens();
  void clearAll();

private:
  Preferences preferences;
  const char* NAMESPACE = "iot_config";

  // Keys for NVS storage
  const char* KEY_WIFI_SSID = "wifi_ssid";
  const char* KEY_WIFI_PASS = "wifi_pass";
  const char* KEY_DEVICE_ID = "device_id";
  const char* KEY_AUTH_TOKEN = "auth_token";
  const char* KEY_REFRESH_TOKEN = "refresh_token";
  const char* KEY_TOKEN_EXPIRY = "token_expiry";
  const char* KEY_SERVER_URL = "server_url";
};

#endif
