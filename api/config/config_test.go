package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	setEnvironment(t, map[string]string{})

	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if config.Port != 1323 {
		t.Fatalf("expected port 1323, got %d", config.Port)
	}
	if config.Debug {
		t.Fatal("expected debug to be disabled")
	}
	expected := Database{
		Addr:     "localhost:5432",
		Database: "acronom",
		User:     "postgres",
		Pass:     "password",
	}
	if config.Database != expected {
		t.Fatalf("unexpected database config: %#v", config.Database)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	setEnvironment(t, map[string]string{
		"PORT":              "8080",
		"DEBUG":             "true",
		"POSTGRES_ADDR":     "database:5433",
		"POSTGRES_DATABASE": "acronom_test",
		"POSTGRES_USER":     "acronom",
		"POSTGRES_PASS":     "secret",
	})

	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if config.Port != 8080 {
		t.Fatalf("expected port 8080, got %d", config.Port)
	}
	if !config.Debug {
		t.Fatal("expected debug to be enabled")
	}
	expected := Database{
		Addr:     "database:5433",
		Database: "acronom_test",
		User:     "acronom",
		Pass:     "secret",
	}
	if config.Database != expected {
		t.Fatalf("unexpected database config: %#v", config.Database)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	setEnvironment(t, map[string]string{"PORT": "invalid"})

	if _, err := Load(); err == nil {
		t.Fatal("expected an invalid port error")
	}
}

func TestLoadRejectsInvalidDatabaseAddress(t *testing.T) {
	setEnvironment(t, map[string]string{"POSTGRES_ADDR": "database"})

	if _, err := Load(); err == nil {
		t.Fatal("expected an invalid database address error")
	}
}

func TestValidateDatabaseRejectsMissingValues(t *testing.T) {
	tests := []Database{
		{Addr: "localhost:5432", User: "postgres"},
		{Addr: "localhost:5432", Database: "acronom"},
	}
	for _, database := range tests {
		if err := validateDatabase(database); err == nil {
			t.Fatalf("expected config to be rejected: %#v", database)
		}
	}
}

func setEnvironment(t *testing.T, values map[string]string) {
	t.Helper()
	for _, name := range []string{
		"PORT",
		"DEBUG",
		"POSTGRES_ADDR",
		"POSTGRES_DATABASE",
		"POSTGRES_USER",
		"POSTGRES_PASS",
	} {
		t.Setenv(name, values[name])
	}
}
