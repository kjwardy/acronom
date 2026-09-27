package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kjwardy/acronom/api/handler"
	"github.com/kjwardy/acronom/api/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HealthChecker defines the interface for checking database connectivity
type HealthChecker interface {
	Ping(context.Context) error
}

// DatabaseProvider defines the interface for providing database access
// This allows Init to work with both real DB instances and mock implementations
type DatabaseProvider interface {
	HealthChecker
	Pool() *pgxpool.Pool
}

// Init initializes the Echo application with middleware, routes, and handlers
// This follows the go-url pattern of centralized route registration
func Init(e *echo.Echo, database DatabaseProvider) {
	// Get the database pool from the provider and create the acronym store
	acronymStore := model.NewPGAcronymStore(database.Pool())

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
	api.GET("/search", h.SearchAcronyms)
	api.POST("/acronyms", h.CreateAcronym)       // Create new acronym
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
