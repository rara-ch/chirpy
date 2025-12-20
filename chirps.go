package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rara-ch/chirpy/internal/auth"
	"github.com/rara-ch/chirpy/internal/database"
)

type chirpResponseBody struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Chirp     string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
	}
	userID, err := auth.ValidateJWT(token, cfg.signiture)
	if err != nil {
		w.WriteHeader(401)
	}

	type parameters struct {
		Chirp string `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
			UserID: userID,
		})

		if err != nil {
			log.Printf("error creating chirp: %s", err)
			respondWithError(w, 500, "error creating chirp")
			return
		}

		responseData := chirpResponseBody{
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

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Printf("error reading chirp id: %s", err)
		respondWithError(w, 500, "error reading chirp id")
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
		}

		log.Printf("error getting chirp: %s", err)
		respondWithError(w, 500, "error getting chirp")
		return
	}

	responseData := chirpResponseBody{
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
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("error getting chirps: %s", err)
		respondWithError(w, 500, "error getting chirps")
		return
	}

	responseBody := []chirpResponseBody{}
	for _, chirp := range chirps {
		chirpResponseData := chirpResponseBody{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Chirp:     chirp.Body,
			UserId:    chirp.UserID,
		}

		responseBody = append(responseBody, chirpResponseData)
	}

	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("error encoding/marshalling chirps: %s", err)
		respondWithError(w, 500, "error encoding/marshalling chirps")
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}
