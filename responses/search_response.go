package responses

import (
	"github.com/jonhealy1/goapi-stac/models"

	"github.com/lib/pq"
)

type SearchResponse struct {
	Status   int        `json:"status"`
	Message  string     `json:"message"`
	Type     string     `json:"type"`
	Context  Context    `json:"context"`
	Features []StacItem `json:"features"`
}

type Context struct {
	Returned int `json:"returned"`
	Limit    int `json:"limit"`
}

type StacItem struct {
	Id             string             `json:"id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Collection     string             `json:"collection,omitempty"`
	StacVersion    string             `json:"stac_version,omitempty"`
	StacExtensions []string           `json:"stac_extensions,omitempty"`
	Bbox           pq.Float64Array    `gorm:"type:float[]"`
	Geometry       models.GeoJSONPoly `json:"geometry,omitempty"`
	Properties     interface{}        `json:"properties,omitempty"`
	Assets         interface{}        `json:"assets,omitempty"`
	Links          []interface{}      `json:"links,omitempty"`
}
