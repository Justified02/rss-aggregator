package main

import (
	"net/http"
	"time"

	"github.com/Justified02/rssagg/internal/database"
)

// Local API response struct (DTO)
type UserResponse struct {
	ID 		string 		`json:"id"`
	Name 	string 		`json:"name"`
	ApiKey 	string 		`json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Handler for GET /v1/users
func (cfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Extract user from context
	dbUser := r.Context().Value("user")
	if dbUser == nil {
		respondWithError(w, http.StatusUnauthorized, "user not found in context")
		return
	}

	user := dbUser.(database.GetUserByAPIKeyRow) 	// type assert to sqlc User

	// 2. Convert to API response
	resp := UserResponse{
		ID: 	user.ID.String(),
		Name: 	user.Name,
		ApiKey:	user.ApiKey,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// 3. Send JSON response
	respondWithJSON(w, http.StatusOK, resp)
}