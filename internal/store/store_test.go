package store

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"github.com/gitvam/platform-go-challenge/internal/models"
	"github.com/lib/pq"
)

func getTestConnStr() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	return fmt.Sprintf("postgres://gwi:password@%s:5432/favorites?sslmode=disable", host)
}

func resetTestDB(db *sql.DB) {
	db.Exec("DELETE FROM favorites")
	db.Exec("DELETE FROM charts")
	db.Exec("DELETE FROM insights")
	db.Exec("DELETE FROM audiences")
}

func TestAddFavorite_Success(t *testing.T) {
	s, err := NewPostgresStore(getTestConnStr())
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	resetTestDB(s.db)

	_, err = s.db.Exec(`
		INSERT INTO charts (external_id, title, x_axis_title, y_axis_title, data, description)
		VALUES ('test_chart_extid', 'Chart Title', 'X', 'Y', ARRAY[1,2,3], 'desc')`)
	if err != nil {
		t.Fatalf("failed to insert test chart: %v", err)
	}

	chart := &models.Chart{
		ExternalID:  "test_chart_extid",
		Title:       "Chart Title",
		XAxisTitle:  "X",
		YAxisTitle:  "Y",
		Data:        pq.Int64Array{1, 2, 3},
		Description: "desc",
		Type:        "chart",
	}

	err = s.AddFavorite("11111111-1111-1111-1111-111111111111", chart)
	if err != nil {
		t.Fatalf("AddFavorite failed: %v", err)
	}

	favs, err := s.ListFavorites("11111111-1111-1111-1111-111111111111", 10, 0)
	if err != nil {
		t.Fatalf("ListFavorites failed: %v", err)
	}

	if len(favs) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(favs))
	}
	if favs[0].GetType() != "chart" {
		t.Errorf("expected type 'chart', got %q", favs[0].GetType())
	}
}

func TestAddFavorite_Duplicate(t *testing.T) {
	s, err := NewPostgresStore(getTestConnStr())
	if err != nil {
		t.Fatal(err)
	}
	resetTestDB(s.db)

	s.db.Exec(`INSERT INTO insights (external_id, text, description) VALUES ('dup_insight', 'some text', 'desc')`)
	insight := &models.Insight{
		ExternalID:  "dup_insight",
		Text:        "some text",
		Description: "desc",
		Type:        "insight",
	}

	err = s.AddFavorite("11111111-1111-1111-1111-111111111111", insight)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AddFavorite("11111111-1111-1111-1111-111111111111", insight)
	if err == nil {
		t.Fatal("expected duplicate error, got nil")
	}
}

func TestAddFavorite_Invalid(t *testing.T) {
	s, err := NewPostgresStore(getTestConnStr())
	if err != nil {
		t.Fatal(err)
	}
	invalid := &models.Chart{ExternalID: "", Title: ""}
	err = s.AddFavorite("11111111-1111-1111-1111-111111111111", invalid)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestListFavorites_EmptyUser(t *testing.T) {
	s, err := NewPostgresStore(getTestConnStr())
	if err != nil {
		t.Fatal(err)
	}
	resetTestDB(s.db)
	favs, err := s.ListFavorites("33333333-3333-3333-3333-333333333333", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(favs) != 0 {
		t.Errorf("expected 0 favorites, got %d", len(favs))
	}
}

func TestAddAndListDifferentAssets(t *testing.T) {
	s, err := NewPostgresStore(getTestConnStr())
	if err != nil {
		t.Fatal(err)
	}
	resetTestDB(s.db)

	s.db.Exec(`INSERT INTO charts (external_id, title, x_axis_title, y_axis_title, data, description) VALUES ('chart_c1', 't', 'x', 'y', ARRAY[1], 'd')`)
	s.db.Exec(`INSERT INTO insights (external_id, text, description) VALUES ('insight_i1', 't', 'd')`)
	s.db.Exec(`INSERT INTO audiences (external_id, gender, birth_country, age_groups, hours_on_social, purchases_last_month, description) VALUES ('audience_a1', 'f', 'GR', ARRAY['18-24'], 2, 1, 'd')`)

	assets := []models.Asset{
		&models.Chart{ExternalID: "chart_c1", Title: "t", XAxisTitle: "x", YAxisTitle: "y", Data: pq.Int64Array{1, 2, 3}, Description: "d", Type: "chart"},
		&models.Insight{ExternalID: "insight_i1", Text: "t", Description: "d", Type: "insight"},
		&models.Audience{ExternalID: "audience_a1", Gender: "f", BirthCountry: "GR", AgeGroups: []string{"18-24"}, HoursOnSocial: 2, PurchasesLastMonth: 1, Description: "d", Type: "audience"},
	}

	for _, asset := range assets {
		if err := s.AddFavorite("22222222-2222-2222-2222-222222222222", asset); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	favs, err := s.ListFavorites("22222222-2222-2222-2222-222222222222", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(favs) != 3 {
		t.Errorf("expected 3 favorites, got %d", len(favs))
	}
}
