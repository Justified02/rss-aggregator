package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Justified02/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pg"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// 
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	r := chi.NewRouter() 	// create a new router

	r.Get("/healthz", healthCheckHandler) 	// register route + handler
	r.Get("/ping", pingHandler)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r) 	// start http server on port 8080
}

