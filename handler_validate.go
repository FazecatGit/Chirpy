package main

import (
	"encoding/json"
	"net/http"
)

type ChirpRequest struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	var chirp Chirp
	err := json.NewDecoder(r.Body).Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanProfanity(chirp.Body)
	respondWithJson(w, http.StatusOK, map[string]string{"cleaned_body": cleanedBody})
}
