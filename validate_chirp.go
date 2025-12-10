package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"something went wrong"}`))
		return
	}

	if 140 < len(params.Chirp) {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"Chirp is too long"}`))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(`{"valid":true}`))
	}
}
