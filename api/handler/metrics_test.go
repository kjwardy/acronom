package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockMetricsStore struct {
	total int64
	err   error
}

func (store *mockMetricsStore) TotalAcronyms(context.Context) (int64, error) {
	return store.total, store.err
}

func TestTotalAcronyms(t *testing.T) {
	t.Run("returns total entries", func(t *testing.T) {
		e := echo.New()
		h := &Handler{Metrics: &mockMetricsStore{total: 2}}
		request := httptest.NewRequest(http.MethodGet, "/api/metrics/total-acronyms", nil)
		response := httptest.NewRecorder()

		if err := h.TotalAcronyms(e.NewContext(request, response)); err != nil {
			t.Fatal(err)
		}
		if response.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", response.Code)
		}
		if response.Body.String() != "{\"total_acronyms\":2}\n" {
			t.Fatalf("Unexpected response: %s", response.Body.String())
		}
	})

	t.Run("returns zero for an empty database", func(t *testing.T) {
		e := echo.New()
		h := &Handler{Metrics: &mockMetricsStore{}}
		request := httptest.NewRequest(http.MethodGet, "/api/metrics/total-acronyms", nil)
		response := httptest.NewRecorder()

		if err := h.TotalAcronyms(e.NewContext(request, response)); err != nil {
			t.Fatal(err)
		}
		if response.Body.String() != "{\"total_acronyms\":0}\n" {
			t.Fatalf("Unexpected response: %s", response.Body.String())
		}
	})

	t.Run("hides database failures", func(t *testing.T) {
		e := echo.New()
		h := &Handler{Metrics: &mockMetricsStore{err: errors.New("database credentials")}}
		request := httptest.NewRequest(http.MethodGet, "/api/metrics/total-acronyms", nil)
		response := httptest.NewRecorder()

		err := h.TotalAcronyms(e.NewContext(request, response))
		httpError, ok := err.(*echo.HTTPError)
		if !ok || httpError.Code != http.StatusInternalServerError {
			t.Fatalf("Expected status 500, got %v", err)
		}
		if httpError.Message != "Failed to load acronym metrics" {
			t.Fatalf("Unexpected error message: %v", httpError.Message)
		}
	})
}
