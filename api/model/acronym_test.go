package model

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kjwardy/acronom/api/config"
	"github.com/kjwardy/acronom/api/db"
)

// setupTestDB creates a test database connection and truncates the acronyms table
// This ensures each test starts with a clean slate
func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	// Use environment variables for database connection if available
	postgresAddr := os.Getenv("POSTGRES_ADDR")
	if postgresAddr == "" {
		postgresAddr = "localhost:5432"
	}

	postgresUser := os.Getenv("POSTGRES_USER")
	if postgresUser == "" {
		postgresUser = "postgres"
	}

	postgresPass := os.Getenv("POSTGRES_PASS")
	if postgresPass == "" {
		postgresPass = "password"
	}

	postgresDB := os.Getenv("POSTGRES_DATABASE")
	if postgresDB == "" {
		postgresDB = "acronom"
	}

	database, err := db.Connect(ctx, config.Database{
		Addr:     postgresAddr,
		Database: postgresDB,
		User:     postgresUser,
		Pass:     postgresPass,
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	pool := database.Pool()

	if err := database.CreateSchema(ctx); err != nil {
		pool.Close()
		t.Fatalf("Failed to create test schema: %v", err)
	}

	// Clean up any existing data before running tests
	_, err = pool.Exec(ctx, "TRUNCATE TABLE acronyms RESTART IDENTITY CASCADE")
	if err != nil {
		pool.Close()
		t.Fatalf("Failed to truncate test table: %v", err)
	}

	return pool
}

func TestPGAcronymStore_Create(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	store := NewPGAcronymStore(pool)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		acronym := &Acronym{
			Acronym:    "API",
			Definition: "Application Programming Interface",
		}

		err := store.Create(ctx, acronym)
		if err != nil {
			t.Errorf("Create() failed: %v", err)
		}
		if acronym.ID == 0 {
			t.Error("Expected non-zero ID after creation")
		}
		if acronym.CreatedAt.IsZero() {
			t.Error("Expected non-zero CreatedAt after creation")
		}
		if acronym.UpdatedAt.IsZero() {
			t.Error("Expected non-zero UpdatedAt after creation")
		}
	})

	t.Run("duplicate acronym", func(t *testing.T) {
		acronym1 := &Acronym{
			Acronym:    "REST",
			Definition: "Representational State Transfer",
		}
		err := store.Create(ctx, acronym1)
		if err != nil {
			t.Fatalf("Failed to create first acronym: %v", err)
		}

		acronym2 := &Acronym{
			Acronym:    "REST",
			Definition: "Different definition",
		}
		err = store.Create(ctx, acronym2)
		if err == nil {
			t.Error("Expected error when creating duplicate acronym")
		}
		if err != nil && err.Error() != "acronym already exists" {
			t.Errorf("Expected 'acronym already exists' error, got: %v", err)
		}
	})
}

func TestPGAcronymStore_Find(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	store := NewPGAcronymStore(pool)
	ctx := context.Background()

	created := &Acronym{
		Acronym:    "JSON",
		Definition: "JavaScript Object Notation",
	}
	err := store.Create(ctx, created)
	if err != nil {
		t.Fatalf("Failed to create test acronym: %v", err)
	}

	t.Run("find existing acronym", func(t *testing.T) {
		found, err := store.Find(ctx, created.ID)
		if err != nil {
			t.Errorf("Find() failed: %v", err)
		}
		if found == nil {
			t.Error("Expected to find acronym, got nil")
		}
		if found != nil {
			if found.ID != created.ID {
				t.Errorf("Expected ID %d, got %d", created.ID, found.ID)
			}
			if found.Acronym != "JSON" {
				t.Errorf("Expected acronym 'JSON', got '%s'", found.Acronym)
			}
			if found.Definition != "JavaScript Object Notation" {
				t.Errorf("Expected definition 'JavaScript Object Notation', got '%s'", found.Definition)
			}
		}
	})

	t.Run("find non-existent acronym", func(t *testing.T) {
		found, err := store.Find(ctx, 99999)
		if err != nil {
			t.Errorf("Find() failed for non-existent ID: %v", err)
		}
		if found != nil {
			t.Error("Expected nil for non-existent acronym, got value")
		}
	})
}

