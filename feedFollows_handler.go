package main

import (
	"encoding/json"
	"net/http"

	"github.com/Justified02/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type createFeedFollowParams struct {
	FeedID uuid.UUID `json:"feed_id"`
}

func (cfg *apiConfig) CreateFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get authenticated user
	user := r.Context().Value("user").(database.GetUserByAPIKeyRow)

	// 2. Decode request body
	var params createFeedFollowParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Call database
	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: params.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Return JSON
	respondWithJSON(w, http.StatusCreated, feedFollow)
}

func (cfg *apiConfig) getFeedFollowsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the authenticated user
	user := r.Context().Value("user").(database.GetUserByAPIKeyRow)

	// 2. Query DB
	feedFollows, err := cfg.DB.GetFeedFollowsForUser(
		r.Context(),
		user.ID,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. Return JSON
	respondWithJSON(w, http.StatusOK, feedFollows)
}

func (cfg *apiConfig) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the feedFollowID from the URL
	// Chi router allows extracting URL parameters
	feedFollowID := chi.URLParam(r, "feedFollowID")

	// 2. convert the id to uuid
	id, err := uuid.Parse(feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. Get the authenticated user - we need user.ID to make sure the deletion belongs to the user
	user := r.Context().Value("user").(database.GetUserByAPIKeyRow)

	// 4. Call the DB to delete
	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID: 	id,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 5. respond with success
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "unfollowed successfully",
	})
}