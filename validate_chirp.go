package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// TODO: Check if responseWithError works properly

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	type parameters struct {
		Chirp string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding paramters: %s", err)
		respondWithError(w, 500, "error decoding request body")
		return
	}

	if 140 < len(params.Chirp) {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"Chirp is too long"}`))
	} else {
		type returnVals struct {
			CleanedBody string `json:"cleaned_body"`
		}
		responseBody := returnVals{
			CleanedBody: RemoveProfane(params.Chirp),
		}
		dat, err := json.Marshal(responseBody)
		if err != nil {
			log.Printf("Error marshalling response: %s", err)
			respondWithError(w, 500, "error marshalling response")
			return
		}
		w.WriteHeader(200)
		w.Write(dat)
	}
}

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
