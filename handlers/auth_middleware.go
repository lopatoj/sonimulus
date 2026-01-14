package handlers

import (
	"context"
	"net/http"
	"strings"

	"lopa.to/sonimulus/api"
)

// AuthMiddleware is a middleware that validates JWT tokens and sets the user in the context.
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
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

		user, err := h.auth.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
