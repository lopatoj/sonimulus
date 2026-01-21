package repo

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

// SessionRepository is a repository providing methods for modifying user session state.
type SessionRepository struct {
	kv KVStorer
}

// Session is the type of data stored for an individual user session.
type Session struct {
	Token *oauth2.Token `json:"token,omitempty"`
	User  User          `json:"data"`
}

// MarshalBinary is a method for marshaling SessionData to JSON.
func (td Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(td)
}

// UnmarshalBinary is a method for unmarshaling SessionData from JSON.
func (td *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, td)
}

// NewSessionRepository creates a new SessionRepository.
func NewSessionRepository(kv KVStorer) *SessionRepository {
	return &SessionRepository{
		kv: kv,
	}
}

// CreateSession creates a new session.
func (sr *SessionRepository) CreateSession(ctx context.Context, data Session) (sessionID string, err error) {
	sessionID = rand.Text()

	err = sr.kv.Set(ctx, "session:"+sessionID, data, time.Until(data.Token.Expiry)).Err()
	if err != nil {
		slog.Error("Error creating session", "error", err)
		return "", err
	}

	return sessionID, nil
}

// GetSession retrieves a session by its ID.
func (sr *SessionRepository) GetSession(ctx context.Context, sessionID string) (session Session, found bool, err error) {
	err = sr.kv.Get(ctx, "session:"+sessionID).Scan(&session)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return session, false, nil
		}

		slog.Error("Error getting session", "error", err)

		return session, false, err
	}

	return session, true, nil
}

func (sr *SessionRepository) DeleteSession(ctx context.Context, sessionID string) (bool, error) {
	err := sr.kv.Del(ctx, "session:"+sessionID).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		slog.Error("Error deleting session", "error", err)

		return false, err
	}

	return true, nil
}
