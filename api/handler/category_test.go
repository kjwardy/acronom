package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kjwardy/acronom/api/model"
	"github.com/labstack/echo/v4"
)

type mockCategoryStore struct {
	categories []model.Category
	create     func(*model.Category) error
	err        error
}

func (store *mockCategoryStore) List(context.Context) ([]model.Category, error) {
	return store.categories, store.err
}

func (store *mockCategoryStore) Create(_ context.Context, category *model.Category) error {
	if store.create != nil {
		return store.create(category)
	}
	return store.err
}

func TestListCategories(t *testing.T) {
	e := echo.New()
	h := &Handler{Categories: &mockCategoryStore{categories: []model.Category{}}}
	request := httptest.NewRequest(http.MethodGet, "/api/categories", nil)
	response := httptest.NewRecorder()

	if err := h.ListCategories(e.NewContext(request, response)); err != nil {
		t.Fatal(err)
	}
	if response.Body.String() != "[]\n" {
		t.Fatalf("Expected empty array, got %s", response.Body.String())
	}
}

func TestCreateCategory(t *testing.T) {
	t.Run("normalizes and creates an approved category", func(t *testing.T) {
		store := &mockCategoryStore{create: func(category *model.Category) error {
			category.ID = 1
			return nil
		}}
		e := echo.New()
		h := &Handler{Categories: store}
		request := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(`{"name":"  platform tools ","colour":"#1976d2"}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		response := httptest.NewRecorder()

		if err := h.CreateCategory(e.NewContext(request, response)); err != nil {
			t.Fatal(err)
		}
		if response.Code != http.StatusCreated {
			t.Fatalf("Expected 201, got %d", response.Code)
		}
		var category model.Category
		if err := json.Unmarshal(response.Body.Bytes(), &category); err != nil {
			t.Fatal(err)
		}
		if category.ID != 1 || category.Name != "Platform Tools" || category.Colour != "#1976D2" {
			t.Fatalf("Unexpected response: %#v", category)
		}
	})

	tests := []struct {
		name string
		body string
	}{
		{name: "empty name", body: `{"name":" ","colour":"#1976D2"}`},
		{name: "long name", body: `{"name":"This category name is definitely too long","colour":"#1976D2"}`},
		{name: "unapproved colour", body: `{"name":"Platform","colour":"#FFFFFF"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			h := &Handler{Categories: &mockCategoryStore{}}
			request := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(test.body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			err := h.CreateCategory(e.NewContext(request, httptest.NewRecorder()))
			httpError, ok := err.(*echo.HTTPError)
			if !ok || httpError.Code != http.StatusBadRequest {
				t.Fatalf("Expected 400, got %v", err)
			}
		})
	}

	t.Run("duplicate name", func(t *testing.T) {
		e := echo.New()
		h := &Handler{Categories: &mockCategoryStore{err: model.ErrCategoryExists}}
		request := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(`{"name":"Platform","colour":"#1976D2"}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		err := h.CreateCategory(e.NewContext(request, httptest.NewRecorder()))
		httpError, ok := err.(*echo.HTTPError)
		if !ok || httpError.Code != http.StatusConflict {
			t.Fatalf("Expected 409, got %v", err)
		}
	})

	t.Run("hides database failure", func(t *testing.T) {
		e := echo.New()
		h := &Handler{Categories: &mockCategoryStore{err: errors.New("database details")}}
		request := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(`{"name":"Platform","colour":"#1976D2"}`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		err := h.CreateCategory(e.NewContext(request, httptest.NewRecorder()))
		httpError, ok := err.(*echo.HTTPError)
		if !ok || httpError.Code != http.StatusInternalServerError {
			t.Fatalf("Expected 500, got %v", err)
		}
	})
}
