package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const MaxCategoriesPerAcronym = 5
const MaxCategoryNameLength = 30

var ApprovedCategoryColours = []string{
	"#1976D2",
	"#00796B",
	"#388E3C",
	"#FFA000",
	"#F57C00",
	"#D32F2F",
	"#7B1FA2",
	"#455A64",
}

var ErrCategoryExists = errors.New("category already exists")
var ErrCategoryNotFound = errors.New("category not found")

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Colour    string    `json:"colour"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CategoryStore interface {
	List(context.Context) ([]Category, error)
	Create(context.Context, *Category) error
}

type PGCategoryStore struct {
	pool *pgxpool.Pool
}

func NewPGCategoryStore(pool *pgxpool.Pool) *PGCategoryStore {
	return &PGCategoryStore{pool: pool}
}

func NormalizeCategoryName(name string) string {
	return cases.Title(language.English).String(strings.Join(strings.Fields(name), " "))
}

func IsApprovedCategoryColour(colour string) bool {
	colour = strings.ToUpper(strings.TrimSpace(colour))
	for _, approved := range ApprovedCategoryColours {
		if colour == approved {
			return true
		}
	}
	return false
}

func (store *PGCategoryStore) List(ctx context.Context) ([]Category, error) {
	rows, err := store.pool.Query(ctx, `
		SELECT id, name, colour, created_at, updated_at
		FROM categories
		ORDER BY LOWER(name), id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Colour,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to read category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read categories: %w", err)
	}
	return categories, nil
}

func (store *PGCategoryStore) Create(ctx context.Context, category *Category) error {
	category.Name = NormalizeCategoryName(category.Name)
	category.Colour = strings.ToUpper(strings.TrimSpace(category.Colour))
	err := store.pool.QueryRow(ctx, `
		INSERT INTO categories (name, colour)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`, category.Name, category.Colour).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) &&
			postgresError.Code == "23505" &&
			postgresError.ConstraintName == "idx_categories_name" {
			return ErrCategoryExists
		}
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}
