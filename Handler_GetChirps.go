package main

import (
	"database/sql"
	"net/http"
	"sort"
	"time"

	"github.com/FazecatGit/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) listallChirpsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort") // "asc" or "desc"

	var chirps []database.Chirp
	var err error

	if authorID != "" {
		userUUID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}
		chirps, err = cfg.DB.ListChirpsByAuthor(ctx, userUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps by author")
			return
		}
	} else {
		chirps, err = cfg.DB.ListChirps(ctx)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
			return
		}
	}

	// In-memory sorting
	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		// default ascending
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	// Prepare response
	var response []struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	for _, dbChirp := range chirps {
		response = append(response, struct {
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
		})
	}

	respondWithJson(w, http.StatusOK, response)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirp_id")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirp")
		}
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

	respondWithJson(w, http.StatusOK, chirp)
}
