package models

import (
	"time"

	"github.com/google/uuid"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Device struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	Location    string    `json:"location" db:"location"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	LastSeen    time.Time `json:"last_seen" db:"last_seen"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

type SensorData struct {
	ID          string    `json:"id" db:"id"`
	DeviceID    string    `json:"device_id" db:"device_id"`
	Temperature float64   `json:"temperature" db:"temperature"`
	Humidity    float64   `json:"humidity" db:"humidity"`
	CO2         float64   `json:"co2" db:"co2"`
	Pressure    float64   `json:"pressure" db:"pressure"`
	Light       float64   `json:"light" db:"light"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type DeviceCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type DeviceUpdateRequest struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Location string `json:"location,omitempty"`
	Status   string `json:"status,omitempty"`
}

type SensorDataRequest struct {
	DeviceID    string  `json:"device_id" validate:"required"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	CO2         float64 `json:"co2"`
	Pressure    float64 `json:"pressure"`
	Light       float64 `json:"light"`
}

func NewDevice(name, deviceType, location string) *Device {
	return &Device{
		ID:        uuid.New().String(),
		Name:      name,
		Type:      deviceType,
		Location:  location,
		Status:    "offline",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastSeen:  time.Now(),
		IsActive:  true,
	}
}

func NewSensorData(deviceID string, temp, humidity, co2, pressure, light float64) *SensorData {
	return &SensorData{
		ID:          uuid.New().String(),
		DeviceID:    deviceID,
		Temperature: temp,
		Humidity:    humidity,
		CO2:         co2,
		Pressure:    pressure,
		Light:       light,
		Timestamp:   time.Now(),
	}
}
