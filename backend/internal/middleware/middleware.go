package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/config"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.AppConfig.CORS.AllowedOrigins,
		AllowMethods:     config.AppConfig.CORS.AllowedMethods,
		AllowHeaders:     config.AppConfig.CORS.AllowedHeaders,
		AllowCredentials: true,
		MaxAge:           86400,
	})
}

func Logger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogError:     true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.Log.WithFields(logrus.Fields{
				"method":     values.Method,
				"uri":        values.URI,
				"status":     values.Status,
				"latency":    values.Latency,
				"remote_ip":  values.RemoteIP,
				"user_agent": values.UserAgent,
				"error":      values.Error,
			}).Info("HTTP Request")
			return nil
		},
	})
}

func Recovery() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          4, // ERROR level
	})
}

func RateLimit() echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10))
}

func Security() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "default-src 'self'",
	})
}

func RequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return generateRequestID()
		},
	})
}

func Timeout() echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	})
}

func IoTDeviceAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Missing authorization header",
				})
			}

			// Check if it's a Bearer token or API key
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if !isValidDeviceToken(token) {
					return c.JSON(http.StatusUnauthorized, map[string]interface{}{
						"success": false,
						"error":   "Invalid device token",
					})
				}
			} else {
				// API key authentication
				if !isValidAPIKey(authHeader) {
					return c.JSON(http.StatusUnauthorized, map[string]interface{}{
						"success": false,
						"error":   "Invalid API key",
					})
				}
			}

			return next(c)
		}
	}
}

func WebAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Missing authorization header",
				})
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Invalid authorization format",
				})
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if !isValidJWTToken(token) {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"error":   "Invalid or expired token",
				})
			}

			return next(c)
		}
	}
}

// Helper functions
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	// Generate cryptographically secure random bytes
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fail-safe: log error and return a fallback ID
		logger.WithField("error", err).Error("Failed to generate secure random string")
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	// Convert to hex string (2 chars per byte)
	return hex.EncodeToString(bytes)[:length]
}

func isValidDeviceToken(token string) bool {
	if len(token) == 0 {
		return false
	}

	// Check for demo mode
	if os.Getenv("DEMO_AUTH") != "" {
		return strings.HasPrefix(token, "demo-")
	}

	// Check if JWT secret is properly configured
	if config.AppConfig.JWT.SecretKey == "your-secret-key" || config.AppConfig.JWT.SecretKey == "" {
		env := os.Getenv("APP_ENV")
		if env == "production" || env == "staging" {
			logger.Warn("Device token validation rejected - JWT secret not configured in production")
			return false
		}
		logger.Warn("Device token validation using default JWT secret - not secure for production")
	}

	// Production mode: validate JWT signature and claims
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.AppConfig.JWT.SecretKey), nil
	})

	if err != nil {
		logger.WithField("error", err).Debug("Device token validation failed")
		return false
	}

	// Check if token is valid and not expired
	if !parsedToken.Valid {
		return false
	}

	// Validate claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		// Check if subject exists
		if _, exists := claims["sub"]; !exists {
			return false
		}
		// Check if role exists and is device
		if role, exists := claims["role"]; !exists || role != "device" {
			return false
		}
	}

	return true
}

func isValidAPIKey(apiKey string) bool {
	if len(apiKey) == 0 {
		return false
	}

	// Check for demo mode
	if os.Getenv("DEMO_AUTH") != "" {
		return strings.HasPrefix(apiKey, "demo-")
	}

	// In production, reject API key validation until database layer is implemented
	// This prevents accidental deployment with insecure API key validation
	env := os.Getenv("APP_ENV")
	if env == "production" || env == "staging" {
		logger.WithField("api_key_prefix", apiKey[:min(8, len(apiKey))]).Warn("API key validation not implemented - rejecting in production")
		return false
	}

	// Development mode: accept any non-empty API key with warning
	logger.WithField("api_key_prefix", apiKey[:min(8, len(apiKey))]).Warn("API key validation not implemented - accepting in development")
	return true
}

func isValidJWTToken(token string) bool {
	if len(token) == 0 {
		return false
	}

	// Check for demo mode
	if os.Getenv("DEMO_AUTH") != "" {
		return strings.HasPrefix(token, "demo-")
	}

	// Check if JWT secret is properly configured
	if config.AppConfig.JWT.SecretKey == "your-secret-key" || config.AppConfig.JWT.SecretKey == "" {
		env := os.Getenv("APP_ENV")
		if env == "production" || env == "staging" {
			logger.Warn("JWT token validation rejected - JWT secret not configured in production")
			return false
		}
		logger.Warn("JWT token validation using default JWT secret - not secure for production")
	}

	// Production mode: validate JWT signature and claims
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.AppConfig.JWT.SecretKey), nil
	})

	if err != nil {
		logger.WithField("error", err).Debug("JWT validation failed")
		return false
	}

	// Check if token is valid and not expired
	if !parsedToken.Valid {
		return false
	}

	// Validate claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		// Check if subject exists
		if _, exists := claims["sub"]; !exists {
			return false
		}
		// Check if role exists
		if _, exists := claims["role"]; !exists {
			return false
		}
	}

	return true
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
