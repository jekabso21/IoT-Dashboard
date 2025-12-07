package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/repository"
	"github.com/labstack/echo/v4"
)

func IngestSensorData(c echo.Context) error {
	deviceID, ok := c.Get("device_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	var req models.SensorDataRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	sensorData, err := repository.InsertSensorData(c.Request().Context(), deviceID, userID, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to insert sensor data",
		})
	}

	if err := repository.UpdateDeviceLastSeen(c.Request().Context(), deviceID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update device last seen",
		})
	}

	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data: models.SensorDataResponse{
			ID:          sensorData.ID,
			Time:        sensorData.Time,
			Temperature: sensorData.Temperature,
			Humidity:    sensorData.Humidity,
			CO2:         sensorData.CO2,
		},
	})
}

func GetSensorData(c echo.Context) error {
	// Verify user is authenticated (userID is used for RLS via middleware)
	_, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	var query models.SensorDataQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid query parameters",
		})
	}

	if query.Limit == 0 {
		query.Limit = 100
	}

	sensorData, err := repository.GetSensorData(c.Request().Context(), &query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve sensor data",
		})
	}

	responses := make([]models.SensorDataResponse, len(sensorData))
	for i, data := range sensorData {
		responses[i] = models.SensorDataResponse{
			ID:          data.ID,
			Time:        data.Time,
			Temperature: data.Temperature,
			Humidity:    data.Humidity,
			CO2:         data.CO2,
		}
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    responses,
	})
}
