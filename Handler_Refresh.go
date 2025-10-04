package main

import (
	"net/http"
	"time"

	"github.com/FazecatGit/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil || tokenStr == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), tokenStr)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt make JWT")
		return
	}

	respondWithJson(w, http.StatusOK, map[string]interface{}{
		"token": newToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil || tokenStr == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}

	_, err = cfg.DB.RevokeRefreshToken(r.Context(), tokenStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// refreshedTokenRecord, err := cfg.DB.GetUserFromRefreshToken(r.Context(), tokenStr)
// if err != nil {
// 	respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
// 	return
// }
