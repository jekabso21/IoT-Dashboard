package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/repository"
	"github.com/labstack/echo/v4"
)

func GetUserDevices(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	devices, err := repository.GetUserDevices(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve devices",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    devices,
	})
}

func GetDevice(c echo.Context) error {
	deviceIDStr := c.Param("id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid device ID",
		})
	}

	device, err := repository.GetDeviceByID(c.Request().Context(), deviceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Device not found",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    device,
	})
}

func UpdateDevice(c echo.Context) error {
	deviceIDStr := c.Param("id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid device ID",
		})
	}

	var req models.DeviceUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	if err := repository.UpdateDevice(c.Request().Context(), deviceID, &req); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update device",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device updated successfully",
	})
}

func DeleteDevice(c echo.Context) error {
	deviceIDStr := c.Param("id")
	deviceID, err := uuid.Parse(deviceIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid device ID",
		})
	}

	if err := repository.DeleteDevice(c.Request().Context(), deviceID); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete device",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Device deleted successfully",
	})
}
