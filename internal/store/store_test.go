package store

import (
	"testing"
	"github.com/gitvam/platform-go-challenge/internal/models"
)

func makeTestChart(id string) *models.Chart {
	return &models.Chart{
		ID:          id,
		Title:       "Test Chart",
		XAxisTitle:  "X",
		YAxisTitle:  "Y",
		Data:        []int{1, 2, 3},
		Description: "desc",
	}
}
func makeTestInsight(id string) *models.Insight {
	return &models.Insight{
		ID:          id,
		Text:        "Test insight",
		Description: "desc",
	}
}
func makeTestAudience(id string) *models.Audience {
	return &models.Audience{
		ID:                 id,
		Gender:             "female",
		BirthCountry:       "Italy",
		AgeGroups:          []string{"18-24"},
		HoursOnSocial:      3,
		PurchasesLastMonth: 2,
		Description:        "desc",
	}
}

func TestAddFavorite_Success(t *testing.T) {
	s := NewInMemoryStore()
	asset := makeTestChart("chart1")

	if err := s.AddFavorite("user1", asset); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	favs, err := s.ListFavorites("user1")
	if err != nil {
		t.Fatalf("list favorites: %v", err)
	}
	if len(favs) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(favs))
	}
	if favs[0].GetID() != "chart1" {
		t.Errorf("expected ID 'chart1', got %q", favs[0].GetID())
	}
}

func TestAddFavorite_Duplicate(t *testing.T) {
	s := NewInMemoryStore()
	asset := makeTestChart("dup1")
	if err := s.AddFavorite("user1", asset); err != nil {
		t.Fatal(err)
	}
	err := s.AddFavorite("user1", asset)
	if err == nil || err.Error() != "asset already in favorites" {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestAddFavorite_Invalid(t *testing.T) {
	s := NewInMemoryStore()
	invalid := &models.Chart{ID: "", Title: ""} // Invalid chart
	err := s.AddFavorite("user1", invalid)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestListFavorites_EmptyUser(t *testing.T) {
	s := NewInMemoryStore()
	favs, err := s.ListFavorites("noone")
	if err != nil {
		t.Fatal(err)
	}
	if len(favs) != 0 {
		t.Errorf("expected 0 favorites, got %d", len(favs))
	}
}

func TestRemoveFavorite(t *testing.T) {
	s := NewInMemoryStore()
	asset := makeTestChart("chart1")
	_ = s.AddFavorite("u", asset)
	// Remove existing
	err := s.RemoveFavorite("u", "chart1")
	if err != nil {
		t.Fatal(err)
	}
	// Remove again, should fail
	err = s.RemoveFavorite("u", "chart1")
	if err == nil {
		t.Fatal("expected error on removing non-existing asset")
	}
	// Remove from non-existent user
	err = s.RemoveFavorite("nope", "whatever")
	if err == nil {
		t.Fatal("expected error on removing from non-existent user")
	}
}

func TestEditFavoriteDescription(t *testing.T) {
	s := NewInMemoryStore()
	asset := makeTestInsight("insight1")
	_ = s.AddFavorite("testu", asset)
	// Success
	err := s.EditFavoriteDescription("testu", "insight1", "new desc")
	if err != nil {
		t.Fatal(err)
	}
	favs, _ := s.ListFavorites("testu")
	if favs[0].GetDescription() != "new desc" {
		t.Errorf("expected updated description, got %q", favs[0].GetDescription())
	}
	// Edit non-existent asset
	err = s.EditFavoriteDescription("testu", "noid", "whatever")
	if err == nil {
		t.Fatal("expected error for missing asset")
	}
	// Edit for non-existent user
	err = s.EditFavoriteDescription("nouser", "noid", "whatever")
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}

func TestAddAndListDifferentAssets(t *testing.T) {
	s := NewInMemoryStore()
	assets := []models.Asset{
		makeTestChart("c1"),
		makeTestInsight("i1"),
		makeTestAudience("a1"),
	}
	for _, a := range assets {
		if err := s.AddFavorite("userx", a); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	favs, err := s.ListFavorites("userx")
	if err != nil {
		t.Fatal(err)
	}
	if len(favs) != 3 {
		t.Errorf("expected 3 favorites, got %d", len(favs))
	}
}
