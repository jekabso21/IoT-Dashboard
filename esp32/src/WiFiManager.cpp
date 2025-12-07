#include "WiFiManager.h"

void WiFiManager::begin() {
  WiFi.mode(WIFI_MODE_NULL);
  WiFi.disconnect(true);
  delay(100);
}

bool WiFiManager::startAccessPoint(const char* ssid, const char* password) {
  WiFi.mode(WIFI_AP);

  bool success;
  if (password == nullptr || strlen(password) == 0) {
    success = WiFi.softAP(ssid);
  } else {
    success = WiFi.softAP(ssid, password);
  }

  if (success) {
    delay(100);
    return true;
  }

  return false;
}

void WiFiManager::stopAccessPoint() {
  WiFi.softAPdisconnect(true);
  WiFi.mode(WIFI_OFF);
  delay(100);
}

bool WiFiManager::isAccessPointActive() {
  return WiFi.getMode() == WIFI_AP || WiFi.getMode() == WIFI_AP_STA;
}

String WiFiManager::getAccessPointIP() {
  if (isAccessPointActive()) {
    return WiFi.softAPIP().toString();
  }
  return "";
}

std::vector<WiFiNetwork> WiFiManager::scanNetworks() {
  std::vector<WiFiNetwork> networks;

  // Ensure WiFi is in station mode for scanning
  WiFi.mode(WIFI_STA);
  WiFi.disconnect();
  delay(100);

  int numNetworks = WiFi.scanNetworks();

  if (numNetworks > 0) {
    for (int i = 0; i < numNetworks; i++) {
      WiFiNetwork network;
      network.ssid = WiFi.SSID(i);
      network.rssi = WiFi.RSSI(i);
      network.encrypted = (WiFi.encryptionType(i) != WIFI_AUTH_OPEN);

      // Only add networks with valid SSIDs
      if (network.ssid.length() > 0) {
        networks.push_back(network);
      }
    }
  }

  WiFi.scanDelete();
  return networks;
}

bool WiFiManager::connect(const char* ssid, const char* password, uint32_t timeout_ms) {
  // Store credentials for auto-reconnect
  lastSSID = String(ssid);
  lastPassword = String(password);

  // Disconnect any existing connection
  WiFi.disconnect(true);
  delay(100);

  // Set to station mode
  WiFi.mode(WIFI_STA);

  // Begin connection
  if (password == nullptr || strlen(password) == 0) {
    WiFi.begin(ssid);
  } else {
    WiFi.begin(ssid, password);
  }

  // Wait for connection with timeout
  unsigned long startTime = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - startTime < timeout_ms) {
    delay(100);
  }

  if (WiFi.status() == WL_CONNECTED) {
    return true;
  }

  // Connection failed
  WiFi.disconnect(true);
  return false;
}

void WiFiManager::disconnect() {
  WiFi.disconnect(true);
  WiFi.mode(WIFI_OFF);
  delay(100);
}

bool WiFiManager::isConnected() {
  return WiFi.status() == WL_CONNECTED;
}

String WiFiManager::getIPAddress() {
  if (isConnected()) {
    return WiFi.localIP().toString();
  }
  return "";
}

int32_t WiFiManager::getRSSI() {
  if (isConnected()) {
    return WiFi.RSSI();
  }
  return 0;
}

void WiFiManager::enableAutoReconnect(bool enable) {
  autoReconnectEnabled = enable;
  WiFi.setAutoReconnect(enable);
}

void WiFiManager::handleReconnection() {
  if (!autoReconnectEnabled) {
    return;
  }

  // Only attempt reconnect if we're not connected and have credentials
  if (isConnected() || lastSSID.length() == 0) {
    return;
  }

  // Throttle reconnection attempts
  unsigned long currentTime = millis();
  if (currentTime - lastReconnectAttempt < RECONNECT_INTERVAL) {
    return;
  }

  lastReconnectAttempt = currentTime;

  // Attempt to reconnect
  WiFi.disconnect();
  delay(100);
  WiFi.begin(lastSSID.c_str(), lastPassword.c_str());
}
