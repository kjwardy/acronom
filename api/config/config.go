package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Database struct {
	Addr     string
	Database string
	User     string
	Pass     string
}

type Config struct {
	Debug    bool
	Port     int
	Hosts    string
	AppURI   string
	Database Database
}

func Load() (Config, error) {
	config := Config{
		Port: 1323,
		Database: Database{
			Addr:     envOrDefault("POSTGRES_ADDR", "localhost:5432"),
			Database: envOrDefault("POSTGRES_DATABASE", "acronom"),
			User:     envOrDefault("POSTGRES_USER", "postgres"),
			Pass:     envOrDefault("POSTGRES_PASS", "password"),
		},
		Hosts:  envOrDefault("HOSTS", ""),
		AppURI: envOrDefault("APP_URI", ""),
	}

	if value := os.Getenv("PORT"); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return Config{}, fmt.Errorf("PORT must be a number between 1 and 65535")
		}
		config.Port = port
	}

	if value := os.Getenv("DEBUG"); value != "" {
		debug, err := strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("DEBUG must be a boolean")
		}
		config.Debug = debug
	}

	if config.Hosts == "" {
		return Config{}, fmt.Errorf("HOSTS is required")
	}

	if config.AppURI == "" {
		return Config{}, fmt.Errorf("APP_URI is required")
	}

	if err := validateDatabase(config.Database); err != nil {
		return Config{}, err
	}

	return config, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func validateDatabase(config Database) error {
	host, portValue, err := net.SplitHostPort(config.Addr)
	if err != nil || host == "" {
		return fmt.Errorf("POSTGRES_ADDR must use the host:port format")
	}
	port, err := strconv.Atoi(portValue)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("POSTGRES_ADDR must contain a valid port")
	}
	if config.Database == "" {
		return fmt.Errorf("POSTGRES_DATABASE must not be empty")
	}
	if config.User == "" {
		return fmt.Errorf("POSTGRES_USER must not be empty")
	}
	return nil
}
