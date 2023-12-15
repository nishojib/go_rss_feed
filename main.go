package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"rss_server/internal/database"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	apiConfig := &apiConfig{
		DB: dbQueries,
	}

	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiV1Router := chi.NewRouter()
	apiV1Router.Get("/readiness", apiConfig.readinessHandler)
	apiV1Router.Get("/err", apiConfig.errHandler)
	apiV1Router.Post("/users", apiConfig.createUserHandler)
	apiV1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.getUsersHandler))
	apiV1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.createFeedHandler))
	apiV1Router.Get("/feeds", apiConfig.getFeedsHandler)
	apiV1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.createFeedFollowHandler))
	apiV1Router.Delete("/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.deleteFeedFollowHandler))
	apiV1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.getFeedFollowsHandler))
	apiV1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.getPostsByUserHandler))
	mux.Mount("/api/v1", apiV1Router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("server listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}
