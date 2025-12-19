package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type responseErr struct {
		Msg string `json:"error"`
	}

	responseBody := responseErr{
		Msg: msg,
	}
	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("error in responseWithError function, specifically marshalling: %v", err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"something unexpected went wrong"}`))
	}

	w.WriteHeader(code)
	w.Write(dat)
}
