package handler

import (
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/kjwardy/acronom/api/model"
	"github.com/labstack/echo/v4"
)

type CreateCategoryRequest struct {
	Name   string `json:"name"`
	Colour string `json:"colour"`
}

func (h *Handler) ListCategories(c echo.Context) error {
	categories, err := h.Categories.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load categories")
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *Handler) CreateCategory(c echo.Context) error {
	var request CreateCategoryRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	name := model.NormalizeCategoryName(request.Name)
	colour := strings.ToUpper(strings.TrimSpace(request.Colour))
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Category name is required")
	}
	if utf8.RuneCountInString(name) > model.MaxCategoryNameLength {
		return echo.NewHTTPError(http.StatusBadRequest, "Category name must be 30 characters or fewer")
	}
	if !model.IsApprovedCategoryColour(colour) {
		return echo.NewHTTPError(http.StatusBadRequest, "Select an approved category colour")
	}

	category := &model.Category{Name: name, Colour: colour}
	if err := h.Categories.Create(c.Request().Context(), category); err != nil {
		if errors.Is(err, model.ErrCategoryExists) {
			return echo.NewHTTPError(http.StatusConflict, "A category with this name already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create category")
	}
	return c.JSON(http.StatusCreated, category)
}
