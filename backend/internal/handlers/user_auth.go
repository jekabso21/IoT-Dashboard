package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/auth"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/repository"
	"github.com/labstack/echo/v4"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if !emailRegex.MatchString(req.Email) {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid email address",
		})
	}

	// Validate password
	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Password must be at least 8 characters",
		})
	}

	// Check if email already exists
	exists, err := repository.EmailExists(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to check email availability",
		})
	}
	if exists {
		return c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Error:   "Email already registered",
		})
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to process registration",
		})
	}

	// Get full name (optional)
	fullName := ""
	if req.FullName != nil {
		fullName = strings.TrimSpace(*req.FullName)
	}

	// Create user
	user, err := repository.CreateUser(c.Request().Context(), req.Email, passwordHash, fullName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create user",
		})
	}

	// Generate JWT token
	token, err := auth.GenerateUserToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
	}

	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data: models.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate email
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Email is required",
		})
	}

	// Validate password
	if req.Password == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Password is required",
		})
	}

	// Get user by email
	user, passwordHash, err := repository.GetUserByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid email or password",
		})
	}

	// Verify password
	if !auth.CheckPasswordHash(req.Password, passwordHash) {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := auth.GenerateUserToken(user.ID, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: models.LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func RefreshUserToken(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	email, ok := c.Get("email").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	// Generate new JWT token
	token, err := auth.GenerateUserToken(userID, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]string{
			"token": token,
		},
	})
}

func Logout(c echo.Context) error {
	// For stateless JWT, logout is handled client-side by removing the token
	// This endpoint can be used for token blacklisting in the future
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]string{
			"message": "Logged out successfully",
		},
	})
}

func GetCurrentUser(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	user, err := repository.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
	}

	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}
