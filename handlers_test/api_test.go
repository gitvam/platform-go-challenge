package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gitvam/platform-go-challenge/internal/handlers"
	"github.com/gitvam/platform-go-challenge/internal/middleware"
	"github.com/gitvam/platform-go-challenge/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "my_super_secret"

func generateJWT(sub string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(2 * 365 * 24 * time.Hour).Unix(),
	})
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic("failed to sign test token: " + err.Error())
	}
	return signed
}

func makeRequest(method, path, token, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	s := store.NewInMemoryStore()
	store.SeedDummyData(s)
	h := handlers.NewHandler(s)

	r := chi.NewRouter()
	r.Use(middleware.JWTAuthMiddleware)
	r.Route("/v1/users/{userID}/favorites", func(sr chi.Router) {
		sr.Get("/", h.ListFavorites)
		sr.Post("/", h.AddFavorite)
	})

	r.ServeHTTP(rr, req)
	return rr
}

func TestListFavorites_Johnsmith(t *testing.T) {
	token := generateJWT("johnsmith")
	rr := makeRequest("GET", "/v1/users/ignored/favorites", token, "")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "chart_engagement_2024") {
		t.Errorf("expected seeded chart, got: %s", rr.Body.String())
	}
}

func TestListFavorites_Maria(t *testing.T) {
	token := generateJWT("mariapapadopoulou")
	rr := makeRequest("GET", "/v1/users/ignored/favorites", token, "")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "chart_ecom_conversion") {
		t.Errorf("expected maria's chart, got: %s", rr.Body.String())
	}
}
