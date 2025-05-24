package main

import (
    "log"
    "net/http"
	"os"
    "github.com/go-chi/chi/v5"
    "github.com/gitvam/platform-go-challenge/internal/store"
    "github.com/gitvam/platform-go-challenge/internal/handlers"
	"github.com/gitvam/platform-go-challenge/internal/middleware"
)

func main() {
    // Initialize the in-memory store
    s := store.NewInMemoryStore()

	// Seed dummy data if running in dev environment
    if os.Getenv("APP_ENV") == "dev" {
        store.SeedDummyData(s)
    }

    h := handlers.NewHandler(s)

    // Set up the router
    r := chi.NewRouter()

	// Wire middleware
	r.Use(middleware.Logging)

    // Routes
    r.Route("/v1/users/{userID}/favorites", func(r chi.Router) {
        r.Get("/", h.ListFavorites)                 // GET    /v1/users/{userID}/favorites
        r.Post("/", h.AddFavorite)                  // POST   /v1/users/{userID}/favorites
        r.Delete("/{assetID}", h.RemoveFavorite)    // DELETE /v1/users/{userID}/favorites/{assetID}
        r.Patch("/{assetID}", h.EditFavoriteDescription) // PATCH /v1/users/{userID}/favorites/{assetID}
    })

    // Start the server
    log.Println("Server starting on :8080...")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("could not start server: %v", err)
    }
}
