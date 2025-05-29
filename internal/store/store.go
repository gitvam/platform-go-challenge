package store

import "github.com/gitvam/platform-go-challenge/internal/models"

type Store interface {
	ListFavorites(userID string) ([]models.Asset, error)
	AddFavorite(userID string, asset models.Asset) error
	RemoveFavorite(userID, assetType, externalID string) error
	EditFavoriteDescription(userID, assetType, externalID, desc string) error
}
