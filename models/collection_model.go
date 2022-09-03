package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// JSONB Interface for JSONB Field of yourTableName Table
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

type Collection struct {
	gorm.Model

	// StacVersion string `gorm:"stac_version,omitempty"`
	Id   string `gorm:"id,omitempty"`
	Data JSONB  `gorm:"type:jsonb" json:"fieldnameofjsonb"`
	// Title       string `gorm:"title,omitempty"`
	// Description string `gorm:"description,omitempty"`
	// License     string `gorm:"license,omitempty"`
	// ItemType    string `gorm:"itemType,omitempty"`
}

type Root struct {
	StacVersion string `json:"stac_version,omitempty"`
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Links       []Link `json:"links,omitempty"`
}
