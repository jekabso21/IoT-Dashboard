#ifndef WEBSERVERMANAGER_H
#define WEBSERVERMANAGER_H

#include <Arduino.h>
#include <WiFi.h>
#include <ESPAsyncWebServer.h>
#include <ArduinoJson.h>
#include <functional>

class WebServerManager {
public:
    WebServerManager();
    ~WebServerManager();

    // Initialize the web server
    void begin();

    // Stop the web server
    void stop();

    // Check if server is running
    bool isRunning() const { return _running; }

    // Set callback for WiFi connection
    void onWiFiConnect(std::function<void(const String& ssid, const String& password)> callback);

    // Set callback for device pairing
    void onDevicePair(std::function<bool(const String& code, String& error)> callback);

    // Set paired status for status endpoint
    void setPairedStatus(bool paired) { _isPaired = paired; }

private:
    AsyncWebServer* _server;
    bool _running;
    bool _isPaired;

    // Callbacks
    std::function<void(const String& ssid, const String& password)> _onWiFiConnect;
    std::function<bool(const String& code, String& error)> _onDevicePair;

    // Pending connection data
    String _pendingSSID;
    String _pendingPassword;
    bool _connectionPending;

    // Setup all routes
    void setupRoutes();
};

#endif
