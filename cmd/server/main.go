package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/gitvam/platform-go-challenge/docs"
	"github.com/gitvam/platform-go-challenge/internal/handlers"
	"github.com/gitvam/platform-go-challenge/internal/middleware"
	"github.com/gitvam/platform-go-challenge/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	s, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	h := handlers.NewHandler(s)

	r := chi.NewRouter()

	// Global middleware
	r.Use(httprate.LimitByIP(10, 1*time.Minute))
	r.Use(middleware.Logging)

	// No auth for Swagger docs
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// JWT-protected API routes
	r.Group(func(api chi.Router) {
		api.Use(middleware.JWTAuthMiddleware)

		api.Route("/v1/users/{userID}/favorites", func(sr chi.Router) {
			sr.Get("/", h.ListFavorites)
			sr.Post("/", h.AddFavorite)
			sr.Delete("/{assetID}", h.RemoveFavorite)
			sr.Patch("/{assetID}", h.EditFavoriteDescription)
		})
	})

	log.Println("Server running on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
