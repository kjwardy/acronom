package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHealth(t *testing.T) {
	e := echo.New()
	Init(e)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if response.Body.String() != "{\"UP\":true}\n" {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}

func TestStatus(t *testing.T) {
	e := echo.New()
	Init(e)

	request := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if response.Body.String() != "{\"name\":\"Acronom API\",\"status\":\"ready\"}\n" {
		t.Fatalf("unexpected response: %s", response.Body.String())
	}
}

func TestRootRedirect(t *testing.T) {
	e := echo.New()
	Init(e)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)

	if response.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status 307, got %d", response.Code)
	}
	if response.Header().Get("Location") != "/acronom" {
		t.Fatalf("unexpected redirect: %s", response.Header().Get("Location"))
	}
}
