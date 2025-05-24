package store

import (
	"errors"
	"sync"

	"github.com/gitvam/platform-go-challenge/internal/models"
)

// Store interface, so we can swap out storage (e.g., DB, cache) later.
type Store interface {
	ListFavorites(userID string) ([]models.Asset, error)
	AddFavorite(userID string, asset models.Asset) error
	RemoveFavorite(userID, assetID string) error
	EditFavoriteDescription(userID, assetID, desc string) error
}

// InMemoryStore implements Store using a map and RWMutex for concurrency.
type InMemoryStore struct {
	mu        sync.RWMutex
	favorites map[string]map[string]models.Asset // userID -> assetID -> Asset
}

// NewInMemoryStore creates an empty store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		favorites: make(map[string]map[string]models.Asset),
	}
}

func (s *InMemoryStore) ListFavorites(userID string) ([]models.Asset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userFavs, ok := s.favorites[userID]
	if !ok {
		return []models.Asset{}, nil
	}
	assets := make([]models.Asset, 0, len(userFavs))
	for _, asset := range userFavs {
		assets = append(assets, asset)
	}
	return assets, nil
}

func (s *InMemoryStore) AddFavorite(userID string, asset models.Asset) error {
    if err := asset.Validate(); err != nil {
        return err
    }
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.favorites[userID] == nil {
        s.favorites[userID] = make(map[string]models.Asset)
    }
    if _, exists := s.favorites[userID][asset.GetID()]; exists {
        return errors.New("asset already in favorites")
    }
    s.favorites[userID][asset.GetID()] = asset
    return nil
}


func (s *InMemoryStore) RemoveFavorite(userID, assetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	userFavs, ok := s.favorites[userID]
	if !ok {
		return errors.New("user has no favorites")
	}
	if _, found := userFavs[assetID]; !found {
		return errors.New("asset not found")
	}
	delete(userFavs, assetID)
	return nil
}

func (s *InMemoryStore) EditFavoriteDescription(userID, assetID, desc string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	userFavs, ok := s.favorites[userID]
	if !ok {
		return errors.New("user has no favorites")
	}
	asset, found := userFavs[assetID]
	if !found {
		return errors.New("asset not found")
	}
	asset.SetDescription(desc)
	userFavs[assetID] = asset
	return nil
}
