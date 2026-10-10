package handler

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNormalizeCategoryIDs(t *testing.T) {
	ids, err := normalizeCategoryIDs([]int64{3, 1, 3, 2})
	if err != nil {
		t.Fatal(err)
	}
	expected := []int64{3, 1, 2}
	for index := range expected {
		if ids[index] != expected[index] {
			t.Fatalf("Expected deduplicated order %#v, got %#v", expected, ids)
		}
	}

	_, err = normalizeCategoryIDs([]int64{1, 2, 3, 4, 5, 6})
	httpError, ok := err.(*echo.HTTPError)
	if !ok || httpError.Code != http.StatusBadRequest {
		t.Fatalf("Expected category limit error, got %v", err)
	}

	_, err = normalizeCategoryIDs([]int64{0})
	httpError, ok = err.(*echo.HTTPError)
	if !ok || httpError.Code != http.StatusBadRequest {
		t.Fatalf("Expected positive ID error, got %v", err)
	}
}
