package main

import (
	"strings"
)

// TODO: Check if responseWithError works properly

func RemoveProfane(chirp string) string {
	profaneWords := map[string]bool{
		"KERFUFFLE": true,
		"SHARBERT":  true,
		"FORNAX":    true,
	}

	chirpTokens := strings.Split(chirp, " ")

	for index, token := range chirpTokens {
		if _, ok := profaneWords[strings.ToUpper(token)]; ok {
			chirpTokens[index] = "****"
		}
	}

	return strings.Join(chirpTokens, " ")
}
