package middleware

import (
	"net/http"
	"strings"
	"time"

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
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Rate:  10, // requests per second
		Burst: 20, // burst size
		KeyGenerator: func(c echo.Context) string {
			return c.RealIP()
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"success": false,
				"error":   "Rate limit exceeded",
			})
		},
	})
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
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func isValidDeviceToken(token string) bool {
	// TODO: Implement proper device token validation
	// For now, accept any non-empty token
	return len(token) > 0
}

func isValidAPIKey(apiKey string) bool {
	// TODO: Implement proper API key validation
	// For now, accept any non-empty API key
	return len(apiKey) > 0
}

func isValidJWTToken(token string) bool {
	// TODO: Implement proper JWT token validation
	// For now, accept any non-empty token
	return len(token) > 0
}
