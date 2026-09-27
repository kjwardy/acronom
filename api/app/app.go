package app

import (
	"context"
	"net/http"
	"time"

	"github.com/kjwardy/acronom/api/db"
	"github.com/kjwardy/acronom/api/handler"
	"github.com/kjwardy/acronom/api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HealthChecker interface {
	Ping(context.Context) error
}

// Init initializes the Echo application with middleware, routes, and handlers
// This follows the go-url pattern of centralized route registration
func Init(e *echo.Echo, database HealthChecker) {
	// Type assert to get the concrete DB instance and create the acronym store
	dbInstance := database.(*db.DB)
	acronymStore := model.NewPGAcronymStore(dbInstance.Pool())

	// Create the handler with the acronym store dependency
	h := &handler.Handler{
		Acronyms: acronymStore,
	}

	// Apply global middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := database.Ping(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]bool{"UP": false})
		}
		return c.JSON(http.StatusOK, map[string]bool{"UP": true})
	})
	e.GET("/api/status", h.Status)

	// API routes for acronym management
	api := e.Group("/api")
	api.POST("/acronyms", h.CreateAcronym)        // Create new acronym
	api.PUT("/acronyms/:id", h.UpdateAcronym)    // Update existing acronym
	api.DELETE("/acronyms/:id", h.DeleteAcronym) // Delete acronym

	// Frontend static file serving
	frontend := e.Group("/acronom")
	frontend.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "public",
		HTML5: true,
	}))

	// Root redirect to frontend
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/acronom")
	})
}
