#ifndef WEBSERVERMANAGER_H
#define WEBSERVERMANAGER_H
#include <WebServer.h>

class WebServerManager {
public:
  void begin(const char* pin);
  void handleClient();
  void stop();
  bool isConfigured() const { return configured; }
  bool isServerConfigured() const { return serverConfigured; }
  String getConfiguredSSID() const { return configSSID; }
  String getConfiguredPassword() const { return configPassword; }
  String getServerURL() const { return serverURL; }
  String getAuthCode() const { return authCode; }
  void setWiFiConnected(bool connected) { wifiConnected = connected; }
  
private:
  WebServer server{80};
  String accessPin;
  bool pinVerified = false;
  bool configured = false;
  bool wifiConnected = false;
  bool serverConfigured = false;
  String configSSID;
  String configPassword;
  String serverURL;
  String authCode;
  
  static WebServerManager* instance;
  
  static void handleRoot();
  static void handleVerifyPin();
  static void handleScan();
  static void handleConnect();
  static void handleServerConfig();
  static void handleNotFound();
};

#endif