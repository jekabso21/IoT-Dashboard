package routes

import (
	"github.com/jekabso21/IoT-Dashboard/backend/internal/handlers"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Initialize handlers
	deviceHandler := handlers.NewDeviceHandler()
	iotHandler := handlers.NewIoTHandler()
	authHandler := handlers.NewAuthHandler()

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":  "healthy",
			"service": "IoT Dashboard API",
		})
	})

	// API v1 group
	v1 := e.Group("/api/v1")

	// Public routes (no authentication required)
	public := v1.Group("/public")
	public.POST("/auth/login", authHandler.Login)
	public.POST("/auth/register", authHandler.Register)

	// Web application routes (JWT authentication required)
	web := v1.Group("/web")
	web.Use(middleware.WebAuth())

	// Authentication routes
	web.POST("/auth/logout", authHandler.Logout)
	web.GET("/auth/profile", authHandler.GetProfile)
	web.PUT("/auth/profile", authHandler.UpdateProfile)

	// Device management routes
	web.GET("/devices", deviceHandler.GetDevices)
	web.GET("/devices/:id", deviceHandler.GetDevice)
	web.POST("/devices", deviceHandler.CreateDevice)
	web.PUT("/devices/:id", deviceHandler.UpdateDevice)
	web.DELETE("/devices/:id", deviceHandler.DeleteDevice)
	web.GET("/devices/:id/data", deviceHandler.GetDeviceData)

	// IoT device routes (API key or device token authentication required)
	iot := v1.Group("/iot")
	iot.Use(middleware.IoTDeviceAuth())

	// Sensor data endpoints
	iot.POST("/data", iotHandler.SendSensorData)
	iot.GET("/devices/:id/status", iotHandler.GetDeviceStatus)
	iot.PUT("/devices/:id/status", iotHandler.UpdateDeviceStatus)
	iot.GET("/devices/:id/commands", iotHandler.GetDeviceCommands)
	iot.POST("/devices/:id/commands/:commandId/response", iotHandler.SendCommandResponse)

	// Analytics and reporting routes (web authentication required)
	analytics := v1.Group("/analytics")
	analytics.Use(middleware.WebAuth())

	// TODO: Add analytics endpoints
	analytics.GET("/devices/:id/summary", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Analytics endpoint - to be implemented",
		})
	})
}
