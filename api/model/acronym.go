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
	ID                  int64      `json:"id"`             // Unique identifier
	Acronym             string     `json:"acronym"`        // The acronym text (e.g., "API")
	Definition          string     `json:"definition"`     // The definition of the acronym
	Link                *string    `json:"link,omitempty"` // Optional link to external resource
	CreatedAt           time.Time  `json:"created_at"`     // Timestamp when the acronym was created
	UpdatedAt           time.Time  `json:"updated_at"`     // Timestamp when the acronym was last updated
	Categories          []Category `json:"categories"`
	CategoryIDs         []int64    `json:"-"`
	CategoriesSpecified bool       `json:"-"`
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
var ErrTooManyCategories = errors.New("too many categories")

func isAcronymMeaningConflict(err error) bool {
	var postgresError *pgconn.PgError
	return errors.As(err, &postgresError) &&
		postgresError.Code == "23505" &&
		postgresError.ConstraintName == "idx_acronyms_unique_meaning"
}

type categoryQuerier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func loadAcronymCategories(ctx context.Context, query categoryQuerier, acronyms []*Acronym) error {
	ids := make([]int64, 0, len(acronyms))
	byID := make(map[int64]*Acronym, len(acronyms))
	for _, acronym := range acronyms {
		acronym.Categories = make([]Category, 0)
		ids = append(ids, acronym.ID)
		byID[acronym.ID] = acronym
	}
	if len(ids) == 0 {
		return nil
	}

	rows, err := query.Query(ctx, `
		SELECT ac.acronym_id, c.id, c.name, c.colour, c.created_at, c.updated_at
		FROM acronym_categories ac
		JOIN categories c ON c.id = ac.category_id
		WHERE ac.acronym_id = ANY($1)
		ORDER BY LOWER(c.name), c.id
	`, ids)
	if err != nil {
		return fmt.Errorf("failed to load acronym categories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var acronymID int64
		var category Category
		if err := rows.Scan(
			&acronymID,
			&category.ID,
			&category.Name,
			&category.Colour,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return fmt.Errorf("failed to read acronym category: %w", err)
		}
		if acronym := byID[acronymID]; acronym != nil {
			acronym.Categories = append(acronym.Categories, category)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to read acronym categories: %w", err)
	}
	return nil
}

func assignAcronymCategories(ctx context.Context, tx pgx.Tx, acronymID int64, categoryIDs []int64, replace bool) error {
	if len(categoryIDs) > MaxCategoriesPerAcronym {
		return ErrTooManyCategories
	}
	if len(categoryIDs) > 0 {
		var count int
		if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM categories WHERE id = ANY($1)", categoryIDs).Scan(&count); err != nil {
			return fmt.Errorf("failed to validate categories: %w", err)
		}
		if count != len(categoryIDs) {
			return ErrCategoryNotFound
		}
	}
	if replace {
		if _, err := tx.Exec(ctx, "DELETE FROM acronym_categories WHERE acronym_id = $1", acronymID); err != nil {
			return fmt.Errorf("failed to replace acronym categories: %w", err)
		}
	}
	for _, categoryID := range categoryIDs {
		if _, err := tx.Exec(ctx, `
			INSERT INTO acronym_categories (acronym_id, category_id)
			VALUES ($1, $2)
		`, acronymID, categoryID); err != nil {
			return fmt.Errorf("failed to assign acronym category: %w", err)
		}
	}
	return nil
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
	if err := loadAcronymCategories(ctx, s.pool, []*Acronym{acronym}); err != nil {
		return nil, err
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
	if err := loadAcronymCategories(ctx, s.pool, []*Acronym{acronymModel}); err != nil {
		return nil, err
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
	acronyms := make([]*Acronym, len(results))
	for index := range results {
		acronyms[index] = &results[index]
	}
	if err := loadAcronymCategories(ctx, s.pool, acronyms); err != nil {
		return nil, err
	}
	return results, nil
}

// Create inserts a new acronym into the database
func (s *PGAcronymStore) Create(ctx context.Context, acronym *Acronym) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin acronym creation: %w", err)
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	acronym.CreatedAt = now
	acronym.UpdatedAt = now
	err = tx.QueryRow(ctx,
		`INSERT INTO acronyms (acronym, definition, link, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		strings.TrimSpace(acronym.Acronym),
		strings.TrimSpace(acronym.Definition),
		acronym.Link,
		now,
		now).Scan(&acronym.ID)
	if err != nil {
		if isAcronymMeaningConflict(err) {
			return ErrAcronymMeaningExists
		}
		return fmt.Errorf("failed to create acronym: %w", err)
	}
	if err := assignAcronymCategories(ctx, tx, acronym.ID, acronym.CategoryIDs, false); err != nil {
		return err
	}
	if err := loadAcronymCategories(ctx, tx, []*Acronym{acronym}); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit acronym creation: %w", err)
	}
	return nil
}

// Update modifies an existing acronym's definition and link
// Returns an error if the acronym doesn't exist (no rows affected)
func (s *PGAcronymStore) Update(ctx context.Context, acronym *Acronym) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin acronym update: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
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
	if acronym.CategoriesSpecified {
		if err := assignAcronymCategories(ctx, tx, acronym.ID, acronym.CategoryIDs, true); err != nil {
			return err
		}
	}
	if err := loadAcronymCategories(ctx, tx, []*Acronym{acronym}); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit acronym update: %w", err)
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
