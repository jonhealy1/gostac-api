package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type StacItem struct {
	Id             string          `json:"id,omitempty"`
	Type           string          `json:"type,omitempty"`
	Collection     string          `json:"collection,omitempty"`
	StacVersion    string          `json:"stac_version,omitempty"`
	StacExtensions []string        `json:"stac_extensions,omitempty"`
	Bbox           pq.Float64Array `gorm:"type:float[]"`
	Geometry       interface{}     `json:"geometry,omitempty"`
	Properties     interface{}     `json:"properties,omitempty"`
	Assets         interface{}     `json:"assets,omitempty"`
	Links          []interface{}   `json:"links,omitempty"`
}

type Item struct {
	gorm.Model

	Id string `json:"id,omitempty"`
	// Type       string `json:"type,omitempty"`
	Collection string `json:"collection,omitempty"`
	// StacVersion string `json:"stac_version,omitempty"`
	// StacExtensions []string      `json:"stac_extensions,omitempty"`
	// Bbox pq.Float64Array `gorm:"type:float[]"`
	Data JSONB `gorm:"type:jsonb" json:"data,omitempty"`
	// Geometry       interface{}   `json:"geometry,omitempty"`
	// Properties     interface{}   `json:"properties,omitempty"`
	// Assets         interface{}   `json:"assets,omitempty"`
	// Links          []interface{} `json:"links,omitempty"`
}

type Context struct {
	Returned int `json:"returned,omitempty"`
	Limit    int `json:"limit,omitempty"`
}

type ItemCollection struct {
	gorm.Model

	Type     string  `json:"type,omitempty"`
	Context  Context `json:"context,omitempty"`
	Features []Item  `json:"features,omitempty"`
}
