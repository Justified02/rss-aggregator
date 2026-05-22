package main

import (
	"context"
	"net/http"
	"strings"
)

// Middleware for API Key authentication - checks the authorization header, loads the user from the DB and stores it as context
func (cfg *apiConfig) middlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {

		// 1. Read Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		// 2. Expected format: "ApiKey <key>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "ApiKey" {
			respondWithError(w, http.StatusUnauthorized, "invalid Authorization header format")
			return
		}
		
		apiKey := parts[1]

		// 3. Look up user in the database
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid API Key")
			return
		}

		// 4. Store user in context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}