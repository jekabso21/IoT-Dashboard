#ifndef STATEMANAGER_H
#define STATEMANAGER_H

#include <Arduino.h>

// Device operation states
enum DeviceState {
  STATE_INIT,                // Initial startup state
  STATE_AP_MODE,             // Access Point mode (no WiFi credentials)
  STATE_CONNECTING_WIFI,     // Attempting to connect to WiFi
  STATE_WIFI_CONNECTED,      // WiFi connected but no device token
  STATE_NEEDS_PAIRING,       // Waiting for pairing code
  STATE_PAIRING,             // Currently pairing with server
  STATE_NORMAL_OPERATION,    // Normal operation (sending sensor data)
  STATE_ERROR                // Error state
};

class StateManager {
public:
  StateManager();

  // Initialize state manager
  void begin();

  // State management
  DeviceState getCurrentState() const;
  void setState(DeviceState newState);
  const char* getStateName() const;
  const char* getStateName(DeviceState state) const;

  // State queries
  bool shouldShowPairingScreen() const;
  bool shouldSendData() const;
  bool isInAPMode() const;
  bool isConnectedToWiFi() const;
  bool isPaired() const;

  // Failure tracking
  void incrementFailureCount();
  void resetFailureCount();
  uint8_t getFailureCount() const;
  bool shouldResetToAP() const;

  // Error management
  void setError(const String& errorMessage);
  String getLastError() const;
  void clearError();

  // Timing helpers
  void updateLastActivityTime();
  unsigned long getTimeSinceLastActivity() const;

private:
  DeviceState _currentState;
  DeviceState _previousState;
  uint8_t _failureCount;
  String _lastError;
  unsigned long _lastActivityTime;
  unsigned long _stateEntryTime;

  static const uint8_t MAX_FAILURES = 10;
};

#endif
