package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Justified02/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// 1. Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	//port := os.Getenv("PORT")

	// 2. Open DB connection
	conn, err := sql.Open("postgres", dbURL) 	// lib/pq driver
	if err != nil {
		log.Fatal("Cant connect to DB:", err)
	}

	// 3. Test Connection
	if err := conn.Ping(); err != nil {
		log.Fatal("Cannot ping DB:", err)
	}

	log.Println("Connected to Postgres!")

	// 4. Wrap connection with sqlc Queries
	cfg := &apiConfig{
		DB: database.New(conn),
	}

	// 5. Router setup will come here
	r := chi.NewRouter() 	// create a new router

	r.Get("/healthz", healthCheckHandler) 	// register route + handler
	r.Get("/ping", pingHandler)

	// User routes
	r.Post("/users", cfg.handlerCreateUser)
	r.Get("/users", cfg.middlewareAuth(cfg.getUserHandler))

	// Feed routes
	r.Post("/feeds", cfg.middlewareAuth(cfg.createFeedHandler))
	r.Get("/feeds", cfg.getFeedsHandler)

	// Feed follows routes
	r.Post("/feed_follows", cfg.middlewareAuth(cfg.CreateFeedFollowHandler))
	r.Get("/feed_follows", cfg.middlewareAuth(cfg.getFeedFollowsHandler))
	r.Delete("/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.deleteFeedFollowHandler))

	r.Mount("/v1", r)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r) 	// start http server on port 8080
}

