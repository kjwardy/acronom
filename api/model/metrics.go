package model

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsStore interface {
	TotalAcronyms(context.Context) (int64, error)
}

type PGMetricsStore struct {
	pool *pgxpool.Pool
}

func NewPGMetricsStore(pool *pgxpool.Pool) *PGMetricsStore {
	return &PGMetricsStore{pool: pool}
}

func (store *PGMetricsStore) TotalAcronyms(ctx context.Context) (int64, error) {
	var total int64
	if err := store.pool.QueryRow(ctx, "SELECT COUNT(*) FROM acronyms").Scan(&total); err != nil {
		return 0, fmt.Errorf("failed to count acronyms: %w", err)
	}
	return total, nil
}
