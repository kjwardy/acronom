package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	setEnvironment(t, map[string]string{
		"HOSTS":   "localhost",
		"APP_URI": "http://localhost:1323",
	})

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
	if config.Hosts != "localhost" {
		t.Fatalf("expected hosts 'localhost', got '%s'", config.Hosts)
	}
	if config.AppURI != "http://localhost:1323" {
		t.Fatalf("expected app URI 'http://localhost:1323', got '%s'", config.AppURI)
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
		"HOSTS":             "example.com",
		"APP_URI":           "http://example.com",
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
	setEnvironment(t, map[string]string{
		"PORT":    "invalid",
		"HOSTS":   "localhost",
		"APP_URI": "http://localhost:1323",
	})

	if _, err := Load(); err == nil {
		t.Fatal("expected an invalid port error")
	}
}

func TestLoadRejectsInvalidDatabaseAddress(t *testing.T) {
	setEnvironment(t, map[string]string{
		"HOSTS":         "localhost",
		"APP_URI":       "http://localhost:1323",
		"POSTGRES_ADDR": "database",
	})

	if _, err := Load(); err == nil {
		t.Fatal("expected an invalid database address error")
	}
}

func TestLoadRejectsMissingHosts(t *testing.T) {
	setEnvironment(t, map[string]string{
		"APP_URI": "http://localhost:1323",
	})

	if _, err := Load(); err == nil {
		t.Fatal("expected an error for missing HOSTS")
	}
}

func TestLoadRejectsMissingAppURI(t *testing.T) {
	setEnvironment(t, map[string]string{
		"HOSTS": "localhost",
	})

	if _, err := Load(); err == nil {
		t.Fatal("expected an error for missing APP_URI")
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
		"HOSTS",
		"APP_URI",
		"POSTGRES_ADDR",
		"POSTGRES_DATABASE",
		"POSTGRES_USER",
		"POSTGRES_PASS",
	} {
		t.Setenv(name, values[name])
	}
}
