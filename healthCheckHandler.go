package main

import (
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 500, map[string]string{"status": "ok"})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"message": "pong"})
}