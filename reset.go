package main

import (
	"net/http"
)

func (cfg *apiConfig) adminResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset database")
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"status": "database reset"})
}
