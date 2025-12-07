package models

import (
	"time"

	"github.com/google/uuid"
)

type SensorData struct {
	ID          int64      `json:"id" db:"id"`
	Time        time.Time  `json:"time" db:"time"`
	DeviceID    uuid.UUID  `json:"device_id" db:"device_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Temperature *float64   `json:"temperature,omitempty" db:"temperature"`
	Humidity    *float64   `json:"humidity,omitempty" db:"humidity"`
	CO2         *int       `json:"co2,omitempty" db:"co2"`
}

type SensorDataRequest struct {
	Temperature *float64 `json:"temperature,omitempty"`
	Humidity    *float64 `json:"humidity,omitempty"`
	CO2         *int     `json:"co2,omitempty"`
}

type SensorDataResponse struct {
	ID          int64     `json:"id"`
	Time        time.Time `json:"time"`
	Temperature *float64  `json:"temperature,omitempty"`
	Humidity    *float64  `json:"humidity,omitempty"`
	CO2         *int      `json:"co2,omitempty"`
}

type SensorDataQuery struct {
	DeviceID  *uuid.UUID `query:"device_id"`
	StartTime *time.Time `query:"start_time"`
	EndTime   *time.Time `query:"end_time"`
	Limit     int        `query:"limit"`
	Offset    int        `query:"offset"`
}
