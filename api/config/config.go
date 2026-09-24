package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Debug bool
	Port  int
}

func Load() (Config, error) {
	config := Config{Port: 1323}

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

	return config, nil
}
