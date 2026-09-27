package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Acronym represents an acronym entry in the database
type Acronym struct {
	ID         int64     `json:"id"`             // Unique identifier
	Acronym    string    `json:"acronym"`        // The acronym text (e.g., "API")
	Definition string    `json:"definition"`     // The definition of the acronym
	Link       *string   `json:"link,omitempty"` // Optional link to external resource
	CreatedAt  time.Time `json:"created_at"`     // Timestamp when the acronym was created
	UpdatedAt  time.Time `json:"updated_at"`     // Timestamp when the acronym was last updated
}

// AcronymStore defines the interface for acronym persistence operations
// This allows for easy testing with mock implementations
type AcronymStore interface {
	Find(ctx context.Context, id int64) (*Acronym, error)
	FindByAcronym(ctx context.Context, acronym string) (*Acronym, error)
	Search(ctx context.Context, query string) ([]Acronym, error)
	Create(ctx context.Context, acronym *Acronym) error
	Update(ctx context.Context, acronym *Acronym) error
	Delete(ctx context.Context, id int64) error
}

// PGAcronymStore implements AcronymStore using PostgreSQL with pgx
type PGAcronymStore struct {
	pool *pgxpool.Pool // PostgreSQL connection pool
}

// NewPGAcronymStore creates a new PGAcronymStore with the given connection pool
func NewPGAcronymStore(pool *pgxpool.Pool) *PGAcronymStore {
	return &PGAcronymStore{pool: pool}
}

const SearchResultLimit = 50

var ErrAcronymMeaningExists = errors.New("acronym meaning already exists")

func isAcronymMeaningConflict(err error) bool {
	var postgresError *pgconn.PgError
	return errors.As(err, &postgresError) &&
		postgresError.Code == "23505" &&
		postgresError.ConstraintName == "idx_acronyms_unique_meaning"
}

// Find retrieves an acronym by its ID
// Returns (nil, nil) if the acronym doesn't exist (not an error)
func (s *PGAcronymStore) Find(ctx context.Context, id int64) (*Acronym, error) {
	acronym := &Acronym{}
	err := s.pool.QueryRow(ctx,
		"SELECT id, acronym, definition, link, created_at, updated_at FROM acronyms WHERE id = $1",
		id).Scan(
		&acronym.ID, &acronym.Acronym, &acronym.Definition, &acronym.Link, &acronym.CreatedAt, &acronym.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is not an error
		}
		return nil, fmt.Errorf("failed to find acronym: %w", err)
	}
	return acronym, nil
}

// FindByAcronym retrieves an acronym by its text (case-insensitive)
// Returns (nil, nil) if the acronym doesn't exist (not an error)
func (s *PGAcronymStore) FindByAcronym(ctx context.Context, acronym string) (*Acronym, error) {
	acronymModel := &Acronym{}
	err := s.pool.QueryRow(ctx,
		"SELECT id, acronym, definition, link, created_at, updated_at FROM acronyms WHERE LOWER(acronym) = LOWER($1)",
		strings.TrimSpace(acronym)).Scan(
		&acronymModel.ID, &acronymModel.Acronym, &acronymModel.Definition, &acronymModel.Link, &acronymModel.CreatedAt, &acronymModel.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is not an error
		}
		return nil, fmt.Errorf("failed to find acronym by name: %w", err)
	}
	return acronymModel, nil
}

func (s *PGAcronymStore) Search(ctx context.Context, query string) ([]Acronym, error) {
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	escapedQuery := strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	).Replace(normalizedQuery)
	prefixPattern := escapedQuery + "%"
	containsPattern := "%" + escapedQuery + "%"

	rows, err := s.pool.Query(ctx, `
		SELECT id, acronym, definition, link, created_at, updated_at
		FROM acronyms
		WHERE LOWER(acronym) LIKE $3 ESCAPE '\'
		   OR LOWER(definition) LIKE $3 ESCAPE '\'
		ORDER BY CASE
			WHEN LOWER(acronym) = $1 THEN 0
			WHEN LOWER(acronym) LIKE $2 ESCAPE '\' THEN 1
			ELSE 2
		END, LOWER(acronym), id
		LIMIT $4
	`, normalizedQuery, prefixPattern, containsPattern, SearchResultLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to search acronyms: %w", err)
	}
	defer rows.Close()

	results := make([]Acronym, 0)
	for rows.Next() {
		var acronym Acronym
		if err := rows.Scan(
			&acronym.ID,
			&acronym.Acronym,
			&acronym.Definition,
			&acronym.Link,
			&acronym.CreatedAt,
			&acronym.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to read acronym search result: %w", err)
		}
		results = append(results, acronym)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read acronym search results: %w", err)
	}
	return results, nil
}

// Create inserts a new acronym into the database
func (s *PGAcronymStore) Create(ctx context.Context, acronym *Acronym) error {
	now := time.Now()
	acronym.CreatedAt = now
	acronym.UpdatedAt = now

	err := s.pool.QueryRow(ctx,
		`INSERT INTO acronyms (acronym, definition, link, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		strings.TrimSpace(acronym.Acronym),    // Trim whitespace from acronym
		strings.TrimSpace(acronym.Definition), // Trim whitespace from definition
		acronym.Link,
		now,
		now).Scan(&acronym.ID)

	if err != nil {
		if isAcronymMeaningConflict(err) {
			return ErrAcronymMeaningExists
		}
		return fmt.Errorf("failed to create acronym: %w", err)
	}
	return nil
}

// Update modifies an existing acronym's definition and link
// Returns an error if the acronym doesn't exist (no rows affected)
func (s *PGAcronymStore) Update(ctx context.Context, acronym *Acronym) error {
	err := s.pool.QueryRow(ctx,
		`UPDATE acronyms
		 SET definition = $1, link = $2, updated_at = NOW()
		 WHERE id = $3
		 RETURNING updated_at`,
		strings.TrimSpace(acronym.Definition),
		acronym.Link,
		acronym.ID).Scan(&acronym.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("acronym not found")
		}
		if isAcronymMeaningConflict(err) {
			return ErrAcronymMeaningExists
		}
		return fmt.Errorf("failed to update acronym: %w", err)
	}

	return nil
}

// Delete removes an acronym from the database by its ID
// Returns an error if the acronym doesn't exist (no rows affected)
func (s *PGAcronymStore) Delete(ctx context.Context, id int64) error {
	result, err := s.pool.Exec(ctx, "DELETE FROM acronyms WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete acronym: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("acronym not found")
	}

	return nil
}
