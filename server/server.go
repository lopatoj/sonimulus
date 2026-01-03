package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"lopa.to/sonimulus/api"
	"lopa.to/sonimulus/config"
	"lopa.to/sonimulus/controllers"
)

type Server struct {
	auth   *controllers.AuthController
	users  *controllers.UsersController
	config *config.Config
}

func NewServer(auth *controllers.AuthController, users *controllers.UsersController, config *config.Config) *Server {
	return &Server{
		auth:   auth,
		users:  users,
		config: config,
	}
}

// Login performs a login request
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest := api.LoginJSONRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		slog.Error("failed to decode login request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := s.auth.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		slog.Error("failed to login", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	authResponse := api.AuthResponse{
		AccessToken:  token,
		ExpiresIn:    s.config.JWT.Expiration,
		RefreshToken: token,
		TokenType:    "Bearer",
		User: api.User{
			Id:    user.Id,
			Email: user.Email,
			Name:  user.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
	w.WriteHeader(http.StatusOK)
}

// Register pserfoms an account registration request
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	registerRequest := api.RegisterJSONRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		slog.Error("failed to decode auth request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := s.auth.Register(registerRequest.Name, registerRequest.Email, registerRequest.Password)
	if err != nil {
		slog.Error("failed to register user", "error", err)

		if errors.Is(err, controllers.ErrInvalidName) || errors.Is(err, controllers.ErrInvalidEmail) || errors.Is(err, controllers.ErrInvalidPassword) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authResponse := api.AuthResponse{
		AccessToken:  token,
		ExpiresIn:    s.config.JWT.Expiration,
		RefreshToken: token,
		TokenType:    "Bearer",
		User: api.User{
			Id:    user.Id,
			Email: user.Email,
			Name:  user.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(authResponse)
}

// Gets a user by their ID
func (s *Server) GetUserByName(w http.ResponseWriter, r *http.Request, name string) {
	user, err := s.users.GetUserByName(name)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		protected := r.Context().Value(api.JwtAuthScopes) != nil

		if !protected {
			next.ServeHTTP(w, r)
			return
		}

		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(bearer, "Bearer ")
		if token == "" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		user, err := s.auth.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
