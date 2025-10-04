package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/FazecatGit/Chirpy/internal/auth"
	"github.com/FazecatGit/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)

	rtParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
	_, err = cfg.DB.CreateRefreshToken(r.Context(), rtParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to persist refresh token")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"id":            user.ID,
		"email":         user.Email,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
		"token":         token,
		"refresh_token": refreshToken,
		"is_chirpy_red": user.IsChirpyRed,
	})

}
