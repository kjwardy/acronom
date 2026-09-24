package app

import (
	"net/http"

	"github.com/kjwardy/acronom/api/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo) {
	h := &handler.Handler{}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/health", func(c echo.Context) error {
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
