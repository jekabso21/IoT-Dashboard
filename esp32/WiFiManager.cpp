#include "WiFiManager.h"

void WiFiManager::begin() {
  WiFi.mode(WIFI_OFF);
  delay(100);
}

void WiFiManager::startAPMode() {
  WiFi.mode(WIFI_AP);
  WiFi.softAPConfig(apIP, apIP, IPAddress(255, 255, 255, 0));
  WiFi.softAP(apSSID, apPassword);
  WiFi.setSleep(false);
  
  state = WIFI_AP_MODE;
  
  Serial.println("=== AP Mode Started ===");
  Serial.print("SSID: ");
  Serial.println(apSSID);
  Serial.print("IP: ");
  Serial.println(WiFi.softAPIP());
  Serial.println("=======================");
}

bool WiFiManager::connectToWiFi(const char* ssid, const char* password) {
  Serial.print("Connecting to: ");
  Serial.println(ssid);
  
  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);
  
  state = WIFI_CONNECTING;
  connectStartTime = millis();
  
  // Wait up to 15 seconds
  int attempts = 0;
  while (WiFi.status() != WL_CONNECTED && attempts < 30) {
    delay(500);
    Serial.print(".");
    attempts++;
  }
  Serial.println();
  
  if (WiFi.status() == WL_CONNECTED) {
    state = WIFI_CONNECTED;
    Serial.println("WiFi Connected!");
    Serial.print("IP: ");
    Serial.println(WiFi.localIP());
    return true;
  } else {
    state = WIFI_FAILED;
    Serial.println("WiFi Connection Failed!");
    return false;
  }
}

void WiFiManager::stopAP() {
  WiFi.softAPdisconnect(true);
  Serial.println("AP Mode stopped");
}

IPAddress WiFiManager::getAPIP() const {
  return apIP;
}

bool WiFiManager::isConnected() const {
  return (state == WIFI_CONNECTED && WiFi.status() == WL_CONNECTED);
}

WiFiState WiFiManager::getState() const {
  return state;
}

int WiFiManager::getStationCount() const {
  return WiFi.softAPgetStationNum();
}