func TestPGAcronymStore_Update(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	store := NewPGAcronymStore(pool)
	ctx := context.Background()

	created := &Acronym{
		Acronym:    "HTTP",
		Definition: "HyperText Transfer Protocol",
	}
	err := store.Create(ctx, created)
	if err != nil {
		t.Fatalf("Failed to create test acronym: %v", err)
	}

	t.Run("successful update", func(t *testing.T) {
		newDefinition := "HyperText Transfer Protocol - The foundation of data communication"
		link := "https://en.wikipedia.org/wiki/HTTP"
		created.Definition = newDefinition
		created.Link = &link

		err := store.Update(ctx, created)
		if err != nil {
			t.Errorf("Update() failed: %v", err)
		}

		updated, err := store.Find(ctx, created.ID)
		if err != nil {
			t.Fatalf("Failed to find updated acronym: %v", err)
		}
		if updated.Definition != newDefinition {
			t.Errorf("Expected definition '%s', got '%s'", newDefinition, updated.Definition)
		}
		if updated.Link == nil {
			t.Error("Link was not set")
		}
		if updated.Link != nil && *updated.Link != "https://en.wikipedia.org/wiki/HTTP" {
			t.Errorf("Expected link 'https://en.wikipedia.org/wiki/HTTP', got '%s'", *updated.Link)
		}
		if !updated.UpdatedAt.After(created.CreatedAt) {
			t.Error("UpdatedAt should be after CreatedAt")
		}
	})

	t.Run("update non-existent acronym", func(t *testing.T) {
		nonExistent := &Acronym{
			ID:         99999,
			Definition: "New definition",
		}
		err := store.Update(ctx, nonExistent)
		if err == nil {
			t.Error("Expected error when updating non-existent acronym")
		}
		if err != nil && err.Error() != "acronym not found" {
			t.Errorf("Expected 'acronym not found' error, got: %v", err)
		}
	})
}

func TestPGAcronymStore_Delete(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	store := NewPGAcronymStore(pool)
	ctx := context.Background()

	created := &Acronym{
		Acronym:    "SQL",
		Definition: "Structured Query Language",
	}
	err := store.Create(ctx, created)
	if err != nil {
		t.Fatalf("Failed to create test acronym: %v", err)
	}

	t.Run("successful deletion", func(t *testing.T) {
		err := store.Delete(ctx, created.ID)
		if err != nil {
			t.Errorf("Delete() failed: %v", err)
		}

		found, err := store.Find(ctx, created.ID)
		if err != nil {
			t.Errorf("Find() failed after deletion: %v", err)
		}
		if found != nil {
			t.Error("Expected nil after deletion, got value")
		}
	})

	t.Run("delete non-existent acronym", func(t *testing.T) {
		err := store.Delete(ctx, 99999)
		if err == nil {
			t.Error("Expected error when deleting non-existent acronym")
		}
		if err != nil && err.Error() != "acronym not found" {
			t.Errorf("Expected 'acronym not found' error, got: %v", err)
		}
	})
}

func TestPGAcronymStore_FindByAcronym(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	store := NewPGAcronymStore(pool)
	ctx := context.Background()

	created := &Acronym{
		Acronym:    "XML",
		Definition: "eXtensible Markup Language",
	}
	err := store.Create(ctx, created)
	if err != nil {
		t.Fatalf("Failed to create test acronym: %v", err)
	}

	t.Run("find by exact acronym", func(t *testing.T) {
		found, err := store.FindByAcronym(ctx, "XML")
		if err != nil {
			t.Errorf("FindByAcronym() failed: %v", err)
		}
		if found == nil {
			t.Error("Expected to find acronym, got nil")
		}
		if found != nil && found.ID != created.ID {
			t.Errorf("Expected ID %d, got %d", created.ID, found.ID)
		}
	})

	t.Run("find by case-insensitive acronym", func(t *testing.T) {
		found, err := store.FindByAcronym(ctx, "xml")
		if err != nil {
			t.Errorf("FindByAcronym() failed with lowercase: %v", err)
		}
		if found == nil {
			t.Error("Expected to find acronym with lowercase, got nil")
		}
		if found != nil && found.ID != created.ID {
			t.Errorf("Expected ID %d with lowercase, got %d", created.ID, found.ID)
		}
	})

	t.Run("find by acronym with spaces", func(t *testing.T) {
		found, err := store.FindByAcronym(ctx, "  XML  ")
		if err != nil {
			t.Errorf("FindByAcronym() failed with spaces: %v", err)
		}
		if found == nil {
			t.Error("Expected to find acronym with spaces, got nil")
		}
		if found != nil && found.ID != created.ID {
			t.Errorf("Expected ID %d with spaces, got %d", created.ID, found.ID)
		}
	})

	t.Run("find non-existent acronym", func(t *testing.T) {
		found, err := store.FindByAcronym(ctx, "NONEXISTENT")
		if err != nil {
			t.Errorf("FindByAcronym() failed for non-existent: %v", err)
		}
		if found != nil {
			t.Error("Expected nil for non-existent acronym, got value")
		}
	})
}
