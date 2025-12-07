package routes

import (
	"github.com/jekabso21/IoT-Dashboard/backend/internal/handlers"
	customMiddleware "github.com/jekabso21/IoT-Dashboard/backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":  "healthy",
			"service": "IoT Dashboard API",
		})
	})

	api := e.Group("/api")

	// Auth routes (public - no middleware required)
	auth := api.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/logout", handlers.Logout)

	// Auth routes (protected - require user auth)
	authProtected := api.Group("/auth", customMiddleware.UserAuth())
	authProtected.POST("/refresh", handlers.RefreshUserToken)
	authProtected.GET("/me", handlers.GetCurrentUser)

	// Device routes (for ESP32 devices)
	devices := api.Group("/devices")
	devices.POST("/pair", handlers.PairDevice)
	devices.POST("/token/refresh", handlers.RefreshDeviceToken, customMiddleware.DeviceAuth())
	devices.POST("/data", handlers.IngestSensorData, customMiddleware.DeviceAuth())

	// User routes (protected - require user auth)
	user := api.Group("/user", customMiddleware.UserAuth())
	user.POST("/pairing/generate", handlers.GeneratePairingCode)
	user.GET("/devices", handlers.GetUserDevices)
	user.GET("/devices/:id", handlers.GetDevice)
	user.PUT("/devices/:id", handlers.UpdateDevice)
	user.DELETE("/devices/:id", handlers.DeleteDevice)
	user.GET("/sensor-data", handlers.GetSensorData)
}
