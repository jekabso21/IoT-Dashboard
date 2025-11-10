#include "DisplayManager.h"

void DisplayManager::showWaitingForConnection(Adafruit_SSD1306& display, const char* ssid) {
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SSD1306_WHITE);
  display.cp437(true); // Enable Code Page 437 for better character support
  
  display.setCursor(0, 5);
  display.println("ESP32 Setup Mode");
  display.println();
  display.println("Connect to WiFi:");
  display.println();
  display.println(ssid);
  
  display.display();
}

void DisplayManager::showAccessInfo(Adafruit_SSD1306& display, const char* ip, const char* pin) {
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SSD1306_WHITE);
  display.cp437(true);
  
  display.setCursor(0, 0);
  display.println("User Connected!");
  display.println();
  display.println("Open browser:");
  display.println(ip);
  display.println();
  display.println("Enter PIN:");
  display.setCursor(40, 48);
  display.setTextSize(2);
  display.println(pin);
  
  display.display();
}

void DisplayManager::showConnecting(Adafruit_SSD1306& display, const char* ssid) {
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SSD1306_WHITE);
  display.cp437(true);
  
  display.setCursor(0, 15);
  display.println("Connecting to:");
  display.println();
  display.println(ssid);
  display.println();
  display.println("Please wait...");
  
  display.display();
}

void DisplayManager::showData(Adafruit_SSD1306& display, int co2, float temp, float humidity, bool wifiConnected) {
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SSD1306_WHITE);
  display.cp437(true);
  
  // WiFi icon in top right corner
  if (wifiConnected) {
    display.setCursor(100, 0);
    display.print("[WiFi]");
  }
  
  display.setCursor(0, 0);
  display.println("Sensor Data:");
  display.println();
  
  display.print("CO2: ");
  display.print(co2);
  display.println(" ppm");
  display.println();
  
  display.print("Temp: ");
  display.print(temp, 1);
  display.println(" C");
  display.println();
  
  display.print("Humidity: ");
  display.print(humidity, 1);
  display.println(" %");
  
  display.display();
}

void DisplayManager::showError(Adafruit_SSD1306& display, const char* message) {
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SSD1306_WHITE);
  display.cp437(true);
  
  display.setCursor(0, 15);
  display.println("ERROR:");
  display.println();
  display.println(message);
  
  display.display();
}