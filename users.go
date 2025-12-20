package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rara-ch/chirpy/internal/auth"
	"github.com/rara-ch/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %v", err)
		respondWithError(w, 500, "error decoding parameters")
		return
	}

	// TODO: Implement strong password checks

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		respondWithError(w, 500, "error hashing password")
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Printf("error creating user: %v", err)
		respondWithError(w, 500, "error creating user")
		return
	}

	responseBody := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	dat, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("error marshalling data: %v", err)
		respondWithError(w, 500, "error marshalling data")
		return
	}

	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	// TODO: DRY the below code

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %v", err)
		respondWithError(w, 500, "error decoding parameters")
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("error getting user: %v", err)
		w.WriteHeader(401)
		w.Write([]byte(`{"message": "Incorrect email or password"}`))
	}

	isMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !isMatch {

		w.WriteHeader(401)
		w.Write([]byte(`{"message": "Incorrect email or password"}`))
	}

	if isMatch {
		responseBody := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		dat, err := json.Marshal(responseBody)
		if err != nil {
			log.Printf("error marshalling data: %v", err)
			respondWithError(w, 500, "error marshalling data")
			return
		}

		w.WriteHeader(200)
		w.Write(dat)
	}

}
