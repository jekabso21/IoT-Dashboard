#include "StateManager.h"

StateManager::StateManager()
  : _currentState(STATE_INIT),
    _previousState(STATE_INIT),
    _failureCount(0),
    _lastActivityTime(0),
    _stateEntryTime(0) {
}

void StateManager::begin() {
  _currentState = STATE_INIT;
  _previousState = STATE_INIT;
  _failureCount = 0;
  _lastActivityTime = millis();
  _stateEntryTime = millis();
  _lastError = "";
}

DeviceState StateManager::getCurrentState() const {
  return _currentState;
}

void StateManager::setState(DeviceState newState) {
  if (_currentState != newState) {
    Serial.printf("State change: %s -> %s\n", getStateName(_currentState), getStateName(newState));
    _previousState = _currentState;
    _currentState = newState;
    _stateEntryTime = millis();

    if (_previousState == STATE_ERROR) {
      _lastError = "";
    }

    if (newState == STATE_NORMAL_OPERATION || newState == STATE_WIFI_CONNECTED) {
      _failureCount = 0;
    }
  }
}

const char* StateManager::getStateName() const {
  return getStateName(_currentState);
}

const char* StateManager::getStateName(DeviceState state) const {
  switch (state) {
    case STATE_INIT:
      return "Initializing";
    case STATE_AP_MODE:
      return "AP Mode";
    case STATE_CONNECTING_WIFI:
      return "Connecting WiFi";
    case STATE_WIFI_CONNECTED:
      return "WiFi Connected";
    case STATE_NEEDS_PAIRING:
      return "Needs Pairing";
    case STATE_PAIRING:
      return "Pairing";
    case STATE_NORMAL_OPERATION:
      return "Normal";
    case STATE_ERROR:
      return "Error";
    default:
      return "Unknown";
  }
}

bool StateManager::shouldShowPairingScreen() const {
  return _currentState == STATE_NEEDS_PAIRING ||
         _currentState == STATE_PAIRING;
}

bool StateManager::shouldSendData() const {
  return _currentState == STATE_NORMAL_OPERATION;
}

bool StateManager::isInAPMode() const {
  return _currentState == STATE_AP_MODE;
}

bool StateManager::isConnectedToWiFi() const {
  return _currentState == STATE_WIFI_CONNECTED ||
         _currentState == STATE_NEEDS_PAIRING ||
         _currentState == STATE_PAIRING ||
         _currentState == STATE_NORMAL_OPERATION;
}

bool StateManager::isPaired() const {
  return _currentState == STATE_NORMAL_OPERATION;
}

void StateManager::incrementFailureCount() {
  _failureCount++;
  Serial.printf("Failure count: %d/%d\n", _failureCount, MAX_FAILURES);
}

void StateManager::resetFailureCount() {
  if (_failureCount > 0) {
    Serial.println("Failure count reset");
    _failureCount = 0;
  }
}

uint8_t StateManager::getFailureCount() const {
  return _failureCount;
}

bool StateManager::shouldResetToAP() const {
  return _failureCount >= MAX_FAILURES;
}

void StateManager::setError(const String& errorMessage) {
  _lastError = errorMessage;
  Serial.printf("Error: %s\n", errorMessage.c_str());
  setState(STATE_ERROR);
}

String StateManager::getLastError() const {
  return _lastError;
}

void StateManager::clearError() {
  _lastError = "";
  if (_currentState == STATE_ERROR) {
    setState(_previousState);
  }
}

void StateManager::updateLastActivityTime() {
  _lastActivityTime = millis();
}

unsigned long StateManager::getTimeSinceLastActivity() const {
  return millis() - _lastActivityTime;
}
