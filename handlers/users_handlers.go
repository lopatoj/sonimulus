package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Gets a user by their ID
func (h *Handler) GetUserByName(w http.ResponseWriter, r *http.Request, name string) {
	user, err := h.users.GetUserByName(name)
	if err != nil {
		slog.Error("failed to get user by id", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
