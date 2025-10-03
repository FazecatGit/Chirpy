package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/FazecatGit/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int64  `json:"expires_in_seconds"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	ok, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if req.ExpiresInSeconds == 0 || req.ExpiresInSeconds > 3600 {
		req.ExpiresInSeconds = 3600
	}
	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, (time.Duration(req.ExpiresInSeconds) * time.Second))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"token":      token,
	})

}
