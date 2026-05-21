package main

import (
	"encoding/json"
	"net/http"
)

// func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
// 	w.Header().Set("Content-Type", "application/json") 	// tell the client it is JSON
// 	w.WriteHeader(status) 	// set HTTP status code
// 	json.NewEncoder(w).Encode(payload) 	// convert Go struct/map to JSON
// }

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("failed to marshal JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{
		"error": message,
	})
}