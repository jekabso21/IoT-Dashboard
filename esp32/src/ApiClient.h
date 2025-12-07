#ifndef APICLIENT_H
#define APICLIENT_H

#include <Arduino.h>
#include <WiFiClientSecure.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>

// API Response structure
struct ApiResponse {
  bool success;
  int httpCode;
  String errorMessage;
  JsonDocument data;
};

// Pairing response data
struct PairingData {
  String deviceId;
  String accessToken;
  String refreshToken;
  int expiresIn;
};

class ApiClient {
public:
  ApiClient();

  // Initialize with server URL
  void begin(const String& serverUrl);

  // Set current auth token for authenticated requests
  void setAuthToken(const String& token);

  // API Methods
  ApiResponse pairDevice(const String& pairingCode);
  ApiResponse refreshToken(const String& refreshToken);
  ApiResponse sendSensorData(uint16_t co2, float temperature, float humidity);

  // Helper to parse pairing response
  bool parsePairingResponse(const ApiResponse& response, PairingData& pairingData);

  // Helper to parse refresh token response
  bool parseRefreshResponse(const ApiResponse& response, String& newAccessToken, int& expiresIn);

  // Configuration
  void setRequestTimeout(int timeoutMs);
  void setInsecure(bool insecure);

private:
  String _serverUrl;
  String _authToken;
  int _requestTimeout;
  bool _insecure;

  // Internal HTTP request handler
  ApiResponse makeRequest(const String& method, const String& endpoint,
                          const String& body = "", bool requiresAuth = false);

  String buildUrl(const String& endpoint);
};

#endif
