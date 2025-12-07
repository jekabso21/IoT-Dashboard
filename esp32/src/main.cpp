#include <Arduino.h>
#include <Wire.h>
#include <U8g2lib.h>
#include "SensorManager.h"
#include "DisplayManager.h"
#include "WiFiManager.h"
#include "ConfigManager.h"
#include "WebServerManager.h"
#include "ApiClient.h"
#include "StateManager.h"

// Hardware components
SensorManager sensor;
DisplayManager displayManager;
U8G2_SSD1306_128X64_NONAME_F_HW_I2C u8g2(U8G2_R0, U8X8_PIN_NONE, 22, 21);

// Network and config components
WiFiManager wifiMgr;
ConfigManager config;
WebServerManager webServer;
ApiClient apiClient;
StateManager stateManager;

// Timing
unsigned long lastSensorRead = 0;
unsigned long lastDataSend = 0;
unsigned long lastStatusChange = 0;
const unsigned long SENSOR_READ_INTERVAL = 5000;
const unsigned long DATA_SEND_INTERVAL = 5000;
const unsigned long STATUS_CHANGE_INTERVAL = 5000;

// Sensor data
int lastCO2 = 0;
float lastTemp = 0;
float lastHumidity = 0;
bool hasValidData = false;
uint8_t currentStatusIndex = 0;

// Server URL - configure this for your deployment
const char* SERVER_URL = "http://192.168.1.100:8080";  // Change to your backend URL

void showMessage(const char* line1, const char* line2 = nullptr) {
  u8g2.clearBuffer();
  u8g2.setFont(u8g2_font_6x10_tr);
  u8g2.drawStr(10, 28, line1);
  if (line2) {
    u8g2.drawStr(10, 42, line2);
  }
  u8g2.sendBuffer();
}

void showAPModeScreen() {
  u8g2.clearBuffer();
  u8g2.setFont(u8g2_font_6x10_tr);
  u8g2.drawStr(5, 12, "WiFi Setup Mode");
  u8g2.drawStr(5, 28, "Connect to:");
  u8g2.drawStr(5, 40, "ESP32_Setup");
  u8g2.drawStr(5, 54, "Open: 192.168.4.1");
  u8g2.sendBuffer();
}

void showPairingScreen() {
  u8g2.clearBuffer();
  u8g2.setFont(u8g2_font_6x10_tr);
  u8g2.drawStr(5, 12, "Device Pairing");
  u8g2.drawStr(5, 28, "Open dashboard to");
  u8g2.drawStr(5, 40, "get pairing code");
  u8g2.drawStr(5, 54, wifiMgr.getIPAddress().c_str());
  u8g2.sendBuffer();
}

void handleWiFiConnect(const String& ssid, const String& password) {
  Serial.printf("WiFi connected callback: %s\n", ssid.c_str());
  config.saveWiFiCredentials(ssid, password);
  stateManager.setState(STATE_NEEDS_PAIRING);
}

bool handlePairing(const String& code, String& error) {
  Serial.printf("Pairing attempt with code: %s\n", code.c_str());
  stateManager.setState(STATE_PAIRING);
  showMessage("Pairing...", "Please wait");

  ApiResponse response = apiClient.pairDevice(code);

  if (response.success) {
    PairingData pairingData;
    if (apiClient.parsePairingResponse(response, pairingData)) {
      // Save tokens to NVS
      config.saveDeviceTokens(
        pairingData.deviceId,
        pairingData.accessToken,
        pairingData.refreshToken,
        millis() / 1000 + pairingData.expiresIn
      );

      apiClient.setAuthToken(pairingData.accessToken);
      webServer.setPairedStatus(true);
      stateManager.setState(STATE_NORMAL_OPERATION);
      stateManager.resetFailureCount();

      Serial.println("Pairing successful!");
      showMessage("Paired!", "Starting...");
      delay(1000);
      return true;
    } else {
      error = "Invalid response";
    }
  } else {
    error = response.errorMessage;
    if (error.isEmpty()) {
      error = "Pairing failed";
    }
  }

  Serial.printf("Pairing failed: %s\n", error.c_str());
  stateManager.incrementFailureCount();
  stateManager.setState(STATE_NEEDS_PAIRING);
  return false;
}

void sendSensorData() {
  if (!hasValidData) return;

  ApiResponse response = apiClient.sendSensorData(lastCO2, lastTemp, lastHumidity);

  if (response.success) {
    stateManager.resetFailureCount();
    stateManager.updateLastActivityTime();
    Serial.println("Data sent successfully");
  } else if (response.httpCode == 401) {
    // Token expired - try refresh
    Serial.println("Token expired, attempting refresh");
    attemptTokenRefresh();
  } else if (response.httpCode == 403) {
    // Device deleted - need to re-pair
    Serial.println("Device unpaired (403), clearing tokens");
    config.clearDeviceTokens();
    webServer.setPairedStatus(false);
    stateManager.setState(STATE_NEEDS_PAIRING);
  } else {
    stateManager.incrementFailureCount();
    Serial.printf("Data send failed: %s (HTTP %d)\n",
                  response.errorMessage.c_str(), response.httpCode);

    if (stateManager.shouldResetToAP()) {
      Serial.println("Too many failures, resetting to AP mode");
      config.clearAll();
      ESP.restart();
    }
  }
}

