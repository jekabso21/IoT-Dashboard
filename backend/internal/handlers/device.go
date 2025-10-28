package handlers

import (
	"net/http"
	"strconv"

	"github.com/jekabso21/IoT-Dashboard/backend/internal/logger"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type DeviceHandler struct {
	// TODO: Add database service when implemented
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

func (h *DeviceHandler) GetDevices(c echo.Context) error {
	logger.Info("Getting all devices")

	// TODO: Implement database query
	devices := []models.Device{
		{
			ID:       "1",
			Name:     "Temperature Sensor 1",
			Type:     "temperature",
			Location: "Living Room",
			Status:   "online",
			IsActive: true,
		},
		{
			ID:       "2",
			Name:     "Humidity Sensor 1",
			Type:     "humidity",
			Location: "Bedroom",
			Status:   "offline",
			IsActive: true,
		},
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Devices retrieved successfully",
		Data:    devices,
	})
}

func (h *DeviceHandler) GetDevice(c echo.Context) error {
	deviceID := c.Param("id")
	logger.WithField("device_id", deviceID).Info("Getting device")

	// TODO: Implement database query
	device := models.Device{
		ID:       deviceID,
		Name:     "Temperature Sensor 1",
		Type:     "temperature",
		Location: "Living Room",
		Status:   "online",
		IsActive: true,
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device retrieved successfully",
		Data:    device,
	})
}

func (h *DeviceHandler) CreateDevice(c echo.Context) error {
	var req models.DeviceCreateRequest
	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind device creation request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	logger.WithField("device_name", req.Name).Info("Creating device")

	// TODO: Implement database insert
	device := models.NewDevice(req.Name, req.Type, req.Location)

	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Device created successfully",
		Data:    device,
	})
}

func (h *DeviceHandler) UpdateDevice(c echo.Context) error {
	deviceID := c.Param("id")
	var req models.DeviceUpdateRequest

	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind device update request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	logger.WithField("device_id", deviceID).Info("Updating device")

	// TODO: Implement database update
	device := models.Device{
		ID:       deviceID,
		Name:     req.Name,
		Type:     req.Type,
		Location: req.Location,
		Status:   req.Status,
		IsActive: true,
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device updated successfully",
		Data:    device,
	})
}

func (h *DeviceHandler) DeleteDevice(c echo.Context) error {
	deviceID := c.Param("id")
	logger.WithField("device_id", deviceID).Info("Deleting device")

	// TODO: Implement database delete

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device deleted successfully",
	})
}

func (h *DeviceHandler) GetDeviceData(c echo.Context) error {
	deviceID := c.Param("id")
	limitStr := c.QueryParam("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	logger.WithFields(map[string]interface{}{
		"device_id": deviceID,
		"limit":     limit,
	}).Info("Getting device data")

	// TODO: Implement database query
	sensorData := []models.SensorData{
		{
			ID:          "1",
			DeviceID:    deviceID,
			Temperature: 22.5,
			Humidity:    45.2,
			CO2:         400,
			Pressure:    1013.25,
			Light:       850,
		},
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device data retrieved successfully",
		Data:    sensorData,
	})
}
