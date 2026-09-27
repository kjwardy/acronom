package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kjwardy/acronom/api/model"
	"github.com/labstack/echo/v4"
)

// mockAcronymStore is a mock implementation of AcronymStore for testing
// It allows customizing the behavior of each method via function fields
type mockAcronymStore struct {
	createFunc        func(ctx context.Context, acronym *model.Acronym) error
	updateFunc        func(ctx context.Context, acronym *model.Acronym) error
	deleteFunc        func(ctx context.Context, id int64) error
	findFunc          func(ctx context.Context, id int64) (*model.Acronym, error)
	findByAcronymFunc func(ctx context.Context, acronym string) (*model.Acronym, error)
	searchFunc        func(ctx context.Context, query string) ([]model.Acronym, error)
}

func (m *mockAcronymStore) Create(ctx context.Context, acronym *model.Acronym) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, acronym)
	}
	return nil
}

func (m *mockAcronymStore) Update(ctx context.Context, acronym *model.Acronym) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, acronym)
	}
	return nil
}

func (m *mockAcronymStore) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockAcronymStore) Find(ctx context.Context, id int64) (*model.Acronym, error) {
	if m.findFunc != nil {
		return m.findFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockAcronymStore) FindByAcronym(ctx context.Context, acronym string) (*model.Acronym, error) {
	if m.findByAcronymFunc != nil {
		return m.findByAcronymFunc(ctx, acronym)
	}
	return nil, nil
}

func (m *mockAcronymStore) Search(ctx context.Context, query string) ([]model.Acronym, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, query)
	}
	return []model.Acronym{}, nil
}

// setupTestHandler creates a test handler and Echo instance with the given mock store
func setupTestHandler(store model.AcronymStore) (*Handler, *echo.Echo) {
	e := echo.New()
	h := &Handler{Acronyms: store}
	return h, e
}

