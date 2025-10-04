package main

import (
	"encoding/json"
	"net/http"

	"github.com/FazecatGit/Chirpy/internal/auth"
	"github.com/FazecatGit/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	user := struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          dbUser.ID.String(),
		Email:       dbUser.Email,
		CreatedAt:   dbUser.CreatedAt.Format("2006-01-02T15:04:05.000Z07:00"),
		UpdatedAt:   dbUser.UpdatedAt.Format("2006-01-02T15:04:05.000Z07:00"),
		IsChirpyRed: dbUser.IsChirpyRed,
	}
	respondWithJson(w, http.StatusCreated, user)
}
