package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gitvam/platform-go-challenge/internal/handlers"
	"github.com/gitvam/platform-go-challenge/internal/middleware"
	"github.com/gitvam/platform-go-challenge/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func getSignedToken(userID string) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "my_super_secret"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	})
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func setupTestRouter() http.Handler {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	dsn := fmt.Sprintf("postgres://gwi:password@%s:5432/favorites?sslmode=disable", host)

	s, err := store.NewPostgresStore(dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to db: %w", err))
	}

	// Ensure test chart exists
	_, _ = s.DB().Exec(`
		INSERT INTO charts (external_id, title, x_axis_title, y_axis_title, data, description)
		VALUES ('chart_engagement_2024', 'Engagement Q1', 'Month', 'Engagement', ARRAY[10,20,30], 'A seeded chart')
		ON CONFLICT (external_id) DO NOTHING
	`)

	h := handlers.NewHandler(s)

	r := chi.NewRouter()
	r.Use(middleware.JWTAuthMiddleware)
	r.Get("/v1/users/{userID}/favorites", h.ListFavorites)
	r.Post("/v1/users/{userID}/favorites", h.AddFavorite)
	r.Delete("/v1/users/{userID}/favorites/{assetID}", h.RemoveFavorite)
	r.Patch("/v1/users/{userID}/favorites/{assetID}", h.EditFavoriteDescription)

	return r
}


func TestAddAndListFavorite(t *testing.T) {
	router := setupTestRouter()
	userID := "11111111-1111-1111-1111-111111111111"
	token := getSignedToken(userID)

	addBody := `{
		"type": "chart",
		"external_id": "chart_engagement_2024",
		"title": "Engagement Chart",
		"x_axis_title": "Month",
		"y_axis_title": "Engagement",
		"data": [10, 20, 30],
		"description": "A test chart"
	}`

	req := httptest.NewRequest("POST", "/v1/users/"+userID+"/favorites", strings.NewReader(addBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.Code)
	}

	// List to verify
	reqList := httptest.NewRequest("GET", "/v1/users/"+userID+"/favorites", nil)
	reqList.Header.Set("Authorization", "Bearer "+token)
	respList := httptest.NewRecorder()
	router.ServeHTTP(respList, reqList)
	if respList.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", respList.Code)
	}
}