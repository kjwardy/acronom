package app

import (
	"context"
	"net/http"
	"time"

	"github.com/kjwardy/acronom/api/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HealthChecker interface {
	Ping(context.Context) error
}

func Init(e *echo.Echo, database HealthChecker) {
	h := &handler.Handler{}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := database.Ping(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]bool{"UP": false})
		}
		return c.JSON(http.StatusOK, map[string]bool{"UP": true})
	})
	e.GET("/api/status", h.Status)

	frontend := e.Group("/acronom")
	frontend.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "public",
		HTML5: true,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/acronom")
	})
}
