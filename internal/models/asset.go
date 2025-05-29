package models

import (
	"fmt"

	"github.com/lib/pq"
)

// AssetType is used to distinguish types
type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// Asset is the interface for all assets
// swagger:model Asset
//
// swagger:oneOf Chart Insight Audience
// swagger:discriminator type
type Asset interface {
	GetID() string
	GetType() AssetType
	GetDescription() string
	SetDescription(desc string)
	Validate() error
}

// Chart Asset
// swagger:model Chart
type Chart struct {
	ID          int           `json:"id"`
	ExternalID  string        `json:"external_id"` // business-slug used by API
	Title       string        `json:"title"`
	XAxisTitle  string        `json:"x_axis_title"`
	YAxisTitle  string        `json:"y_axis_title"`
	Data        pq.Int64Array `json:"data" db:"data"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
}

func (c *Chart) GetID() string              { return c.ExternalID }
func (c *Chart) GetType() AssetType         { return AssetTypeChart }
func (c *Chart) GetDescription() string     { return c.Description }
func (c *Chart) SetDescription(desc string) { c.Description = desc }
func (c *Chart) Validate() error {
	if c.ExternalID == "" || c.Title == "" {
		return fmt.Errorf("chart must have external_id and title")
	}
	return nil
}

// Insight Asset
// swagger:model Insight
type Insight struct {
	ID          int    `json:"id"`
	ExternalID  string `json:"external_id"`
	Text        string `json:"text"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func (i *Insight) GetID() string              { return i.ExternalID }
func (i *Insight) GetType() AssetType         { return AssetTypeInsight }
func (i *Insight) GetDescription() string     { return i.Description }
func (i *Insight) SetDescription(desc string) { i.Description = desc }
func (i *Insight) Validate() error {
	if i.ExternalID == "" || i.Text == "" {
		return fmt.Errorf("insight must have id and text")
	}
	return nil
}

// Audience Asset
// swagger:model Audience
type Audience struct {
	ID                 int            `json:"id"`
	ExternalID         string         `json:"external_id"`
	Gender             string         `json:"gender"`
	BirthCountry       string         `json:"birth_country"`
	AgeGroups          pq.StringArray `db:"age_groups"`
	HoursOnSocial      int            `json:"hours_on_social"`
	PurchasesLastMonth int            `json:"purchases_last_month"`
	Description        string         `json:"description"`
	Type               string         `json:"type"`
}

func (a *Audience) GetID() string              { return a.ExternalID }
func (a *Audience) GetType() AssetType         { return AssetTypeAudience }
func (a *Audience) GetDescription() string     { return a.Description }
func (a *Audience) SetDescription(desc string) { a.Description = desc }
func (a *Audience) Validate() error {
	if a.ExternalID == "" || a.Gender == "" || a.BirthCountry == "" {
		return fmt.Errorf("audience must have id, gender, and birth country")
	}
	return nil
}
