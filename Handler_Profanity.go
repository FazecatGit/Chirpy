package main

import (
	"strings"
)

func cleanProfanity(text string) string {

	words := strings.Split(text, " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "fornax" || strings.ToLower(word) == "sharbert" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
