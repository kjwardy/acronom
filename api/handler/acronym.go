package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/kjwardy/acronom/api/model"
	"github.com/labstack/echo/v4"
)

// CreateAcronymRequest represents the JSON request body for creating an acronym
type CreateAcronymRequest struct {
	Acronym    string  `json:"acronym"`    // The acronym text (required)
	Definition string  `json:"definition"` // The definition of the acronym (required)
	Link       *string `json:"link,omitempty"` // Optional link to external resource
}

// UpdateAcronymRequest represents the JSON request body for updating an acronym
type UpdateAcronymRequest struct {
	Definition string  `json:"definition"` // The updated definition (required)
	Link       *string `json:"link,omitempty"` // Optional updated link
}

// CreateAcronym handles POST /api/acronyms requests
// Creates a new acronym entry and returns 201 Created with the persisted record
func (h *Handler) CreateAcronym(c echo.Context) error {
	var req CreateAcronymRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate and normalize the request
	acronym, err := h.validateCreateRequest(req)
	if err != nil {
		return err
	}

	// Attempt to create the acronym in the database
	if err := h.Acronyms.Create(c.Request().Context(), acronym); err != nil {
		if err.Error() == "acronym already exists" {
			return echo.NewHTTPError(http.StatusConflict, "An acronym with this name already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create acronym")
	}

	return c.JSON(http.StatusCreated, acronym)
}

// UpdateAcronym handles PUT /api/acronyms/:id requests
// Updates an existing acronym by ID and returns 200 OK with the updated record
func (h *Handler) UpdateAcronym(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
	}

	var req UpdateAcronymRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate and prepare the update request
	acronym, err := h.validateUpdateRequest(c, id, req)
	if err != nil {
		return err
	}

	// Attempt to update the acronym in the database
	if err := h.Acronyms.Update(c.Request().Context(), acronym); err != nil {
		if err.Error() == "acronym not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Acronym not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update acronym")
	}

	return c.JSON(http.StatusOK, acronym)
}

// DeleteAcronym handles DELETE /api/acronyms/:id requests
// Deletes an acronym by ID and returns 204 No Content
func (h *Handler) DeleteAcronym(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
	}

	// Parse and validate the ID
	parsedID, err := parseID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	// Attempt to delete the acronym from the database
	if err := h.Acronyms.Delete(c.Request().Context(), parsedID); err != nil {
		if err.Error() == "acronym not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Acronym not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete acronym")
	}

	return c.NoContent(http.StatusNoContent)
}

// validateCreateRequest validates and normalizes the create acronym request
// Returns an error if validation fails, otherwise returns a validated Acronym model
func (h *Handler) validateCreateRequest(req CreateAcronymRequest) (*model.Acronym, error) {
	acronym := strings.TrimSpace(req.Acronym)
	definition := strings.TrimSpace(req.Definition)

	if acronym == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Acronym is required")
	}

	if definition == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Definition is required")
	}

	return &model.Acronym{
		Acronym:    acronym,
		Definition: definition,
		Link:       req.Link,
	}, nil
}

// validateUpdateRequest validates and normalizes the update acronym request
// Fetches the existing acronym, validates the update data, and returns the updated model
func (h *Handler) validateUpdateRequest(ctx echo.Context, id string, req UpdateAcronymRequest) (*model.Acronym, error) {
	parsedID, err := parseID(id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	definition := strings.TrimSpace(req.Definition)
	if definition == "" {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Definition is required")
	}

	// Fetch the existing acronym to ensure it exists
	existing, err := h.Acronyms.Find(ctx.Request().Context(), parsedID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve acronym")
	}
	if existing == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Acronym not found")
	}

	// Update the fields that can be modified
	existing.Definition = definition
	existing.Link = req.Link

	return existing, nil
}

// parseID converts a string ID to int64
// Returns an error if the string is not a valid integer
func parseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}
