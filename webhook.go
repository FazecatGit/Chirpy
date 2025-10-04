package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Event {
	case "user.upgraded":
		userUUID, err := uuid.Parse(req.Data.UserID)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}
		user, err := cfg.DB.UpgradeUserToChirpyRed(r.Context(), userUUID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		_ = user
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
