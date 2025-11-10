#include "SensorManager.h"

void SensorManager::begin(TwoWire& wirePort, uint8_t address) {
  scd4x.begin(wirePort, address);
}

void SensorManager::startMeasurement() {
  scd4x.startPeriodicMeasurement();
}

bool SensorManager::read() {
  return scd4x.readMeasurement(co2, temperature, humidity) == 0;
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
