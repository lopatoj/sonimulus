package controllers

import (
	"bytes"
	"crypto/rand"
	"errors"
	"log/slog"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
	"lopa.to/sonimulus/config"
	"lopa.to/sonimulus/repository"
)

var (
	// ErrUserAlreadyExists is returned when a user with the given email already exists.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotFound is returned when a user with the given email does not exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials is returned when the provided email and password do not match.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidToken is returned when the provided token is invalid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidName is returned when the provided name is invalid.
	ErrInvalidName = errors.New("invalid name")

	// ErrInvalidEmail is returned when the provided email is invalid.
	ErrInvalidEmail = errors.New("invalid email")

	// ErrInvalidPassword is returned when the provided password is invalid.
	ErrInvalidPassword = errors.New("invalid password")
)

const (
	nameRegex  = `^[a-zA-Z0-9._%+-]{2,}$`
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

// AuthController handles authentication-related operations.
type AuthController struct {
	users  UsersRepository
	config config.Config
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController(users UsersRepository, config config.Config) *AuthController {
	return &AuthController{users: users, config: config}
}

func (ac *AuthController) Login(email, password string) (user *repository.User, token string, err error) {
	user, err = ac.users.GetUserByEmail(email)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, "", err
	}

	if user == nil {
		slog.Error("user not found", "error", err)
		return nil, "", ErrUserNotFound
	}

	if !verifyPassword(password, user) {
		slog.Error("failed to verify password")
		return nil, "", ErrInvalidCredentials
	}

	token, err = ac.refreshToken(user)
	if err != nil {
		slog.Error("failed to create token", "error", err)
		return nil, "", err
	}

	return user, token, nil
}

// Register creates a new user.
func (ac *AuthController) Register(name, email, password string) (user *repository.User, token string, err error) {
	if matched, err := regexp.MatchString(nameRegex, name); !matched || err != nil {
		slog.Error("invalid name")
		return nil, "", ErrInvalidName
	}

	if matched, err := regexp.MatchString(emailRegex, email); !matched || err != nil {
		slog.Error("invalid email")
		return nil, "", ErrInvalidEmail
	}

	if !isValidPassword(password) {
		slog.Error("invalid password")
		return nil, "", ErrInvalidPassword
	}

	hash, salt, err := hashPassword(password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, "", err
	}

	user, err = ac.users.GetUserByName(name)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, "", err
	}

	if user != nil {
		slog.Error("user already exists", "error", err)
		return nil, "", ErrUserAlreadyExists
	}

	user, err = ac.users.Create(name, email, salt, hash)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, "", err
	}

	token, err = ac.refreshToken(user)
	if err != nil {
		slog.Error("failed to create token", "error", err)
		return nil, "", err
	}

	return user, token, nil
}

func (ac *AuthController) ValidateToken(token string) (user *repository.User, err error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(ac.config.JWT.Secret), nil
	})

	if err != nil {
		slog.Error("failed to parse token", "error", err)
		return nil, ErrInvalidToken
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userId := claims["userId"].(string)
		user, err = ac.users.GetUserById(userId)
		if err != nil {
			slog.Error("failed to get user", "error", err)
			return nil, ErrUserNotFound
		}
	} else {
		return nil, ErrInvalidToken
	}

	return user, nil
}

// generateToken generates a JWT token for the given user.
func (ac *AuthController) refreshToken(user *repository.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"exp":    time.Now().Add(time.Duration(ac.config.JWT.Expiration) * time.Second).Unix(),
		"iat":    time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(ac.config.JWT.Secret))
	if err != nil {
		slog.Error("failed to sign token", "error", err)
		return "", err
	}

	return tokenString, nil
}

// hashPassword salts and hashes a password using Argon2.
func hashPassword(password string) (hash []byte, salt []byte, err error) {
	salt = make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}
	hash = argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hash, salt, nil
}

// verifyPassword verifies a password against a user's hashed password.
func verifyPassword(password string, user *repository.User) bool {
	hash := argon2.IDKey([]byte(password), user.Salt, 1, 64*1024, 4, 32)
	return bytes.Equal(hash, user.Hash)
}

// isValidPassword checks if a password meets the required criteria.
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	validChars := regexp.MustCompile(`^[A-Za-z\d@$!%*?&]+$`).MatchString(password)

	return hasLetter && hasDigit && validChars
}
