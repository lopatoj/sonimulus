package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserKey string

const (
	UserKeyID       UserKey = "id"
	UserKeyPersonID UserKey = "person_id"
)

// User represents a row in the users table.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	PersonID  *int64    `json:"person_id"`
	Processed bool      `json:"processed"`
}

// UsersRepository represents a repository for managing users.
type UsersRepository struct {
	db *sql.DB
}

// NewUsersRepository creates a new UsersRepository instance.
func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// FindByKey retrieves a user by key.
func (ur *UsersRepository) FindByKey(ctx context.Context, key UserKey, value string) (user User, found bool, err error) {
	if key != UserKeyID && key != UserKeyPersonID {
		return user, false, errors.New("invalid user key")
	}
	query := fmt.Sprintf("SELECT id, created_at, person_id, processed, username FROM users WHERE %s = $1;", key)
	return ur.queryUserRow(ctx, query, value)
}

// Create creates a new user, and returns the created user.
func (ur *UsersRepository) Create(ctx context.Context, id int64, username string) (user User, err error) {
	user, _, err = ur.queryUserRow(ctx, "INSERT INTO users (id, username) VALUES ($1, $2) RETURNING id, created_at, person_id, processed, username;", id, username)
	return user, err
}

func (ur *UsersRepository) queryUserRow(ctx context.Context, query string, args ...any) (user User, found bool, err error) {
	err = ur.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.PersonID, &user.Processed, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, false, nil
		}
		return user, false, err
	}
	return user, true, nil
}
