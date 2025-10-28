package handlers

import (
	"net/http"
	"strings"

	"github.com/jekabso21/IoT-Dashboard/backend/internal/logger"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type IoTHandler struct {
	// TODO: Add database service when implemented
}

func NewIoTHandler() *IoTHandler {
	return &IoTHandler{}
}

func (h *IoTHandler) SendSensorData(c echo.Context) error {
	var req models.SensorDataRequest
	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind sensor data request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if strings.TrimSpace(req.DeviceID) == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "device_id is required",
		})
	}

	logger.WithFields(logrus.Fields{
		"device_id":   req.DeviceID,
		"temperature": req.Temperature,
		"humidity":    req.Humidity,
		"co2":         req.CO2,
	}).Info("Received sensor data from IoT device")

	// TODO: Implement database insert
	sensorData := models.NewSensorData(
		req.DeviceID,
		req.Temperature,
		req.Humidity,
		req.CO2,
		req.Pressure,
		req.Light,
	)

	// TODO: Update device last_seen timestamp
	// TODO: Trigger real-time notifications if needed

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Sensor data received successfully",
		Data:    sensorData,
	})
}

func (h *IoTHandler) GetDeviceStatus(c echo.Context) error {
	deviceID := c.Param("id")
	logger.WithField("device_id", deviceID).Info("Getting device status")

	// TODO: Implement database query
	status := map[string]interface{}{
		"device_id": deviceID,
		"status":    "online",
		"last_seen": "2024-01-01T12:00:00Z",
		"battery":   85,
		"signal":    "strong",
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device status retrieved successfully",
		Data:    status,
	})
}

func (h *IoTHandler) UpdateDeviceStatus(c echo.Context) error {
	deviceID := c.Param("id")

	var req struct {
		Status string `json:"status" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind device status update request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if strings.TrimSpace(req.Status) == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "status is required",
		})
	}

	logger.WithFields(logrus.Fields{
		"device_id": deviceID,
		"status":    req.Status,
	}).Info("Updating device status")

	// TODO: Implement database update

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device status updated successfully",
	})
}

func (h *IoTHandler) GetDeviceCommands(c echo.Context) error {
	deviceID := c.Param("id")
	logger.WithField("device_id", deviceID).Info("Getting device commands")

	// TODO: Implement database query for pending commands
	commands := []map[string]interface{}{
		{
			"id":      "1",
			"command": "restart",
			"params":  map[string]interface{}{},
		},
		{
			"id":      "2",
			"command": "update_config",
			"params": map[string]interface{}{
				"sampling_rate": 30,
			},
		},
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device commands retrieved successfully",
		Data:    commands,
	})
}

func (h *IoTHandler) SendCommandResponse(c echo.Context) error {
	deviceID := c.Param("id")
	commandID := c.Param("commandId")

	var req struct {
		Success bool        `json:"success"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind command response request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	logger.WithFields(logrus.Fields{
		"device_id":  deviceID,
		"command_id": commandID,
		"success":    req.Success,
	}).Info("Received command response from device")

	// TODO: Implement database update for command status

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Command response received successfully",
	})
}
