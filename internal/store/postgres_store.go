package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/gitvam/platform-go-challenge/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func (ps *PostgresStore) DB() *sql.DB {
	return ps.db
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (ps *PostgresStore) ListFavorites(userID string) ([]models.Asset, error) {
	var results []models.Asset

	// Charts
	chartQuery := `
		SELECT c.id, c.external_id, c.title, c.x_axis_title, c.y_axis_title, c.data, c.description
		FROM favorites f
		JOIN charts c ON f.asset_type = 'chart' AND f.asset_id = c.id
		WHERE f.user_id = $1
	`
	chartRows, err := ps.db.Query(chartQuery, userID)
	if err != nil {
		return nil, err
	}
	defer chartRows.Close()
	for chartRows.Next() {
		var c models.Chart
		if err := chartRows.Scan(&c.ID, &c.ExternalID, &c.Title, &c.XAxisTitle, &c.YAxisTitle, &c.Data, &c.Description); err != nil {
			return nil, err
		}
		c.Type = "chart"
		results = append(results, &c)
	}

	// Insights
	insightQuery := `
		SELECT i.id, i.external_id, i.text, i.description
		FROM favorites f
		JOIN insights i ON f.asset_type = 'insight' AND f.asset_id = i.id
		WHERE f.user_id = $1
	`
	insightRows, err := ps.db.Query(insightQuery, userID)
	if err != nil {
		return nil, err
	}
	defer insightRows.Close()
	for insightRows.Next() {
		var i models.Insight
		if err := insightRows.Scan(&i.ID, &i.ExternalID, &i.Text, &i.Description); err != nil {
			return nil, err
		}
		i.Type = "insight"
		results = append(results, &i)
	}

	// Audiences
	audienceQuery := `
		SELECT a.id, a.external_id, a.gender, a.birth_country, a.age_groups, a.hours_on_social, a.purchases_last_month, a.description
		FROM favorites f
		JOIN audiences a ON f.asset_type = 'audience' AND f.asset_id = a.id
		WHERE f.user_id = $1
	`
	audienceRows, err := ps.db.Query(audienceQuery, userID)
	if err != nil {
		return nil, err
	}
	defer audienceRows.Close()
	for audienceRows.Next() {
		var a models.Audience
		if err := audienceRows.Scan(&a.ID, &a.ExternalID, &a.Gender, &a.BirthCountry, &a.AgeGroups, &a.HoursOnSocial, &a.PurchasesLastMonth, &a.Description); err != nil {
			return nil, err
		}
		a.Type = "audience"
		results = append(results, &a)
	}

	return results, nil
}

func (ps *PostgresStore) AddFavorite(userID string, asset models.Asset) error {
	if err := asset.Validate(); err != nil {
		return err
	}

	assetType := string(asset.GetType())
	externalID := asset.GetID()
	var internalID int
	var query string

	switch assetType {
	case "chart":
		query = `SELECT id FROM charts WHERE external_id = $1`
	case "insight":
		query = `SELECT id FROM insights WHERE external_id = $1`
	case "audience":
		query = `SELECT id FROM audiences WHERE external_id = $1`
	default:
		return errors.New("unknown asset type")
	}

	err := ps.db.QueryRow(query, externalID).Scan(&internalID)
	if err != nil {
		return fmt.Errorf("could not resolve asset ID: %v", err)
	}

	insert := `
		INSERT INTO favorites (user_id, asset_id, asset_type, description)
		VALUES ($1, $2, $3, $4)
	`
	_, err = ps.db.Exec(insert, userID, internalID, assetType, asset.GetDescription())

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("asset already in favorites")
		}
		return err
	}
	return nil
}

func (ps *PostgresStore) RemoveFavorite(userID, assetType, externalID string) error {
	var assetID int
	var query string

	switch assetType {
	case "chart":
		query = `SELECT id FROM charts WHERE external_id = $1`
	case "insight":
		query = `SELECT id FROM insights WHERE external_id = $1`
	case "audience":
		query = `SELECT id FROM audiences WHERE external_id = $1`
	default:
		return errors.New("unknown asset type")
	}

	if err := ps.db.QueryRow(query, externalID).Scan(&assetID); err != nil {
		return fmt.Errorf("could not resolve asset ID: %v", err)
	}

	res, err := ps.db.Exec(`DELETE FROM favorites WHERE user_id = $1 AND asset_type = $2 AND asset_id = $3`, userID, assetType, assetID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("asset not found")
	}
	return nil
}

func (ps *PostgresStore) EditFavoriteDescription(userID, assetType, externalID, desc string) error {
	var assetID int
	var query string

	switch assetType {
	case "chart":
		query = `SELECT id FROM charts WHERE external_id = $1`
	case "insight":
		query = `SELECT id FROM insights WHERE external_id = $1`
	case "audience":
		query = `SELECT id FROM audiences WHERE external_id = $1`
	default:
		return errors.New("unknown asset type")
	}

	if err := ps.db.QueryRow(query, externalID).Scan(&assetID); err != nil {
		return fmt.Errorf("could not resolve asset ID: %v", err)
	}

	update := `
		UPDATE favorites
		SET description = $1
		WHERE user_id = $2 AND asset_type = $3 AND asset_id = $4`
	res, err := ps.db.Exec(update, desc, userID, assetType, assetID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("asset not found")
	}
	return nil
}
