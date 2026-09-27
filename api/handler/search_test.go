package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kjwardy/acronom/api/model"
)

func TestSearchAcronyms(t *testing.T) {
	t.Run("returns trimmed search results", func(t *testing.T) {
		calledWith := ""
		store := &mockAcronymStore{
			searchFunc: func(_ context.Context, query string) ([]model.Acronym, error) {
				calledWith = query
				return []model.Acronym{{ID: 1, Acronym: "API", Definition: "Application Programming Interface"}}, nil
			},
		}
		h, e := setupTestHandler(store)
		request := httptest.NewRequest(http.MethodGet, "/api/search?q=%20API%20", nil)
		response := httptest.NewRecorder()
		ctx := e.NewContext(request, response)

		if err := h.SearchAcronyms(ctx); err != nil {
			t.Fatal(err)
		}
		if calledWith != "API" {
			t.Fatalf("expected trimmed query, got %q", calledWith)
		}
		if response.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", response.Code)
		}
	})

	t.Run("returns an empty array when there are no matches", func(t *testing.T) {
		h, e := setupTestHandler(&mockAcronymStore{})
		request := httptest.NewRequest(http.MethodGet, "/api/search?q=missing", nil)
		response := httptest.NewRecorder()
		ctx := e.NewContext(request, response)

		if err := h.SearchAcronyms(ctx); err != nil {
			t.Fatal(err)
		}
		if response.Body.String() != "[]\n" {
			t.Fatalf("expected empty array, got %s", response.Body.String())
		}
	})

	t.Run("rejects empty queries without using the store", func(t *testing.T) {
		called := false
		store := &mockAcronymStore{
			searchFunc: func(context.Context, string) ([]model.Acronym, error) {
				called = true
				return nil, nil
			},
		}
		h, e := setupTestHandler(store)
		request := httptest.NewRequest(http.MethodGet, "/api/search?q=%20%20", nil)
		response := httptest.NewRecorder()
		ctx := e.NewContext(request, response)

		err := h.SearchAcronyms(ctx)
		if err == nil {
			t.Fatal("expected empty query error")
		}
		if called {
			t.Fatal("expected store not to be called")
		}
	})

	t.Run("hides database failures", func(t *testing.T) {
		store := &mockAcronymStore{
			searchFunc: func(context.Context, string) ([]model.Acronym, error) {
				return nil, errors.New("database credentials and details")
			},
		}
		h, e := setupTestHandler(store)
		request := httptest.NewRequest(http.MethodGet, "/api/search?q=api", nil)
		response := httptest.NewRecorder()
		ctx := e.NewContext(request, response)

		err := h.SearchAcronyms(ctx)
		if err == nil {
			t.Fatal("expected database error")
		}
		if err.Error() != "code=500, message=Failed to search acronyms" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
