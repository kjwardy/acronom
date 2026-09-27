package db

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kjwardy/acronom/api/config"
)

type fakePool struct {
	pingErr error
	closed  bool
}

func (pool *fakePool) Ping(context.Context) error {
	return pool.pingErr
}

func (pool *fakePool) Close() {
	pool.closed = true
}

func TestConnectConfiguresAndChecksPool(t *testing.T) {
	fake := &fakePool{}
	databaseConfig := config.Database{
		Addr:     "database:5433",
		Database: "acronom",
		User:     "app",
		Pass:     "secret",
	}

	database, err := connectWith(context.Background(), databaseConfig, func(_ context.Context, poolConfig *pgxpool.Config) (pool, error) {
		if poolConfig.ConnConfig.Host != "database" {
			t.Fatalf("unexpected host: %s", poolConfig.ConnConfig.Host)
		}
		if poolConfig.ConnConfig.Port != 5433 {
			t.Fatalf("unexpected port: %d", poolConfig.ConnConfig.Port)
		}
		if poolConfig.ConnConfig.Database != "acronom" {
			t.Fatalf("unexpected database: %s", poolConfig.ConnConfig.Database)
		}
		if poolConfig.ConnConfig.User != "app" {
			t.Fatalf("unexpected user: %s", poolConfig.ConnConfig.User)
		}
		if poolConfig.ConnConfig.Password != "secret" {
			t.Fatal("unexpected password")
		}
		return fake, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	database.Close()
	if !fake.closed {
		t.Fatal("expected pool to close")
	}
}

func TestConnectClosesPoolWhenPingFails(t *testing.T) {
	fake := &fakePool{pingErr: errors.New("unavailable")}
	databaseConfig := config.Database{
		Addr:     "localhost:5432",
		Database: "acronom",
		User:     "postgres",
		Pass:     "secret",
	}

	_, err := connectWith(context.Background(), databaseConfig, func(context.Context, *pgxpool.Config) (pool, error) {
		return fake, nil
	})
	if err == nil {
		t.Fatal("expected connection error")
	}
	if !fake.closed {
		t.Fatal("expected failed pool to close")
	}
	if contains := errors.Is(err, fake.pingErr); contains {
		t.Fatal("expected internal connection error to be hidden")
	}
}

func TestConnectHidesConnectorError(t *testing.T) {
	internalError := errors.New("postgres://app:secret@localhost/acronom")
	databaseConfig := config.Database{
		Addr:     "localhost:5432",
		Database: "acronom",
		User:     "postgres",
		Pass:     "secret",
	}

	_, err := connectWith(context.Background(), databaseConfig, func(context.Context, *pgxpool.Config) (pool, error) {
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
