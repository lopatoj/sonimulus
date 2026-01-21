package data

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client.
func NewRedisClient(uri string) (*redis.Client, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	slog.Info("Redis connection established")

	return redis.NewClient(opt), nil
}
