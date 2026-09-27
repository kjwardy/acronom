package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) SearchAcronyms(c echo.Context) error {
	query := strings.TrimSpace(c.QueryParam("q"))
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Search query is required")
	}

	results, err := h.Acronyms.Search(c.Request().Context(), query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to search acronyms")
	}
	return c.JSON(http.StatusOK, results)
}