void attemptTokenRefresh() {
  String deviceId, authToken, refreshToken;
  uint64_t tokenExpiry;

  if (!config.loadDeviceTokens(deviceId, authToken, refreshToken, tokenExpiry)) {
    Serial.println("No refresh token available");
    config.clearDeviceTokens();
    stateManager.setState(STATE_NEEDS_PAIRING);
    return;
  }

  ApiResponse response = apiClient.refreshToken(refreshToken);

  String newToken;
  int expiresIn;
  if (apiClient.parseRefreshResponse(response, newToken, expiresIn)) {
    config.saveDeviceTokens(deviceId, newToken, refreshToken,
                            millis() / 1000 + expiresIn);
    apiClient.setAuthToken(newToken);
    stateManager.resetFailureCount();
    Serial.println("Token refreshed successfully");
  } else {
    Serial.println("Token refresh failed, need to re-pair");
    config.clearDeviceTokens();
    webServer.setPairedStatus(false);
    stateManager.setState(STATE_NEEDS_PAIRING);
  }
}

void setup() {
  Serial.begin(115200);
  Serial.println("\n\nIoT Dashboard Device Starting...");

  // Initialize I2C and display
  Wire.begin(21, 22);
  u8g2.begin();
  showMessage("Initializing...");

  // Initialize components
  config.begin();
  wifiMgr.begin();
  stateManager.begin();
  apiClient.begin(SERVER_URL);
  apiClient.setInsecure(true);  // Skip SSL validation for development

  // Initialize sensor
  sensor.begin(Wire, 0x62);
  delay(500);
  sensor.startMeasurement();

  // Setup web server callbacks
  webServer.onWiFiConnect(handleWiFiConnect);
  webServer.onDevicePair(handlePairing);

  // Check for stored WiFi credentials
  if (config.hasWiFiCredentials()) {
    String ssid, password;
    config.loadWiFiCredentials(ssid, password);

    showMessage("Connecting WiFi...", ssid.c_str());
    stateManager.setState(STATE_CONNECTING_WIFI);

    if (wifiMgr.connect(ssid.c_str(), password.c_str(), 15000)) {
      Serial.printf("WiFi connected: %s\n", wifiMgr.getIPAddress().c_str());

      // Check if we have valid tokens
      if (config.hasValidToken()) {
        String deviceId, authToken, refreshToken;
        uint64_t tokenExpiry;
        config.loadDeviceTokens(deviceId, authToken, refreshToken, tokenExpiry);
        apiClient.setAuthToken(authToken);
        webServer.setPairedStatus(true);
        stateManager.setState(STATE_NORMAL_OPERATION);
        Serial.println("Resuming normal operation");
      } else {
        stateManager.setState(STATE_NEEDS_PAIRING);
        Serial.println("Need pairing");
      }

      // Start web server for status/re-pairing
      webServer.begin();
    } else {
      Serial.println("WiFi connection failed, starting AP mode");
      stateManager.setState(STATE_AP_MODE);
    }
  } else {
    Serial.println("No WiFi credentials, starting AP mode");
    stateManager.setState(STATE_AP_MODE);
  }

  // Start AP mode if needed
  if (stateManager.isInAPMode()) {
    wifiMgr.startAccessPoint("ESP32_Setup");
    webServer.begin();
    showAPModeScreen();
  }

  wifiMgr.enableAutoReconnect(true);
  Serial.printf("State: %s\n", stateManager.getStateName());
}

void loop() {
  // Handle WiFi reconnection
  wifiMgr.handleReconnection();

  // Read sensors periodically
  if (millis() - lastSensorRead >= SENSOR_READ_INTERVAL) {
    if (sensor.read()) {
      lastCO2 = sensor.getCO2();
      lastTemp = sensor.getTemperature();
      lastHumidity = sensor.getHumidity();
      hasValidData = true;
    }
    lastSensorRead = millis();
  }

  // Handle state-specific behavior
  switch (stateManager.getCurrentState()) {
    case STATE_AP_MODE:
      // Display is handled by showAPModeScreen() on state entry
      break;

    case STATE_NEEDS_PAIRING:
    case STATE_PAIRING:
      showPairingScreen();
      break;

    case STATE_NORMAL_OPERATION:
      // Send data periodically
      if (millis() - lastDataSend >= DATA_SEND_INTERVAL) {
        sendSensorData();
        lastDataSend = millis();
      }

      // Update display with sensor data
      if (millis() - lastStatusChange >= STATUS_CHANGE_INTERVAL) {
        currentStatusIndex++;
        lastStatusChange = millis();
      }

      if (hasValidData) {
        bool serverOnline = (stateManager.getTimeSinceLastActivity() < 30000);
        displayManager.showData(u8g2, lastCO2, lastTemp, lastHumidity,
                               serverOnline, currentStatusIndex);
      }
      break;

    case STATE_ERROR:
      showMessage("Error:", stateManager.getLastError().c_str());
      delay(5000);
      stateManager.clearError();
      break;

    default:
      break;
  }

  delay(100);
}