func TestCreateAcronym(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		store := &mockAcronymStore{
			createFunc: func(ctx context.Context, acronym *model.Acronym) error {
				acronym.ID = 1
				return nil
			},
		}
		h, e := setupTestHandler(store)

		reqBody := `{"acronym":"API","definition":"Application Programming Interface"}`
		req := httptest.NewRequest(http.MethodPost, "/api/acronyms", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateAcronym(c)
		if err != nil {
			t.Errorf("CreateAcronym() failed: %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		var response model.Acronym
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}
		if response.Acronym != "API" {
			t.Errorf("Expected acronym 'API', got '%s'", response.Acronym)
		}
	})

	t.Run("empty acronym", func(t *testing.T) {
		store := &mockAcronymStore{}
		h, e := setupTestHandler(store)

		reqBody := `{"acronym":"","definition":"test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/acronyms", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateAcronym(c)
		if err == nil {
			t.Error("Expected error for empty acronym")
		}
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 status, got %v", err)
			}
		}
	})

	t.Run("empty definition", func(t *testing.T) {
		store := &mockAcronymStore{}
		h, e := setupTestHandler(store)

		reqBody := `{"acronym":"API","definition":""}`
		req := httptest.NewRequest(http.MethodPost, "/api/acronyms", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateAcronym(c)
		if err == nil {
			t.Error("Expected error for empty definition")
		}
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 status, got %v", err)
			}
		}
	})

	t.Run("duplicate meaning", func(t *testing.T) {
		store := &mockAcronymStore{
			createFunc: func(ctx context.Context, acronym *model.Acronym) error {
				return model.ErrAcronymMeaningExists
			},
		}
		h, e := setupTestHandler(store)

		reqBody := `{"acronym":"PC","definition":"Personal Computer"}`
		req := httptest.NewRequest(http.MethodPost, "/api/acronyms", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateAcronym(c)
		if err == nil {
			t.Fatal("Expected duplicate meaning error")
		}
		httpErr, ok := err.(*echo.HTTPError)
		if !ok || httpErr.Code != http.StatusConflict {
			t.Fatalf("Expected 409 status, got %v", err)
		}
	})

	t.Run("database failure", func(t *testing.T) {
		store := &mockAcronymStore{
			createFunc: func(ctx context.Context, acronym *model.Acronym) error {
				return fmt.Errorf("database unavailable")
			},
		}
		h, e := setupTestHandler(store)

		reqBody := `{"acronym":"API","definition":"Application Programming Interface"}`
		req := httptest.NewRequest(http.MethodPost, "/api/acronyms", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.CreateAcronym(c)
		if err == nil {
			t.Error("Expected database error")
		}
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusInternalServerError {
				t.Errorf("Expected 500 status, got %v", err)
			}
		}
	})
}

func TestUpdateAcronym(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		existing := &model.Acronym{
			ID:         1,
			Acronym:    "API",
			Definition: "Old definition",
		}
		store := &mockAcronymStore{
			findFunc: func(ctx context.Context, id int64) (*model.Acronym, error) {
				return existing, nil
			},
			updateFunc: func(ctx context.Context, acronym *model.Acronym) error {
				return nil
			},
		}
		h, e := setupTestHandler(store)

		reqBody := `{"definition":"New definition"}`
		req := httptest.NewRequest(http.MethodPut, "/api/acronyms/1", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.UpdateAcronym(c)
		if err != nil {
			t.Errorf("UpdateAcronym() failed: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("duplicate meaning", func(t *testing.T) {
		store := &mockAcronymStore{
			findFunc: func(ctx context.Context, id int64) (*model.Acronym, error) {
				return &model.Acronym{ID: id, Acronym: "PC"}, nil
			},
			updateFunc: func(ctx context.Context, acronym *model.Acronym) error {
				return model.ErrAcronymMeaningExists
			},
		}
		h, e := setupTestHandler(store)
		reqBody := `{"definition":"Personal Computer"}`
		req := httptest.NewRequest(http.MethodPut, "/api/acronyms/1", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.UpdateAcronym(c)
		httpErr, ok := err.(*echo.HTTPError)
		if !ok || httpErr.Code != http.StatusConflict {
			t.Fatalf("Expected 409 status, got %v", err)
		}
	})

	t.Run("acronym not found", func(t *testing.T) {
		store := &mockAcronymStore{
			findFunc: func(ctx context.Context, id int64) (*model.Acronym, error) {
				return nil, nil
			},
		}
		h, e := setupTestHandler(store)

		reqBody := `{"definition":"New definition"}`
		req := httptest.NewRequest(http.MethodPut, "/api/acronyms/999", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := h.UpdateAcronym(c)
		if err == nil {
			t.Error("Expected error for non-existent acronym")
		}
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusNotFound {
				t.Errorf("Expected 404 status, got %v", err)
			}
		}
	})

	t.Run("invalid ID format", func(t *testing.T) {
		store := &mockAcronymStore{}
		h, e := setupTestHandler(store)

		reqBody := `{"definition":"New definition"}`
		req := httptest.NewRequest(http.MethodPut, "/api/acronyms/invalid", bytes.NewBufferString(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := h.UpdateAcronym(c)
		if err == nil {
			t.Error("Expected error for invalid ID format")
		}
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 status, got %v", err)
			}
		}
	})
}

func TestDeleteAcronym(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		store := &mockAcronymStore{
			deleteFunc: func(ctx context.Context, id int64) error {
				return nil
			},
		}
		h, e := setupTestHandler(store)

		req := httptest.NewRequest(http.MethodDelete, "/api/acronyms/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := h.DeleteAcronym(c)
		if err != nil {
			t.Errorf("DeleteAcronym() failed: %v", err)
		}

		if rec.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, rec.Code)
		}
	})

	t.Run("acronym not found", func(t *testing.T) {
		store := &mockAcronymStore{
			deleteFunc: func(ctx context.Context, id int64) error {
				return fmt.Errorf("acronym not found")
			},
		}
		h, e := setupTestHandler(store)

		req := httptest.NewRequest(http.MethodDelete, "/api/acronyms/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := h.DeleteAcronym(c)
		if err == nil {
			t.Error("Expected error for non-existent acronym")
		}
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusNotFound {
				t.Errorf("Expected 404 status, got %v", err)
			}
		}
	})

	t.Run("invalid ID format", func(t *testing.T) {
		store := &mockAcronymStore{}
		h, e := setupTestHandler(store)

		req := httptest.NewRequest(http.MethodDelete, "/api/acronyms/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := h.DeleteAcronym(c)
		if err == nil {
			t.Error("Expected error for invalid ID format")
		}
		if err != nil {
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 status, got %v", err)
			}
		}
	})
}

func TestValidateCreateRequest(t *testing.T) {
	h := &Handler{}

	t.Run("valid request", func(t *testing.T) {
		req := CreateAcronymRequest{
			Acronym:    "API",
			Definition: "Application Programming Interface",
		}

		acronym, err := h.validateCreateRequest(req)
		if err != nil {
			t.Errorf("validateCreateRequest() failed: %v", err)
		}
		if acronym.Acronym != "API" {
			t.Errorf("Expected acronym 'API', got '%s'", acronym.Acronym)
		}
	})

	t.Run("empty acronym", func(t *testing.T) {
		req := CreateAcronymRequest{
			Acronym:    "",
			Definition: "test",
		}

		_, err := h.validateCreateRequest(req)
		if err == nil {
			t.Error("Expected error for empty acronym")
		}
	})

	t.Run("empty definition", func(t *testing.T) {
		req := CreateAcronymRequest{
			Acronym:    "API",
			Definition: "",
		}

		_, err := h.validateCreateRequest(req)
		if err == nil {
			t.Error("Expected error for empty definition")
		}
	})

	t.Run("trimming", func(t *testing.T) {
		req := CreateAcronymRequest{
			Acronym:    "  API  ",
			Definition: "  Application Programming Interface  ",
		}

		acronym, err := h.validateCreateRequest(req)
		if err != nil {
			t.Errorf("validateCreateRequest() failed: %v", err)
		}
		if acronym.Acronym != "API" {
			t.Errorf("Expected trimmed acronym 'API', got '%s'", acronym.Acronym)
		}
		if acronym.Definition != "Application Programming Interface" {
			t.Errorf("Expected trimmed definition, got '%s'", acronym.Definition)
		}
	})
}
