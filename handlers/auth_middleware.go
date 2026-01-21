package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := r.Cookie("SESSION_ID")
		if err != nil {
			slog.Error("Failed to get session id cookie", "error", err)
			h.redirectToLogin(w, r)
			return
		}

		session, found, err := h.auth.GetSession(r.Context(), sessionID.Value)
		if err != nil {
			slog.Error("Failed to get session", "error", err)
			h.redirectToLogin(w, r)
			return
		}

		if !found {
			slog.Error("Session not found")
			h.redirectToLogin(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "access_token", session.Token.AccessToken)
		ctx = context.WithValue(ctx, "user", session.User)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s:%d/login", h.env.Client.URL, h.env.Client.Port), http.StatusSeeOther)
}
