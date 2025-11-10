#include "WebServerManager.h"
#include <WiFi.h>

WebServerManager* WebServerManager::instance = nullptr;

const char pinPage[] PROGMEM = R"rawliteral(
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>ESP32 Setup</title>
<style>
body{font-family:Arial,sans-serif;text-align:center;padding:50px;background:#667eea;color:#fff;margin:0}
.box{background:rgba(255,255,255,0.1);padding:40px;border-radius:15px;max-width:400px;margin:0 auto}
input{width:100%;padding:15px;margin:10px 0;border:none;border-radius:8px;font-size:18px;text-align:center;box-sizing:border-box}
button{width:100%;padding:15px;background:#4CAF50;color:#fff;border:none;border-radius:8px;font-size:18px;cursor:pointer}
button:hover{background:#45a049}
.error{color:#ff6b6b;margin:10px 0}
</style>
</head>
<body>
<div class="box">
<h1>ESP32 Setup</h1>
<p>Enter PIN shown on display</p>
<form action="/verify" method="POST">
<input type="text" name="pin" placeholder="Enter PIN" maxlength="4" autofocus>
<button type="submit">Verify</button>
</form>
<div class="error" id="error"></div>
</div>
<script>
const params=new URLSearchParams(window.location.search);
if(params.get('error')=='1')document.getElementById('error').textContent='Invalid PIN!';
</script>
</body>
</html>
)rawliteral";

const char wifiPage[] PROGMEM = R"rawliteral(
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>WiFi Setup</title>
<style>
body{font-family:Arial,sans-serif;padding:20px;background:#667eea;color:#fff;margin:0}
.box{background:rgba(255,255,255,0.1);padding:30px;border-radius:15px;max-width:500px;margin:0 auto}
.network{background:rgba(255,255,255,0.2);padding:15px;margin:10px 0;border-radius:8px;cursor:pointer;display:flex;justify-content:space-between;align-items:center}
.network:hover{background:rgba(255,255,255,0.3)}
input,select{width:100%;padding:12px;margin:10px 0;border:none;border-radius:8px;font-size:16px;box-sizing:border-box}
button{width:100%;padding:15px;background:#4CAF50;color:#fff;border:none;border-radius:8px;font-size:18px;cursor:pointer;margin-top:10px}
button:hover{background:#45a049}
.scanning{text-align:center;padding:20px}
.signal{font-size:20px}
</style>
</head>
<body>
<div class="box">
<h1>WiFi Networks</h1>
<div id="networks" class="scanning">Scanning...</div>
<form id="wifiForm" style="display:none" action="/connect" method="POST">
<h3>Connect to: <span id="selectedSSID"></span></h3>
<input type="hidden" name="ssid" id="ssidInput">
<input type="password" name="password" placeholder="WiFi Password" required>
<button type="submit">Connect</button>
</form>
</div>
<script>
function scanNetworks(){
  fetch('/scan').then(r=>r.json()).then(data=>{
    const div=document.getElementById('networks');
    if(data.networks.length===0){
      div.innerHTML='<p>No networks found</p>';
      return;
    }
    div.innerHTML='';
    data.networks.forEach(net=>{
      const signal=net.rssi>-50?'[***]':net.rssi>-70?'[** ]':'[*  ]';
      const netDiv=document.createElement('div');
      netDiv.className='network';
      netDiv.innerHTML='<span>'+net.ssid+'</span><span class="signal">'+signal+'</span>';
      netDiv.onclick=()=>{
        document.getElementById('selectedSSID').textContent=net.ssid;
        document.getElementById('ssidInput').value=net.ssid;
        document.getElementById('wifiForm').style.display='block';
      };
      div.appendChild(netDiv);
    });
  });
}
scanNetworks();
</script>
</body>
</html>
)rawliteral";

const char connectingPage[] PROGMEM = R"rawliteral(
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Connecting...</title>
<style>
body{font-family:Arial,sans-serif;text-align:center;padding:50px;background:#667eea;color:#fff;margin:0}
.box{background:rgba(255,255,255,0.1);padding:40px;border-radius:15px;max-width:400px;margin:0 auto}
.spinner{border:8px solid rgba(255,255,255,0.3);border-top:8px solid #fff;border-radius:50%;width:60px;height:60px;animation:spin 1s linear infinite;margin:20px auto}
@keyframes spin{0%{transform:rotate(0deg)}100%{transform:rotate(360deg)}}
</style>
</head>
<body>
<div class="box">
<h1>Connecting...</h1>
<div class="spinner"></div>
<p>ESP32 is connecting to WiFi.</p>
<p>Please wait...</p>
</div>
<script>
setTimeout(function(){window.location.href='/';},5000);
</script>
</body>
</html>
)rawliteral";

const char serverConfigPage[] PROGMEM = R"rawliteral(
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Server Setup</title>
<style>
body{font-family:Arial,sans-serif;text-align:center;padding:50px;background:#667eea;color:#fff;margin:0}
.box{background:rgba(255,255,255,0.1);padding:40px;border-radius:15px;max-width:400px;margin:0 auto}
input{width:100%;padding:15px;margin:10px 0;border:none;border-radius:8px;font-size:16px;box-sizing:border-box}
button{width:100%;padding:15px;background:#4CAF50;color:#fff;border:none;border-radius:8px;font-size:18px;cursor:pointer;margin-top:10px}
button:hover{background:#45a049}
.success{color:#4CAF50;font-size:24px;margin:20px 0}
label{display:block;text-align:left;margin-top:15px;font-size:14px}
</style>
</head>
<body>
<div class="box">
<div class="success">WiFi Connected!</div>
<h2>Server Configuration</h2>
<form action="/serverconfig" method="POST">
<label>Server URL:</label>
<input type="text" name="url" placeholder="http://example.com/api" required>
<label>Authentication Code:</label>
<input type="text" name="code" placeholder="Enter auth code" required>
<button type="submit">Save & Start</button>
</form>
</div>
</body>
</html>
)rawliteral";

const char completePage[] PROGMEM = R"rawliteral(
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Setup Complete</title>
<style>
body{font-family:Arial,sans-serif;text-align:center;padding:50px;background:#667eea;color:#fff;margin:0}
.box{background:rgba(255,255,255,0.1);padding:40px;border-radius:15px;max-width:400px;margin:0 auto}
.success{color:#4CAF50;font-size:48px;margin:20px 0}
</style>
</head>
<body>
<div class="box">
<div class="success">&#10003;</div>
<h1>Setup Complete!</h1>
<p>ESP32 is now configured and running.</p>
<p>This page will close in 5 seconds.</p>
</div>
<script>
setTimeout(function(){window.close();},5000);
</script>
</body>
</html>
)rawliteral";

void WebServerManager::begin(const char* pin) {
  instance = this;
  accessPin = String(pin);
  pinVerified = false;
  configured = false;
  wifiConnected = false;
  serverConfigured = false;
  
  server.on("/", handleRoot);
  server.on("/verify", HTTP_POST, handleVerifyPin);
  server.on("/scan", HTTP_GET, handleScan);
  server.on("/connect", HTTP_POST, handleConnect);
  server.on("/serverconfig", HTTP_POST, handleServerConfig);
  server.onNotFound(handleNotFound);
  
  server.begin();
  Serial.println("Web server started");
}

void WebServerManager::handleClient() {
  server.handleClient();
}

void WebServerManager::stop() {
  server.stop();
  Serial.println("Web server stopped");
}

void WebServerManager::handleRoot() {
  if (!instance->pinVerified) {
    instance->server.sendHeader("Content-Type", "text/html; charset=utf-8");
    instance->server.send_P(200, "text/html", pinPage);
  } else if (!instance->wifiConnected) {
    instance->server.sendHeader("Content-Type", "text/html; charset=utf-8");
    instance->server.send_P(200, "text/html", wifiPage);
  } else {
    instance->server.sendHeader("Content-Type", "text/html; charset=utf-8");
    instance->server.send_P(200, "text/html", serverConfigPage);
  }
}

void WebServerManager::handleVerifyPin() {
  String enteredPin = instance->server.arg("pin");
  
  if (enteredPin == instance->accessPin) {
    instance->pinVerified = true;
    Serial.println("PIN verified successfully");
    instance->server.sendHeader("Location", "/", true);
    instance->server.send(302, "text/plain", "");
  } else {
    Serial.println("Invalid PIN entered");
    instance->server.sendHeader("Location", "/?error=1", true);
    instance->server.send(302, "text/plain", "");
  }
}

void WebServerManager::handleScan() {
  Serial.println("Scanning WiFi networks...");
  int n = WiFi.scanNetworks();
  
  String json = "{\"networks\":[";
  bool first = true;
  
  // Keep track of SSIDs we've already added to avoid duplicates
  for (int i = 0; i < n; i++) {
    String ssid = WiFi.SSID(i);
    
    // Skip empty SSIDs
    if (ssid.length() == 0) continue;
    
    // Check if this SSID already exists in our list
    bool isDuplicate = false;
    for (int j = 0; j < i; j++) {
      if (WiFi.SSID(j) == ssid) {
        isDuplicate = true;
        break;
      }
    }
    
    // Skip duplicates
    if (isDuplicate) continue;
    
    if (!first) json += ",";
    json += "{\"ssid\":\"" + ssid + "\",\"rssi\":" + String(WiFi.RSSI(i)) + "}";
    first = false;
  }
  json += "]}";
  
  instance->server.send(200, "application/json", json);
  Serial.print("Found unique networks: ");
  Serial.println(n);
}

void WebServerManager::handleConnect() {
  instance->configSSID = instance->server.arg("ssid");
  instance->configPassword = instance->server.arg("password");
  instance->configured = true;
  
  Serial.print("WiFi config received - SSID: ");
  Serial.println(instance->configSSID);
  
  instance->server.sendHeader("Content-Type", "text/html; charset=utf-8");
  instance->server.send_P(200, "text/html", connectingPage);
}

void WebServerManager::handleNotFound() {
  instance->server.sendHeader("Location", "/", true);
  instance->server.send(302, "text/plain", "");
}

void WebServerManager::handleServerConfig() {
  instance->serverURL = instance->server.arg("url");
  instance->authCode = instance->server.arg("code");
  instance->serverConfigured = true;
  
  Serial.println("Server config received:");
  Serial.print("URL: ");
  Serial.println(instance->serverURL);
  Serial.print("Auth Code: ");
  Serial.println(instance->authCode);
  
  instance->server.sendHeader("Content-Type", "text/html; charset=utf-8");
  instance->server.send_P(200, "text/html", completePage);
}