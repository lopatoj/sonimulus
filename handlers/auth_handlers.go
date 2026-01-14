package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"lopa.to/sonimulus/api"
	"lopa.to/sonimulus/controllers"
)

// Login performs a login request
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest := api.LoginJSONRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		slog.Error("failed to decode login request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.auth.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		slog.Error("failed to login", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	authResponse := api.AuthResponse{
		AccessToken:  token,
		ExpiresIn:    h.env.JWT.Expiration,
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
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	registerRequest := api.RegisterJSONRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		slog.Error("failed to decode auth request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.auth.Register(registerRequest.Name, registerRequest.Email, registerRequest.Password)
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
		ExpiresIn:    h.env.JWT.Expiration,
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
