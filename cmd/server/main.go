// @title           GWI Favorites API
// @version         1.0
// @description     An API to manage user favorites (charts, insights, audiences) at GWI.
// @termsOfService  http://swagger.io/terms/

// @contact.name   George Vamvakousis
// @contact.email  geovam99@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /
package main

import (
	"fmt"
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

// johsmith token: JWT Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDgzNzMyNDEsInN1YiI6ImpvaG5zbWl0aCJ9.BPEdl1zvq3k4qqq3ewPRIdZVmvFsmugB0gYskvv8nEA

func main() {

	// Initialize the in-memory store
	s := store.NewInMemoryStore()

	// Seed dummy data if running in dev environment
	if os.Getenv("APP_ENV") == "dev" {
		store.SeedDummyData(s)
		fmt.Println("Loaded dummy data!")
	}

	h := handlers.NewHandler(s)

	// Set up the router
	r := chi.NewRouter()

	// Apply rate limiting: 10 requests per 10 seconds per IP
	r.Use(httprate.Limit(
		10,
		time.Minute,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"error": "Rate-limited. Please, slow down."}`, http.StatusTooManyRequests)
		}),
	))

	// Add logging middleware globally
	r.Use(middleware.Logging)

	// Swagger route (no auth)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// API routes (protected)
	r.Route("/v1/users/{userID}/favorites", func(sr chi.Router) {
		sr.Use(middleware.JWTAuthMiddleware)
		sr.Get("/", h.ListFavorites)
		sr.Post("/", h.AddFavorite)
		sr.Delete("/{assetID}", h.RemoveFavorite)
		sr.Patch("/{assetID}", h.EditFavoriteDescription)
	})

	// Start the server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
