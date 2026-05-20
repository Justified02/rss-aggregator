package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json") 	// tell the client it is JSON
	w.WriteHeader(status) 	// set HTTP status code
	json.NewEncoder(w).Encode(payload) 	// convert Go struct/map to JSON
}