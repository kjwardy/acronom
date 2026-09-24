package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DEBUG", "")

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
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DEBUG", "true")

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
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	t.Setenv("PORT", "invalid")

	if _, err := Load(); err == nil {
		t.Fatal("expected an invalid port error")
	}
}
