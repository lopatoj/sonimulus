package repository

import "database/sql"

// User represents a row in the users table.
type User struct {
	Id    string // UUID primary key.
	Name  string // Name of the user.
	Email string // Email address of the user.
	Salt  []byte // Salt used for password hashing.
	Hash  []byte // Hashed password for the user.
}

// UsersRepository represents a repository for managing users.
type UsersRepository struct {
	db *sql.DB
}

// NewUsersRepository creates a new UsersRepository instance.
func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// GetUserById retrieves a user by ID.
func (ur *UsersRepository) GetUserById(id string) (*User, error) {
	var user User
	err := ur.db.QueryRow("SELECT id, name, email, salt, hash FROM users.users WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Email, &user.Salt, &user.Hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email.
func (ur *UsersRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := ur.db.QueryRow("SELECT id, name, email, salt, hash FROM users.users WHERE email = $1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Salt, &user.Hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByName retrieves a user by name.
func (ur *UsersRepository) GetUserByName(name string) (*User, error) {
	var user User
	err := ur.db.QueryRow("SELECT id, name, email, salt, hash FROM users.users WHERE name = $1", name).Scan(&user.Id, &user.Name, &user.Email, &user.Salt, &user.Hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create creates a new user, and returns the created user.
func (ur *UsersRepository) Create(name, email string, salt, hash []byte) (*User, error) {
	var user User
	err := ur.db.QueryRow("INSERT INTO users.users (name, email, salt, hash) VALUES ($1, $2, $3, $4) RETURNING *", name, email, salt, hash).Scan(&user.Id, &user.Name, &user.Email, &user.Salt, &user.Hash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
