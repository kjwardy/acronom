package model

import (
	"context"
	"testing"
)

func TestPGMetricsStoreTotalAcronyms(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	ctx := context.Background()
	metrics := NewPGMetricsStore(pool)
	acronyms := NewPGAcronymStore(pool)

	total, err := metrics.TotalAcronyms(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if total != 0 {
		t.Fatalf("Expected empty database count 0, got %d", total)
	}

	entries := []*Acronym{
		{Acronym: "PC", Definition: "Personal Computer"},
		{Acronym: "PC", Definition: "Probable Cause"},
	}
	for _, entry := range entries {
		if err := acronyms.Create(ctx, entry); err != nil {
			t.Fatalf("Failed to create metric fixture: %v", err)
		}
	}

	total, err = metrics.TotalAcronyms(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Fatalf("Expected each glossary entry to be counted, got %d", total)
	}
}
