package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	slog.Info("Recieved authenticate request")
	url, err := h.auth.CreateAuthURL(r.Context())
	if err != nil {
		slog.Error("creating auth URL", "error", err)
		http.Error(w, err.Error(), 500)
		return
	}
	slog.Info("Built authentication URL", "URL", url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	slog.Info("Recieved callback request", "code", code, "state", state)

	sessionID, err := h.auth.ObtainToken(r.Context(), code, state)
	if err != nil {
		slog.Error("obtaining token", "error", err)
		http.Error(w, err.Error(), 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "SESSION_ID",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   true,
		Domain:   "localhost",
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
	})
	http.Redirect(w, r, fmt.Sprintf("%s:%d", h.env.Client.URL, h.env.Client.Port), http.StatusSeeOther)
}

func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("SESSION_ID")
	if err != nil {
		slog.Error("getting session ID", "error", err)
		http.Error(w, err.Error(), 500)
		return
	}

	_, found, err := h.auth.GetSession(r.Context(), sessionID.Value)
	if err != nil {
		slog.Error("getting session", "error", err)
		http.Error(w, err.Error(), 500)
		return
	}

	switch found {
	case true:
		w.WriteHeader(http.StatusOK)
	case false:
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("SESSION_ID")
	if err != nil {
		slog.Error("getting session ID", "error", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if sessionID == nil {
		slog.Error("session ID is empty")
		http.Error(w, "session ID is empty", http.StatusUnauthorized)
		return
	}

	found, err := h.auth.DeleteSession(r.Context(), sessionID.Name)
	if err != nil {
		slog.Error("deleting session", "error", err)
		http.Error(w, err.Error(), 500)
		return
	}

	if found {
		http.SetCookie(w, &http.Cookie{
			Name:     "SESSION_ID",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Domain:   "localhost",
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Unix(0, 0),
		})
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
