package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

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

type StacCollection struct {
	StacVersion string        `json:"stac_version,omitempty"`
	Id          string        `json:"id,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Keywords    []string      `json:"keywords,omitempty"`
	License     string        `json:"license,omitempty"`
	Providers   []interface{} `json:"providers,omitempty"`
	Extent      interface{}   `json:"extent,omitempty"`
	Summaries   interface{}   `json:"summary,omitempty"`
	Links       []Link        `json:"links,omitempty"`
	ItemType    string        `json:"itemType,omitempty"`
	Crs         []string      `json:"crs,omitempty"`
}

type Collection struct {
	gorm.Model

	Id   string `json:"id,omitempty"`
	Data JSONB  `gorm:"type:jsonb" json:"data,omitempty"`
}

type Root struct {
	StacVersion string `json:"stac_version,omitempty"`
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Links       []Link `json:"links,omitempty"`
}
