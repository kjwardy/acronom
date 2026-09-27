package handler

import (
	"net/http"

	"github.com/kjwardy/acronom/api/model"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Acronyms model.AcronymStore
}

type Status struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (h *Handler) Status(c echo.Context) error {
	return c.JSON(http.StatusOK, Status{
		Name:   "Acronom API",
		Status: "ready",
	})
}
