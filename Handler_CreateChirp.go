package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FazecatGit/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Body == "" || req.UserID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp body exceeds 140 characters")
		return
	}
	now := time.Now().UTC()
	chirpID := uuid.New()

	dbChirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        chirpID,
		CreatedAt: now,
		UpdatedAt: now,
		Body:      req.Body,
		UserID:    req.UserID,
	},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	chirp := struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJson(w, http.StatusCreated, chirp)
}
