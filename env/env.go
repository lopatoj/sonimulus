// Package env provides configuration settings from environment variables for the sonimulus application.

package env

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"go-simpler.org/env"
)

// Env holds the configuration settings for the application.
type Env struct {
	Server struct {
		Port int    `env:"PORT" default:"8080"`
		URL  string `env:"URL" default:"https://localhost"`
	} `env:"SERVER_"`
	Soundcloud struct {
		URL          string `env:"URL"`
		ClientID     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
	} `env:"SOUNDCLOUD_"`
	DB struct {
		URI string `env:"URI"`
	} `env:"POSTGRES_"`
	JWT struct {
		Secret     string `env:"SECRET" default:"secret"`
		Expiration int    `env:"EXPIRE" default:"3600"`
	} `env:"JWT_"`
}

// NewEnv initializes a new Env instance, drawing from environment variables.
func NewEnv() (e Env, err error) {
	// Load environment variables from .env file
	if err = godotenv.Load(); err != nil {
		slog.Error("failed to load environment variables from .env file")
		return Env{}, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize fields from environment variables
	if err := env.Load(&e, nil); err != nil {
		slog.Error("failed to set config fields from environment variables")
		return Env{}, fmt.Errorf("failed to load config: %w", err)
	}

	return e, nil
}
