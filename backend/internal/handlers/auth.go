package handlers

import (
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/config"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/logger"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	// TODO: Add database service when implemented
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

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

	// Check if demo authentication is enabled
	if os.Getenv("DEMO_AUTH") == "" {
		return c.JSON(http.StatusNotImplemented, models.APIResponse{
			Success: false,
			Error:   "Authentication not implemented. Set DEMO_AUTH environment variable for demo mode.",
		})
	}

	// TODO: Implement proper authentication
	// For now, accept any username/password (demo mode only)
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
	}

	// Generate JWT token
	var token string
	var err error

	if os.Getenv("DEMO_AUTH") != "" {
		// Demo mode: use simple prefixed token
		token = "demo-" + req.Username
	} else {
		// Production mode: generate proper JWT
		token, err = generateJWTToken(req.Username, "user")
		if err != nil {
			logger.WithField("error", err).Error("Failed to generate JWT token")
			return c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error:   "Failed to generate authentication token",
			})
		}
	}

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

	// Trim whitespace from input fields
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	// Validate that at least one field is provided
	if req.Username == "" && req.Email == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "No fields to update",
		})
	}

	// Validate email format if provided
	if req.Email != "" && !emailRegex.MatchString(req.Email) {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid email format",
		})
	}

	logger.WithField("user_id", userID).Info("Updating user profile")

	// TODO: Fetch current user from database
	// For now, simulate fetching current user data
	currentUser := models.User{
		ID:       userID,
		Username: "testuser",         // This would come from database
		Email:    "test@example.com", // This would come from database
		Role:     "user",
		IsActive: true,
	}

	// Only update fields that are provided (preserve existing values)
	updatedUser := currentUser
	if req.Username != "" {
		updatedUser.Username = req.Username
	}
	if req.Email != "" {
		updatedUser.Email = req.Email
	}

	// TODO: Implement database update
	// This would update the user in the database with updatedUser

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    updatedUser,
	})
}

// generateJWTToken creates a signed JWT token with standard claims
func generateJWTToken(username, role string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(config.AppConfig.JWT.ExpiresIn) * time.Second)),
		NotBefore: jwt.NewNumericDate(now),
	}

	// Add custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  claims.Subject,
		"iat":  claims.IssuedAt,
		"exp":  claims.ExpiresAt,
		"nbf":  claims.NotBefore,
		"role": role,
	})

	return token.SignedString([]byte(config.AppConfig.JWT.SecretKey))
}
