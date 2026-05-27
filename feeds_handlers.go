package main

import (
	"encoding/json"
	"net/http"

	"github.com/Justified02/rssagg/internal/database"
)

type createFeedParams struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (cfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Decode the request body
	var params createFeedParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 2. Get user from context (middlewareAuth stored it here)
	user := r.Context().Value("user").(database.GetUserByAPIKeyRow)

	// 3. Insert feed into database
	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:   params.Name,
		Url:    params.URL,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Return feed as JSON
	respondWithJSON(w, http.StatusCreated, feed)
}

func (cfg *apiConfig) getFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, feeds)
}
