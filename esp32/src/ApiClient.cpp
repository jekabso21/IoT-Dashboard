#include "ApiClient.h"

ApiClient::ApiClient()
  : _requestTimeout(10000), _insecure(true) {
}

void ApiClient::begin(const String& serverUrl) {
  _serverUrl = serverUrl;
  if (_serverUrl.endsWith("/")) {
    _serverUrl.remove(_serverUrl.length() - 1);
  }
}

void ApiClient::setAuthToken(const String& token) {
  _authToken = token;
}

void ApiClient::setRequestTimeout(int timeoutMs) {
  _requestTimeout = timeoutMs;
}

void ApiClient::setInsecure(bool insecure) {
  _insecure = insecure;
}

String ApiClient::buildUrl(const String& endpoint) {
  String url = _serverUrl;
  if (!endpoint.startsWith("/")) {
    url += "/";
  }
  url += endpoint;
  return url;
}

ApiResponse ApiClient::makeRequest(const String& method, const String& endpoint,
                                    const String& body, bool requiresAuth) {
  ApiResponse response;
  response.success = false;
  response.httpCode = 0;

  HTTPClient http;
  String url = buildUrl(endpoint);

  // Check if HTTPS
  if (url.startsWith("https://")) {
    WiFiClientSecure* client = new WiFiClientSecure();
    if (_insecure) {
      client->setInsecure();
    }
    if (!http.begin(*client, url)) {
      response.errorMessage = "Failed to initialize HTTPS client";
      delete client;
      return response;
    }
  } else {
    WiFiClient* client = new WiFiClient();
    if (!http.begin(*client, url)) {
      response.errorMessage = "Failed to initialize HTTP client";
      delete client;
      return response;
    }
  }

  http.setTimeout(_requestTimeout);
  http.addHeader("Content-Type", "application/json");

  if (requiresAuth && _authToken.length() > 0) {
    http.addHeader("Authorization", "Bearer " + _authToken);
  }

  int httpCode;
  if (method == "GET") {
    httpCode = http.GET();
  } else if (method == "POST") {
    httpCode = http.POST(body);
  } else if (method == "PUT") {
    httpCode = http.PUT(body);
  } else {
    response.errorMessage = "Unsupported HTTP method";
    http.end();
    return response;
  }

  response.httpCode = httpCode;

  if (httpCode > 0) {
    String payload = http.getString();

    DeserializationError error = deserializeJson(response.data, payload);

    if (error) {
      response.errorMessage = "JSON parse error: " + String(error.c_str());
      response.success = false;
    } else {
      if (httpCode >= 200 && httpCode < 300) {
        response.success = response.data["success"] | true;
      } else {
        if (response.data.containsKey("error")) {
          response.errorMessage = response.data["error"].as<String>();
        } else if (response.data.containsKey("message")) {
          response.errorMessage = response.data["message"].as<String>();
        } else {
          response.errorMessage = "HTTP error: " + String(httpCode);
        }
      }
    }
  } else {
    response.errorMessage = "Request failed: " + http.errorToString(httpCode);
  }

  http.end();
  return response;
}

ApiResponse ApiClient::pairDevice(const String& pairingCode) {
  JsonDocument doc;
  doc["code"] = pairingCode;

  String body;
  serializeJson(doc, body);

  Serial.printf("Pairing with code: %s\n", pairingCode.c_str());
  return makeRequest("POST", "/api/devices/pair", body, false);
}

ApiResponse ApiClient::refreshToken(const String& refreshToken) {
  JsonDocument doc;
  doc["refresh_token"] = refreshToken;

  String body;
  serializeJson(doc, body);

  return makeRequest("POST", "/api/devices/token/refresh", body, true);
}

ApiResponse ApiClient::sendSensorData(uint16_t co2, float temperature, float humidity) {
  JsonDocument doc;
  doc["co2"] = co2;
  doc["temperature"] = temperature;
  doc["humidity"] = humidity;

  String body;
  serializeJson(doc, body);

  return makeRequest("POST", "/api/devices/data", body, true);
}

bool ApiClient::parsePairingResponse(const ApiResponse& response, PairingData& pairingData) {
  if (!response.success) {
    return false;
  }

  // The response is wrapped in "data" field based on APIResponse model
  JsonVariantConst data = response.data["data"];
  if (data.isNull()) {
    // Try direct access if not wrapped
    data = response.data.as<JsonVariantConst>();
  }

  if (!data.containsKey("device_id") || !data.containsKey("access_token")) {
    return false;
  }

  pairingData.deviceId = data["device_id"].as<String>();
  pairingData.accessToken = data["access_token"].as<String>();
  pairingData.refreshToken = data["refresh_token"].as<String>();
  pairingData.expiresIn = data["expires_in"] | 86400;

  return true;
}

bool ApiClient::parseRefreshResponse(const ApiResponse& response, String& newAccessToken, int& expiresIn) {
  if (!response.success) {
    return false;
  }

  JsonVariantConst data = response.data["data"];
  if (data.isNull()) {
    data = response.data.as<JsonVariantConst>();
  }

  if (!data.containsKey("access_token")) {
    return false;
  }

  newAccessToken = data["access_token"].as<String>();
  expiresIn = data["expires_in"] | 86400;
  return true;
}
