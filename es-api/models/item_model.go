package models

import (
	"time"

	"github.com/lib/pq"
)

// polygon or multiline
type GeoJSONPoly struct {
	Type        string         `json:"type"`
	Coordinates [][][2]float64 `json:"coordinates"`
}

type StacItem struct {
	Id             string          `json:"id,omitempty"`
	Type           string          `json:"type,omitempty"`
	Collection     string          `json:"collection,omitempty"`
	StacVersion    string          `json:"stac_version,omitempty"`
	StacExtensions []string        `json:"stac_extensions,omitempty"`
	Bbox           pq.Float64Array `gorm:"type:float[]"`
	Geometry       GeoJSONPoly     `json:"geometry,omitempty"`
	Properties     interface{}     `json:"properties,omitempty"`
	Assets         interface{}     `json:"assets,omitempty"`
	Links          []interface{}   `json:"links,omitempty"`
}

type Item struct {
	Id         string     `json:"id,omitempty"`
	Collection string     `json:"collection,omitempty"`
	Data       string     `json:"data,omitempty"`
	Geometry   string     `json:"geometry,omitempty"`
	CreatedAt  *time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt  *time.Time `json:"UpdatedAt,omitempty"`
}

type Context struct {
	Returned int `json:"returned,omitempty"`
	Limit    int `json:"limit,omitempty"`
}

type ItemCollection struct {
	Type     string  `json:"type,omitempty"`
	Context  Context `json:"context,omitempty"`
	Features []Item  `json:"features,omitempty"`
}
