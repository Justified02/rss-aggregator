package main

import (
	"net/http"
	"strconv"

	"github.com/Justified02/rssagg/internal/database"
)


func (cfg *apiConfig) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get user from context
	user := r.Context().Value("user").(database.GetUserByAPIKeyRow)

	// 2. Read limit (max_posts) query parameter
	limitStr := r.URL.Query().Get("max_post")
	limit := int32(20)
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil && l > 0 {
			limit = int32(l)
		}
	}

	// 3. Call the DB query
	posts, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		MaxPosts: limit,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Return with JSON
	respondWithJSON(w, http.StatusOK, posts)
}