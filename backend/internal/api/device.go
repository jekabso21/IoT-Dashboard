package api

import (
	"encoding/json"
	"net/http"
)

// SensorData represents a simple structure for the IoT data
type SensorData struct {
	DeviceID    string  `json:"device_id"`
	Temperature float64 `json:"temperature"`
	CO2         float64 `json:"co2"`
	Humidity    float64 `json:"humidity"`
}

// DeviceDataHandler handles POST requests with sensor data
func DeviceDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data SensorData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Here you would typically save data to the database

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received"))
}
