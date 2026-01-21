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
		Port          int    `env:"PUBLIC_API_PORT" default:"3000"`
		URL           string `env:"PUBLIC_API_URL" default:"https://localhost"`
		Route         string `env:"PUBLIC_API_ROUTE" default:"/api/v1"`
		RedirectRoute string `env:"API_REDIRECT_ROUTE" default:"/auth/callback"`
	}
	Client struct {
		URL  string `env:"URL" default:"http://localhost"`
		Port int    `env:"PORT" default:"5173"`
	} `env:"CLIENT_"`
	Soundcloud struct {
		URL          string `env:"URL" default:"https://soundcloud.com"`
		APIURL       string `env:"API_URL" default:"https://api.soundcloud.com"`
		ClientID     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
		AuthURL      string `env:"AUTH_URL"`
		TokenURL     string `env:"TOKEN_URL"`
	} `env:"SOUNDCLOUD_"`
	DB struct {
		PostgresURI string `env:"POSTGRES_URI"`
		RedisURI    string `env:"REDIS_URI"`
	}
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
