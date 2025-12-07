#include "WebServerManager.h"
#include "web_pages.h"

WebServerManager::WebServerManager()
    : _server(nullptr), _running(false), _isPaired(false), _connectionPending(false) {
}

WebServerManager::~WebServerManager() {
    stop();
}

void WebServerManager::begin() {
    if (_running) {
        return;
    }

    _server = new AsyncWebServer(80);
    setupRoutes();
    _server->begin();
    _running = true;

    Serial.println("Web server started on port 80");
}

void WebServerManager::stop() {
    if (_server) {
        _server->end();
        delete _server;
        _server = nullptr;
    }
    _running = false;
}

void WebServerManager::onWiFiConnect(std::function<void(const String& ssid, const String& password)> callback) {
    _onWiFiConnect = callback;
}

void WebServerManager::onDevicePair(std::function<bool(const String& code, String& error)> callback) {
    _onDevicePair = callback;
}

void WebServerManager::setupRoutes() {
    // Serve main WiFi configuration page
    _server->on("/", HTTP_GET, [](AsyncWebServerRequest* request) {
        request->send_P(200, "text/html", WIFI_PAGE);
    });

    // WiFi network scan endpoint
    _server->on("/scan", HTTP_GET, [](AsyncWebServerRequest* request) {
        Serial.println("Scanning WiFi networks...");

        int numNetworks = WiFi.scanNetworks();

        StaticJsonDocument<2048> doc;
        JsonArray networks = doc.createNestedArray("networks");

        for (int i = 0; i < numNetworks; i++) {
            JsonObject network = networks.createNestedObject();
            network["ssid"] = WiFi.SSID(i);
            network["rssi"] = WiFi.RSSI(i);
            network["encrypted"] = (WiFi.encryptionType(i) != WIFI_AUTH_OPEN);
        }

        String response;
        serializeJson(doc, response);
        request->send(200, "application/json", response);

        WiFi.scanDelete();
        Serial.printf("Found %d networks\n", numNetworks);
    });

    // WiFi connect endpoint
    _server->on("/connect", HTTP_POST, [](AsyncWebServerRequest* request) {
        request->send(400, "application/json", "{\"success\":false,\"message\":\"No body\"}");
    }, nullptr, [this](AsyncWebServerRequest* request, uint8_t* data, size_t len, size_t index, size_t total) {
        StaticJsonDocument<512> doc;
        DeserializationError error = deserializeJson(doc, data, len);

        if (error) {
            request->send(400, "application/json", "{\"success\":false,\"message\":\"Invalid JSON\"}");
            return;
        }

        const char* ssid = doc["ssid"];
        const char* password = doc["password"] | "";

        if (!ssid || strlen(ssid) == 0) {
            request->send(400, "application/json", "{\"success\":false,\"message\":\"SSID required\"}");
            return;
        }

        Serial.printf("Connecting to WiFi: %s\n", ssid);

        // Try to connect
        WiFi.mode(WIFI_AP_STA);
        WiFi.begin(ssid, password);

        int attempts = 0;
        while (WiFi.status() != WL_CONNECTED && attempts < 30) {
            delay(500);
            attempts++;
        }

        if (WiFi.status() == WL_CONNECTED) {
            Serial.printf("Connected! IP: %s\n", WiFi.localIP().toString().c_str());

            StaticJsonDocument<256> resp;
            resp["success"] = true;
            resp["message"] = "Connected";
            resp["ip"] = WiFi.localIP().toString();

            String response;
            serializeJson(resp, response);
            request->send(200, "application/json", response);

            if (_onWiFiConnect) {
                _onWiFiConnect(String(ssid), String(password));
            }
        } else {
            Serial.println("WiFi connection failed");
            request->send(500, "application/json", "{\"success\":false,\"message\":\"Connection failed\"}");
        }
    });

    // Serve pairing page
    _server->on("/pairing", HTTP_GET, [](AsyncWebServerRequest* request) {
        request->send_P(200, "text/html", PAIRING_PAGE);
    });

    // Device pairing endpoint
    _server->on("/pair", HTTP_POST, [](AsyncWebServerRequest* request) {
        request->send(400, "application/json", "{\"success\":false,\"message\":\"No body\"}");
    }, nullptr, [this](AsyncWebServerRequest* request, uint8_t* data, size_t len, size_t index, size_t total) {
        StaticJsonDocument<256> doc;
        DeserializationError error = deserializeJson(doc, data, len);

        if (error) {
            request->send(400, "application/json", "{\"success\":false,\"message\":\"Invalid JSON\"}");
            return;
        }

        const char* code = doc["code"];

        if (!code || strlen(code) != 6) {
            request->send(400, "application/json", "{\"success\":false,\"message\":\"Invalid code format\"}");
            return;
        }

        // Validate digits only
        for (int i = 0; i < 6; i++) {
            if (!isdigit(code[i])) {
                request->send(400, "application/json", "{\"success\":false,\"message\":\"Code must be digits only\"}");
                return;
            }
        }

        Serial.printf("Pairing with code: %s\n", code);

        if (_onDevicePair) {
            String errorMsg;
            bool success = _onDevicePair(String(code), errorMsg);

            StaticJsonDocument<256> resp;
            resp["success"] = success;
            resp["message"] = success ? "Paired successfully" : errorMsg.c_str();

            String response;
            serializeJson(resp, response);
            request->send(success ? 200 : 400, "application/json", response);
        } else {
            request->send(500, "application/json", "{\"success\":false,\"message\":\"Pairing not configured\"}");
        }
    });

    // Device status endpoint
    _server->on("/status", HTTP_GET, [this](AsyncWebServerRequest* request) {
        StaticJsonDocument<512> doc;

        doc["wifi_connected"] = (WiFi.status() == WL_CONNECTED);
        doc["ssid"] = WiFi.SSID();
        doc["ip"] = WiFi.localIP().toString();
        doc["rssi"] = WiFi.RSSI();
        doc["paired"] = _isPaired;

        String response;
        serializeJson(doc, response);
        request->send(200, "application/json", response);
    });

    // 404 handler
    _server->onNotFound([](AsyncWebServerRequest* request) {
        request->send(404, "application/json", "{\"success\":false,\"message\":\"Not found\"}");
    });
}
