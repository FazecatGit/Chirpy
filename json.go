package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJson(w, status, map[string]string{"error": message})
}
