package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/auth"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/repository"
	"github.com/labstack/echo/v4"
)

func GeneratePairingCode(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	pairingCode, err := repository.GeneratePairingCode(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate pairing code",
		})
	}

	expiresIn := int(pairingCode.ExpiresAt.Sub(pairingCode.CreatedAt).Seconds())

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.GeneratePairingCodeResponse{
			Code:      pairingCode.Code,
			ExpiresAt: pairingCode.ExpiresAt,
			ExpiresIn: expiresIn,
		},
	})
}

func PairDevice(c echo.Context) error {
	var req models.DevicePairRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	pairingCode, err := repository.ValidatePairingCode(c.Request().Context(), req.Code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid or expired pairing code",
		})
	}

	device, err := repository.CreateDevice(c.Request().Context(), pairingCode.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create device",
		})
	}

	if err := repository.MarkPairingCodeAsUsed(c.Request().Context(), req.Code); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to mark pairing code as used",
		})
	}

	accessToken, refreshToken, expiresIn, err := auth.GenerateDeviceToken(device.ID, pairingCode.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate device tokens",
		})
	}

	if err := repository.StoreDeviceTokens(c.Request().Context(), device.ID, pairingCode.UserID, accessToken, refreshToken, expiresIn); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to store device tokens",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.DevicePairResponse{
			DeviceID:     device.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expiresIn,
		},
	})
}

func RefreshDeviceToken(c echo.Context) error {
	var req models.DeviceTokenRefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request",
		})
	}

	deviceID, ok := c.Get("device_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	valid, err := repository.ValidateDeviceRefreshToken(c.Request().Context(), deviceID, req.RefreshToken)
	if err != nil || !valid {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid refresh token",
		})
	}

	device, err := repository.GetDeviceByID(c.Request().Context(), deviceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Device not found",
		})
	}

	if device.UserID == nil {
		return c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Device not paired",
		})
	}

	accessToken, _, expiresIn, err := auth.GenerateDeviceToken(deviceID, *device.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate access token",
		})
	}

	if err := repository.StoreDeviceTokens(c.Request().Context(), deviceID, *device.UserID, accessToken, req.RefreshToken, expiresIn); err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to store device tokens",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.DeviceTokenRefreshResponse{
			AccessToken: accessToken,
			ExpiresIn:   expiresIn,
		},
	})
}
