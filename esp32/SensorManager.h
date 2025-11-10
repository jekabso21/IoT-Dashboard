#ifndef SENSORMANAGER_H
#define SENSORMANAGER_H

#include <Wire.h>
#include <SensirionI2cScd4x.h>

class SensorManager {
public:
  void begin(TwoWire& wirePort, uint8_t address);
  void startMeasurement();
  bool read();
  uint16_t getCO2() const;
  float getTemperature() const;
  float getHumidity() const;

private:
  SensirionI2cScd4x scd4x;
  uint16_t co2;
  float temperature;
  float humidity;
};

#endif
