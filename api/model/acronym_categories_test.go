package model

import (
	"context"
	"testing"
)

func TestPGAcronymStoreCategoryAssignments(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	ctx := context.Background()
	categories := NewPGCategoryStore(pool)
	acronyms := NewPGAcronymStore(pool)

	platform := &Category{Name: "Platform", Colour: "#1976D2"}
	business := &Category{Name: "Business", Colour: "#388E3C"}
	for _, category := range []*Category{platform, business} {
		if err := categories.Create(ctx, category); err != nil {
			t.Fatal(err)
		}
	}

	first := &Acronym{
		Acronym:     "API",
		Definition:  "Application Programming Interface",
		CategoryIDs: []int64{platform.ID, business.ID},
	}
	if err := acronyms.Create(ctx, first); err != nil {
		t.Fatal(err)
	}
	if len(first.Categories) != 2 || first.Categories[0].Name != "Business" || first.Categories[1].Name != "Platform" {
		t.Fatalf("Expected alphabetical assigned categories, got %#v", first.Categories)
	}

	second := &Acronym{
		Acronym:     "SDK",
		Definition:  "Software Development Kit",
		CategoryIDs: []int64{platform.ID},
	}
	if err := acronyms.Create(ctx, second); err != nil {
		t.Fatal(err)
	}

	first.CategoryIDs = []int64{}
	first.CategoriesSpecified = true
	if err := acronyms.Update(ctx, first); err != nil {
		t.Fatal(err)
	}
	if len(first.Categories) != 0 {
		t.Fatalf("Expected categories to be cleared, got %#v", first.Categories)
	}

	second.Definition = "Development tools"
	if err := acronyms.Update(ctx, second); err != nil {
		t.Fatal(err)
	}
	if len(second.Categories) != 1 || second.Categories[0].ID != platform.ID {
		t.Fatalf("Expected omitted categories to be preserved, got %#v", second.Categories)
	}

	if err := acronyms.Delete(ctx, second.ID); err != nil {
		t.Fatal(err)
	}
	available, err := categories.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(available) != 2 {
		t.Fatalf("Expected shared categories to remain, got %#v", available)
	}
}

func TestPGAcronymStoreRejectsUnknownCategory(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	store := NewPGAcronymStore(pool)

	acronym := &Acronym{
		Acronym:     "API",
		Definition:  "Application Programming Interface",
		CategoryIDs: []int64{999},
	}
	if err := store.Create(context.Background(), acronym); err != ErrCategoryNotFound {
		t.Fatalf("Expected unknown category error, got %v", err)
	}

	var count int
	if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM acronyms").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("Expected failed transaction to roll back, found %d acronyms", count)
	}
}
