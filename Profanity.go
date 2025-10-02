package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func respondWithJson(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func cleanProfanity(text string) string {

	words := strings.Split(text, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "fornax" || strings.ToLower(word) == "sharbert" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
