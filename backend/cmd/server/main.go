package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jekabso21/IoT-Dashboard/backend/internal/config"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/logger"
	customMiddleware "github.com/jekabso21/IoT-Dashboard/backend/internal/middleware"
	"github.com/jekabso21/IoT-Dashboard/backend/internal/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		logger.Fatal("Failed to load configuration:", err)
	}

	// Create Echo instance
	e := echo.New()

	// Hide Echo banner
	e.HideBanner = true

	// Setup global middleware
	setupGlobalMiddleware(e)

	// Setup routes
	routes.SetupRoutes(e)

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	logger.WithFields(logrus.Fields{
		"host": config.AppConfig.Server.Host,
		"port": config.AppConfig.Server.Port,
	}).Info("Starting IoT Dashboard API server")

	// Start server in a goroutine
	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Server exited")
}

func setupGlobalMiddleware(e *echo.Echo) {
	// Request ID middleware
	e.Use(customMiddleware.RequestID())

	// Recovery middleware
	e.Use(customMiddleware.Recovery())

	// Security middleware
	e.Use(customMiddleware.Security())

	// CORS middleware
	e.Use(customMiddleware.CORS())

	// Request logging middleware
	e.Use(customMiddleware.Logger())

	// Rate limiting middleware
	e.Use(customMiddleware.RateLimit())

	// Timeout middleware
	e.Use(customMiddleware.Timeout())

	// Body limit middleware
	e.Use(middleware.BodyLimit("10MB"))

	// Gzip compression
	e.Use(middleware.Gzip())

	// Custom error handler
	e.HTTPErrorHandler = customErrorHandler
}

func customErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		} else {
			message = fmt.Sprintf("%v", he.Message)
		}
	}

	logger.WithFields(logrus.Fields{
		"error":  err.Error(),
		"path":   c.Request().URL.Path,
		"method": c.Request().Method,
		"status": code,
	}).Error("HTTP Error")

	if !c.Response().Committed {
		c.JSON(code, map[string]interface{}{
			"success": false,
			"error":   message,
		})
	}
}
