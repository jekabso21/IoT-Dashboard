#include "ConfigManager.h"

bool ConfigManager::begin() {
  // Preferences library handles initialization internally
  return true;
}

bool ConfigManager::saveWiFiCredentials(const String& ssid, const String& password) {
  if (!preferences.begin(NAMESPACE, false)) {
    return false;
  }

  bool success = true;

  if (preferences.putString(KEY_WIFI_SSID, ssid) == 0) {
    success = false;
  }

  if (preferences.putString(KEY_WIFI_PASS, password) == 0) {
    success = false;
  }

  preferences.end();
  return success;
}

bool ConfigManager::loadWiFiCredentials(String& ssid, String& password) {
  if (!preferences.begin(NAMESPACE, true)) {
    return false;
  }

  ssid = preferences.getString(KEY_WIFI_SSID, "");
  password = preferences.getString(KEY_WIFI_PASS, "");

  preferences.end();

  return (ssid.length() > 0);
}

bool ConfigManager::hasWiFiCredentials() {
  if (!preferences.begin(NAMESPACE, true)) {
    return false;
  }

  String ssid = preferences.getString(KEY_WIFI_SSID, "");
  preferences.end();

  return (ssid.length() > 0);
}

bool ConfigManager::saveDeviceTokens(const String& deviceId, const String& authToken,
                                      const String& refreshToken, uint64_t tokenExpiry) {
  if (!preferences.begin(NAMESPACE, false)) {
    return false;
  }

  bool success = true;

  if (preferences.putString(KEY_DEVICE_ID, deviceId) == 0) {
    success = false;
  }

  if (preferences.putString(KEY_AUTH_TOKEN, authToken) == 0) {
    success = false;
  }

  if (preferences.putString(KEY_REFRESH_TOKEN, refreshToken) == 0) {
    success = false;
  }

  if (preferences.putULong64(KEY_TOKEN_EXPIRY, tokenExpiry) == 0) {
    success = false;
  }

  preferences.end();
  return success;
}

bool ConfigManager::loadDeviceTokens(String& deviceId, String& authToken,
                                      String& refreshToken, uint64_t& tokenExpiry) {
  if (!preferences.begin(NAMESPACE, true)) {
    return false;
  }

  deviceId = preferences.getString(KEY_DEVICE_ID, "");
  authToken = preferences.getString(KEY_AUTH_TOKEN, "");
  refreshToken = preferences.getString(KEY_REFRESH_TOKEN, "");
  tokenExpiry = preferences.getULong64(KEY_TOKEN_EXPIRY, 0);

  preferences.end();

  return (deviceId.length() > 0 && authToken.length() > 0);
}

bool ConfigManager::hasValidToken() {
  if (!preferences.begin(NAMESPACE, true)) {
    return false;
  }

  String deviceId = preferences.getString(KEY_DEVICE_ID, "");
  String authToken = preferences.getString(KEY_AUTH_TOKEN, "");
  uint64_t tokenExpiry = preferences.getULong64(KEY_TOKEN_EXPIRY, 0);

  preferences.end();

  // Check if we have tokens and they haven't expired
  if (deviceId.length() == 0 || authToken.length() == 0) {
    return false;
  }

  // Get current Unix timestamp (seconds since epoch)
  // Note: ESP32 time needs to be set via NTP or external source
  uint64_t currentTime = (uint64_t)(millis() / 1000);

  // If tokenExpiry is 0, we can't determine validity, so return true if tokens exist
  if (tokenExpiry == 0) {
    return true;
  }

  return (currentTime < tokenExpiry);
}

bool ConfigManager::saveServerURL(const String& url) {
  if (!preferences.begin(NAMESPACE, false)) {
    return false;
  }

  bool success = (preferences.putString(KEY_SERVER_URL, url) > 0);

  preferences.end();
  return success;
}

bool ConfigManager::loadServerURL(String& url) {
  if (!preferences.begin(NAMESPACE, true)) {
    return false;
  }

  url = preferences.getString(KEY_SERVER_URL, "");

  preferences.end();

  return (url.length() > 0);
}

bool ConfigManager::hasServerURL() {
  if (!preferences.begin(NAMESPACE, true)) {
    return false;
  }

  String url = preferences.getString(KEY_SERVER_URL, "");
  preferences.end();

  return (url.length() > 0);
}

void ConfigManager::clearWiFiCredentials() {
  if (!preferences.begin(NAMESPACE, false)) {
    return;
  }

  preferences.remove(KEY_WIFI_SSID);
  preferences.remove(KEY_WIFI_PASS);

  preferences.end();
}

void ConfigManager::clearDeviceTokens() {
  if (!preferences.begin(NAMESPACE, false)) {
    return;
  }

  preferences.remove(KEY_DEVICE_ID);
  preferences.remove(KEY_AUTH_TOKEN);
  preferences.remove(KEY_REFRESH_TOKEN);
  preferences.remove(KEY_TOKEN_EXPIRY);

  preferences.end();
}

void ConfigManager::clearAll() {
  if (!preferences.begin(NAMESPACE, false)) {
    return;
  }

  preferences.clear();
  preferences.end();
}
