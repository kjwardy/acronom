package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kjwardy/acronom/api/config"
)

type pool interface {
	Ping(context.Context) error
	Close()
}

type connector func(context.Context, *pgxpool.Config) (pool, error)

type DB struct {
	pool pool
}

func Connect(ctx context.Context, databaseConfig config.Database) (*DB, error) {
	return connectWith(ctx, databaseConfig, func(ctx context.Context, config *pgxpool.Config) (pool, error) {
		return pgxpool.NewWithConfig(ctx, config)
	})
}

func connectWith(ctx context.Context, databaseConfig config.Database, connect connector) (*DB, error) {
	poolConfig, err := poolConfig(databaseConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid PostgreSQL connection configuration")
	}

	connection, err := connect(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("could not initialise PostgreSQL connection")
	}
	if err := connection.Ping(ctx); err != nil {
		connection.Close()
		return nil, fmt.Errorf("could not connect to PostgreSQL")
	}

	return &DB{pool: connection}, nil
}

func poolConfig(databaseConfig config.Database) (*pgxpool.Config, error) {
	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseConfig.User, databaseConfig.Pass),
		Host:   databaseConfig.Addr,
		Path:   databaseConfig.Database,
	}
	return pgxpool.ParseConfig(connectionURL.String())
}

func (database *DB) Ping(ctx context.Context) error {
	return database.pool.Ping(ctx)
}

func (database *DB) Close() {
	database.pool.Close()
}
