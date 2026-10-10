package model

import (
	"context"
	"errors"
	"testing"
)

func TestNormalizeCategoryName(t *testing.T) {
	if actual := NormalizeCategoryName("  platform   tools "); actual != "Platform Tools" {
		t.Fatalf("Unexpected normalized name: %q", actual)
	}
}

func TestApprovedCategoryColours(t *testing.T) {
	if !IsApprovedCategoryColour("#1976d2") {
		t.Fatal("Expected approved colour to be case-insensitive")
	}
	if IsApprovedCategoryColour("#FFFFFF") {
		t.Fatal("Expected unapproved colour to be rejected")
	}
}

func TestPGCategoryStore(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	ctx := context.Background()
	store := NewPGCategoryStore(pool)

	platform := &Category{Name: "  platform  ", Colour: "#1976d2"}
	if err := store.Create(ctx, platform); err != nil {
		t.Fatal(err)
	}
	if platform.Name != "Platform" || platform.Colour != "#1976D2" {
		t.Fatalf("Unexpected normalized category: %#v", platform)
	}
	if err := store.Create(ctx, &Category{Name: "PLATFORM", Colour: "#00796B"}); !errors.Is(err, ErrCategoryExists) {
		t.Fatalf("Expected duplicate category error, got %v", err)
	}
	if err := store.Create(ctx, &Category{Name: "Business", Colour: "#388E3C"}); err != nil {
		t.Fatal(err)
	}

	categories, err := store.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) != 2 || categories[0].Name != "Business" || categories[1].Name != "Platform" {
		t.Fatalf("Expected alphabetical categories, got %#v", categories)
	}
}
