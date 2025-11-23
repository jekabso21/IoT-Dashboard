#include "SensorManager.h"

void SensorManager::begin(TwoWire& wirePort, uint8_t address) {
  scd4x.begin(wirePort, address);
}

void SensorManager::startMeasurement() {
  scd4x.startPeriodicMeasurement();
}

bool SensorManager::read() {
  if (!scd4x.readMeasurement()) {
    return false;
  }

  uint16_t newCO2 = scd4x.getCO2();
  if (newCO2 > 0) {
    co2 = newCO2;
    temperature = scd4x.getTemperature();
    humidity = scd4x.getHumidity();
    return true;
  }

  return false;
}

uint16_t SensorManager::getCO2() const {
  return co2;
}

float SensorManager::getTemperature() const {
  return temperature;
}

float SensorManager::getHumidity() const {
  return humidity;
}
