package handlers

import (
	"net/http"

	"github.com/jekabso21/IoT-Dashboard/backend/internal/logger"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	// TODO: Add database service when implemented
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind login request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	logger.WithField("username", req.Username).Info("User login attempt")

	// TODO: Implement proper authentication
	// For now, accept any username/password
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
	}

	// TODO: Generate JWT token
	token := "mock-jwt-token-" + req.Username

	user := models.User{
		ID:       "1",
		Username: req.Username,
		Email:    req.Username + "@example.com",
		Role:     "user",
		IsActive: true,
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"user":  user,
			"token": token,
		},
	})
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind register request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	logger.WithField("username", req.Username).Info("User registration attempt")

	// TODO: Implement proper validation and user creation
	// For now, create a mock user
	user := models.User{
		ID:       "2",
		Username: req.Username,
		Email:    req.Email,
		Role:     "user",
		IsActive: true,
	}

	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	logger.Info("User logout")

	// TODO: Implement token invalidation

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Logout successful",
	})
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	// TODO: Get user from JWT token
	userID := "1" // Mock user ID
	logger.WithField("user_id", userID).Info("Getting user profile")

	user := models.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
		IsActive: true,
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	// TODO: Get user from JWT token
	userID := "1" // Mock user ID

	var req struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		logger.WithField("error", err).Error("Failed to bind profile update request")
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	logger.WithField("user_id", userID).Info("Updating user profile")

	// TODO: Implement database update
	user := models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Role:     "user",
		IsActive: true,
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    user,
	})
}
