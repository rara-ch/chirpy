package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rara-ch/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	type parameters struct {
		Chirp  string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %v", err)
		respondWithError(w, 500, "error decoding parameters")
		return
	}

	if 140 < len(params.Chirp) {
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"Chirp is too long"}`))
	} else {
		cleanedChirp := RemoveProfane(params.Chirp)

		chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleanedChirp,
			UserID: params.UserID,
		})

		if err != nil {
			log.Printf("error creating chirp: %s", err)
			respondWithError(w, 500, "error creating chirp")
			return
		}

		type responseBody struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Chirp     string    `json:"body"`
			UserId    uuid.UUID `json:"user_id"`
		}

		responseData := responseBody{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Chirp:     chirp.Body,
			UserId:    chirp.UserID,
		}

		dat, err := json.Marshal(responseData)
		if err != nil {
			log.Printf("Error marshalling response: %s", err)
			respondWithError(w, 500, "error marshalling response")
			return
		}
		w.WriteHeader(201)
		w.Write(dat)
	}
}
