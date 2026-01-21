package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
)

type KVGetter interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}

type KVSetter interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}

type KVDeleter interface {
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type KVGetAndDeleter interface {
	GetDel(ctx context.Context, key string) *redis.StringCmd
}

type KVStorer interface {
	KVGetAndDeleter
	KVDeleter
	KVSetter
	KVGetter
}

type DBQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
