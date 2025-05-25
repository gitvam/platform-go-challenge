package models

import "fmt"

// AssetType is used to distinguish types
type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// Asset is the interface for all assets
// swagger:model Asset
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
	ID          string `json:"id"`
	Title       string `json:"title"`
	XAxisTitle  string `json:"x_axis_title"`
	YAxisTitle  string `json:"y_axis_title"`
	Data        []int  `json:"data"`
	Description string `json:"description"`
}

func (c *Chart) GetID() string              { return c.ID }
func (c *Chart) GetType() AssetType         { return AssetTypeChart }
func (c *Chart) GetDescription() string     { return c.Description }
func (c *Chart) SetDescription(desc string) { c.Description = desc }
func (c *Chart) Validate() error {
	if c.ID == "" || c.Title == "" {
		return fmt.Errorf("chart must have id and title")
	}
	return nil
}

// Insight Asset
// swagger:model Insight
type Insight struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Description string `json:"description"`
}

func (i *Insight) GetID() string              { return i.ID }
func (i *Insight) GetType() AssetType         { return AssetTypeInsight }
func (i *Insight) GetDescription() string     { return i.Description }
func (i *Insight) SetDescription(desc string) { i.Description = desc }
func (i *Insight) Validate() error {
	if i.ID == "" || i.Text == "" {
		return fmt.Errorf("insight must have id and text")
	}
	return nil
}

// Audience Asset
// swagger:model Asset
type Audience struct {
	ID                 string   `json:"id"`
	Gender             string   `json:"gender"`
	BirthCountry       string   `json:"birth_country"`
	AgeGroups          []string `json:"age_groups"`
	HoursOnSocial      int      `json:"hours_on_social"`
	PurchasesLastMonth int      `json:"purchases_last_month"`
	Description        string   `json:"description"`
}

func (a *Audience) GetID() string              { return a.ID }
func (a *Audience) GetType() AssetType         { return AssetTypeAudience }
func (a *Audience) GetDescription() string     { return a.Description }
func (a *Audience) SetDescription(desc string) { a.Description = desc }
func (a *Audience) Validate() error {
	if a.ID == "" || a.Gender == "" || a.BirthCountry == "" {
		return fmt.Errorf("audience must have id, gender, and birth country")
	}
	return nil
}
