package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model

	Id             string        `json:"id,omitempty"`
	Type           string        `json:"type,omitempty"`
	Collection     string        `json:"collection,omitempty"`
	StacVersion    string        `json:"stac_version,omitempty"`
	StacExtensions []string      `json:"stac_extensions,omitempty"`
	Bbox           []float64     `json:"bbox,omitempty"`
	Geometry       interface{}   `json:"geometry,omitempty"`
	Properties     interface{}   `json:"properties,omitempty"`
	Assets         interface{}   `json:"assets,omitempty"`
	Links          []interface{} `json:"links,omitempty"`
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
