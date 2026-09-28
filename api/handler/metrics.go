package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TotalAcronymsMetric struct {
	TotalAcronyms int64 `json:"total_acronyms"`
}

func (h *Handler) TotalAcronyms(c echo.Context) error {
	total, err := h.Metrics.TotalAcronyms(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load acronym metrics")
	}
	return c.JSON(http.StatusOK, TotalAcronymsMetric{TotalAcronyms: total})
}
