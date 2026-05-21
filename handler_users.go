package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Justified02/rssagg/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Name string `json:"name"`
}

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func generateAPIKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID: 		dbUser.ID.String(),
		Name: 		dbUser.Name,
		ApiKey: 	dbUser.ApiKey,
		CreatedAt:	dbUser.CreatedAt,
		UpdatedAt:	dbUser.UpdatedAt,
	}
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. Decode the request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	// 2. Generate UUID and API Key
	id := uuid.New() 	// unique user id
	apiKey, err := generateAPIKey()
	if err != nil {
		http.Error(w, "Cannot generate API key", 500)
		return
	}

	// 3. Call sqlc
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: 	id,
		Name: 	params.Name,
		ApiKey: apiKey,
	})
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	// 4. Map DB object to API response
	apiUser := databaseUserToUser(user)
	respondWithJSON(w, 201, apiUser)

}