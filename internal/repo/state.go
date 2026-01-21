package repo

import (
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type StateRepository struct {
	kv KVStorer
}

func NewStateRepository(kv KVStorer) *StateRepository {
	return &StateRepository{kv: kv}
}

func (sr *StateRepository) CreateState(ctx context.Context, verifier string) (state string, err error) {
	state = rand.Text()
	err = sr.kv.Set(ctx, "state:"+state, verifier, 10*time.Minute).Err()
	if err != nil {
		slog.Error("Error creating state", "error", err)
		return "", err
	}
	return state, nil
}

func (sr *StateRepository) GetVerifier(ctx context.Context, state string) (verifier string, found bool, err error) {
	verifier, err = sr.kv.GetDel(ctx, "state:"+state).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		return "", false, err
	}
	return verifier, true, nil
}
