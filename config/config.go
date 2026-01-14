// Package config provides configuration settings from environment variables for the sonimulus application.

package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"go-simpler.org/env"
)

// Config represents the configuration settings for the sonimulus application.
type Config struct {
	Port          int    `env:"PORT" default:"8080"`
	SoundcloudUrl string `env:"SOUNDCLOUD_URL"`
	DBUrl         string `env:"POSTGRES_URI"`
	JWT           struct {
		Secret     string `env:"SECRET" default:"secret"`
		Expiration int    `env:"EXPIRE" default:"3600"`
	} `env:"JWT_"`
}

// NewConfig initializes a new Config instance, drawing from environment variables and command line flags.
func NewConfig() (Config, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load environment variables from .env file")
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize fields from environment variables
	config := Config{}
	if err := env.Load(&config, nil); err != nil {
		slog.Error("failed to set config fields from environment variables")
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return config, nil
}
