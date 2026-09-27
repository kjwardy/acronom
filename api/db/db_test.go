package db

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kjwardy/acronom/api/config"
)

// Note: These tests are currently disabled because they require mocking pgxpool.Pool
// which is complex. For now, we rely on integration tests with a real database.
// These can be re-enabled if we add a more sophisticated mocking approach.

func TestConnectHidesConnectorError(t *testing.T) {
	internalError := errors.New("postgres://app:secret@localhost/acronom")
	databaseConfig := config.Database{
		Addr:     "localhost:5432",
		Database: "acronom",
		User:     "postgres",
		Pass:     "secret",
	}

	_, err := connectWith(context.Background(), databaseConfig, func(context.Context, *pgxpool.Config) (*pgxpool.Pool, error) {
		return nil, internalError
	})
	if err == nil {
		t.Fatal("expected connection error")
	}
	if errors.Is(err, internalError) {
		t.Fatal("expected internal connection error to be hidden")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatal("expected password to be hidden")
	}
}
