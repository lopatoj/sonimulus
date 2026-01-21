package handlers

import (
	"context"
	"net/http"

	"lopa.to/sonimulus/env"
	"lopa.to/sonimulus/internal/repo"
)

type AuthController interface {
	CreateAuthURL(ctx context.Context) (url string, err error)
	ObtainToken(ctx context.Context, code string, state string) (token string, err error)
	GetSession(ctx context.Context, sessionID string) (session repo.Session, found bool, err error)
	DeleteSession(ctx context.Context, sessionID string) (found bool, err error)
}

type Handler struct {
	auth AuthController
	env  env.Env
}

func NewHandler(auth AuthController, e env.Env) *Handler {
	return &Handler{
		auth: auth,
		env:  e,
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
