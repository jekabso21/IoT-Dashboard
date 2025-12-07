package middleware

import (
	"net/http"
	"strings"

	"github.com/jekabso21/IoT-Dashboard/backend/internal/auth"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/models"
	"github.com/labstack/echo/v4"
)

func UserAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Error:   "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Error:   "Invalid authorization header format",
				})
			}

			token := parts[1]
			claims, err := auth.ValidateUserToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Error:   "Invalid or expired token",
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

func DeviceAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Error:   "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Error:   "Invalid authorization header format",
				})
			}

			token := parts[1]
			claims, err := auth.ValidateDeviceToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Error:   "Invalid or expired token",
				})
			}

			c.Set("device_id", claims.DeviceID)
			c.Set("user_id", claims.UserID)

			return next(c)
		}
	}
}
