#include <Wire.h>
#include <Adafruit_SSD1306.h>
#include <HTTPClient.h>
#include "WiFiManager.h"
#include "WebServerManager.h"
#include "SensorManager.h"
#include "DisplayManager.h"

#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64
#define OLED_RESET -1

Adafruit_SSD1306 display(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, OLED_RESET);
SensorManager sensor;
DisplayManager displayManager;
WiFiManager wifi;
WebServerManager webServer;

// Setup state machine
enum SetupState {
  STATE_AP_WAITING,      // Waiting for user to connect
  STATE_AP_CONNECTED,    // User connected, show IP and PIN
  STATE_WIFI_CONNECTING, // Connecting to home WiFi
  STATE_WIFI_CONNECTED,  // WiFi connected, waiting for server config
  STATE_SERVER_CONFIG,   // Configuring server settings
  STATE_NORMAL_OPERATION // Connected, send data every 5 sec
};

SetupState currentState = STATE_AP_WAITING;
String accessPin = "";
String serverURL = "";
String authCode = "";
unsigned long lastSensorRead = 0;
unsigned long wifiConnectedTime = 0;
const unsigned long WIFI_STABLE_TIME = 10000; // 10 seconds

void setup() {
  Serial.begin(115200);
  Wire.begin(21, 22);
  
  // Initialize display
  if (!display.begin(SSD1306_SWITCHCAPVCC, 0x3C)) {
    Serial.println("SSD1306 allocation failed");
    while (true);
  }
  display.clearDisplay();
  display.display();
  
  // Initialize sensor
  sensor.begin(Wire, 0x62);
  sensor.startMeasurement();
  
  // Generate random 4-digit PIN
  randomSeed(analogRead(0));
  accessPin = String(random(1000, 10000));
  
  // Start in AP mode
  wifi.begin();
  wifi.startAPMode();
  displayManager.showWaitingForConnection(display, "ESP32_Setup");
  
  currentState = STATE_AP_WAITING;
  Serial.println("=== Setup Started ===");
  Serial.print("PIN: ");
  Serial.println(accessPin);
}

void sendDataToServer(int co2, float temp, float humidity) {
  if (serverURL.length() == 0) {
    Serial.println("Server URL not configured");
    return;
  }
  
  HTTPClient http;
  http.begin(serverURL);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("Authorization", authCode);
  
  String jsonData = "{\"co2\":" + String(co2) + 
                    ",\"temperature\":" + String(temp, 1) + 
                    ",\"humidity\":" + String(humidity, 1) + "}";
  
  Serial.print("Sending data to server: ");
  Serial.println(jsonData);
  
  int httpCode = http.POST(jsonData);
  
  if (httpCode > 0) {
    Serial.print("Server response code: ");
    Serial.println(httpCode);
    if (httpCode == 200) {
      Serial.println("Data sent successfully!");
    }
  } else {
    Serial.print("HTTP POST failed: ");
    Serial.println(http.errorToString(httpCode));
  }
  
  http.end();
}

void loop() {
  switch (currentState) {
    
    case STATE_AP_WAITING: {
      // Check if someone connected to AP
      if (wifi.getStationCount() > 0) {
        Serial.println("User connected to AP");
        webServer.begin(accessPin.c_str());
        displayManager.showAccessInfo(display, "192.168.4.1", accessPin.c_str());
        currentState = STATE_AP_CONNECTED;
      }
      delay(500);
      break;
    }
    
    case STATE_AP_CONNECTED: {
      webServer.handleClient();
      
      // Check if WiFi has been configured
      if (webServer.isConfigured()) {
        String ssid = webServer.getConfiguredSSID();
        String password = webServer.getConfiguredPassword();
        
        Serial.println("WiFi credentials received");
        displayManager.showConnecting(display, ssid.c_str());
        
        delay(1000);
        
        currentState = STATE_WIFI_CONNECTING;
        
        // Try to connect to home WiFi
        if (wifi.connectToWiFi(ssid.c_str(), password.c_str())) {
          wifiConnectedTime = millis();
          webServer.setWiFiConnected(true);
          currentState = STATE_WIFI_CONNECTED;
          
          display.clearDisplay();
          display.setTextSize(1);
          display.setTextColor(SSD1306_WHITE);
          display.setCursor(0, 20);
          display.println("WiFi Connected!");
          display.println();
          display.println("Configure server");
          display.println("settings in browser");
          display.display();
        } else {
          displayManager.showError(display, "WiFi Failed!\nRestarting...");
          delay(5000);
          ESP.restart();
        }
      }
      delay(100);
      break;
    }
    
    case STATE_WIFI_CONNECTING: {
      // This state is handled by the blocking connectToWiFi call
      // We transition from STATE_AP_CONNECTED directly to STATE_WIFI_CONNECTED
      break;
    }
    
    case STATE_WIFI_CONNECTED: {
      webServer.handleClient();
      
      // Wait for stable connection (10 seconds)
      if (millis() - wifiConnectedTime >= WIFI_STABLE_TIME) {
        // Check if server config has been entered
        if (webServer.isServerConfigured()) {
          serverURL = webServer.getServerURL();
          authCode = webServer.getAuthCode();
          
          Serial.println("=== Configuration Complete ===");
          Serial.print("Server URL: ");
          Serial.println(serverURL);
          Serial.print("Auth Code: ");
          Serial.println(authCode);
          
          // Stop AP and web server
          wifi.stopAP();
          webServer.stop();
          
          currentState = STATE_NORMAL_OPERATION;
          lastSensorRead = millis();
          
          displayManager.showData(display, 0, 0, 0, true);
          Serial.println("Entering normal operation mode");
        }
      }
      
      delay(100);
      break;
    }
    
    case STATE_NORMAL_OPERATION: {
      // Check WiFi connection
      if (!wifi.isConnected()) {
        displayManager.showError(display, "WiFi Disconnected!");
        delay(3000);
        ESP.restart();
      }
      
      // Read sensor every 5 seconds
      if (millis() - lastSensorRead >= 5000) {
        if (sensor.read()) {
          int co2 = sensor.getCO2();
          float temp = sensor.getTemperature();
          float humidity = sensor.getHumidity();
          
          // Display data with WiFi icon
          displayManager.showData(display, co2, temp, humidity, true);
          
          // Send data to server
          Serial.println("=== Sensor Data ===");
          Serial.print("CO2: ");
          Serial.print(co2);
          Serial.println(" ppm");
          Serial.print("Temperature: ");
          Serial.print(temp);
          Serial.println(" C");
          Serial.print("Humidity: ");
          Serial.print(humidity);
          Serial.println(" %");
          Serial.println("==================");
          
          sendDataToServer(co2, temp, humidity);
          
        } else {
          Serial.println("Failed to read sensor");
        }
        lastSensorRead = millis();
      }
      
      delay(100);
      break;
    }
  }
}