#ifndef SENSORMANAGER_H
#define SENSORMANAGER_H

#include <Wire.h>
#include <SparkFun_SCD4x_Arduino_Library.h>

class SensorManager {
public:
  void begin(TwoWire& wirePort, uint8_t address);
  void startMeasurement();
  bool read();
  uint16_t getCO2() const;
  float getTemperature() const;
  float getHumidity() const;

private:
  SCD4x scd4x;
  uint16_t co2;
  float temperature;
  float humidity;
};

#endif
