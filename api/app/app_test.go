package app

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// mockDatabaseProvider is a mock implementation of DatabaseProvider for testing
type mockDatabaseProvider struct {
	pool *pgxpool.Pool
}

func (m *mockDatabaseProvider) Ping(ctx context.Context) error {
	return nil
}

func (m *mockDatabaseProvider) Pool() *pgxpool.Pool {
	return m.pool
}

func TestInitWithMockDatabase(t *testing.T) {
	// Test that Init doesn't panic when passed a mock DatabaseProvider
	mockDB := &mockDatabaseProvider{pool: nil}

	e := echo.New()

	// This should not panic since we're using the interface instead of type assertion
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Init panicked with mock database: %v", r)
		}
	}()

	Init(e, mockDB)
}
