package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FazecatGit/Chirpy/internal/auth"
	"github.com/FazecatGit/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Body == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp body exceeds 140 characters")
		return
	}

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	userID, err := auth.ValidateJWT(tokenStr, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}
	now := time.Now().UTC()
	chirpID := uuid.New()

	dbChirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        chirpID,
		CreatedAt: now,
		UpdatedAt: now,
		Body:      req.Body,
		UserID:    userID,
	})

	type chirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	resp := chirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJson(w, http.StatusCreated, resp)
}

// func (cfg *apiConfig) validateChirpHandler(body string) (string, error) {
// 	badwords := map[string]struct{}{
// 		"kerfuffle": {},
// 		"fornax":    {},
// 		"sharbert":  {},
// 	}

// 	cleaned := cleanProfanity(body)
// 	return cleaned, nil
// }
