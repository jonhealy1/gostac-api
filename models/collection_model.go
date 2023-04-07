package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type Link struct {
	Rel   string `json:"rel,omitempty"`
	Href  string `json:"href,omitempty"`
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

type Spatial struct {
	Bbox [][]float64 `json:"bbox,omitempty"`
}

type Temporal struct {
	Interval [][]string `json:"interval,omitempty"`
}

type Extent struct {
	Spatial  Spatial  `json:"spatial,omitempty"`
	Temporal Temporal `json:"temporal,omitempty"`
	License  string   `json:"license,omitempty"`
}

type Providers struct {
	Name  string   `json:"name,omitempty"`
	Roles []string `json:"roles,omitempty"`
	Url   string   `json:"url,omitempty"`
}

type StacCollection struct {
	StacVersion    string                 `json:"stac_version,omitempty"`
	Id             string                 `json:"id,omitempty"`
	Title          string                 `json:"title,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Keywords       []string               `json:"keywords,omitempty"`
	StacExtensions []string               `json:"stac_extensions,omitempty"`
	License        string                 `json:"license,omitempty"`
	Providers      []Providers            `json:"providers,omitempty"`
	Extent         Extent                 `json:"extent,omitempty"`
	Summaries      map[string]interface{} `json:"summaries,omitempty"`
	Links          []Link                 `json:"links,omitempty"`
	ItemType       string                 `json:"itemType,omitempty"`
	Crs            []string               `json:"crs,omitempty"`
}

type Collection struct {
	gorm.Model

	Id        string     `json:"id,omitempty"`
	Data      JSONB      `gorm:"type:jsonb" json:"data,omitempty"`
	CreatedAt *time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt *time.Time `json:"UpdatedAt,omitempty"`
}

type Root struct {
	StacVersion string `json:"stac_version,omitempty"`
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Links       []Link `json:"links,omitempty"`
}
