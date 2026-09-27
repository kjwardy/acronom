package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kjwardy/acronom/api/config"
)

type connector func(context.Context, *pgxpool.Config) (*pgxpool.Pool, error)

type DB struct {
	pool *pgxpool.Pool
}

func (database *DB) Pool() *pgxpool.Pool {
	return database.pool
}

func (database *DB) Ping(ctx context.Context) error {
	return database.pool.Ping(ctx)
}

func Connect(ctx context.Context, databaseConfig config.Database) (*DB, error) {
	return connectWith(ctx, databaseConfig, func(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
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

func (database *DB) Close() {
	database.pool.Close()
}

// CreateSchema creates the acronyms table and indexes if they don't exist
// This is called during application startup to ensure the database schema is up to date
func (database *DB) CreateSchema(ctx context.Context) error {
	// Create the acronyms table with all required fields and constraints
	_, err := database.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS acronyms (
			id BIGSERIAL PRIMARY KEY,
			acronym TEXT NOT NULL UNIQUE,
			definition TEXT NOT NULL,
			link TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create acronyms table: %w", err)
	}

	// Create an index on the lowercase version of acronym for case-insensitive lookups
	_, err = database.pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_acronyms_acronym ON acronyms(LOWER(acronym))
	`)
	if err != nil {
		return fmt.Errorf("failed to create acronym index: %w", err)
	}

	return nil
}
