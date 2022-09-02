package models

import "gorm.io/gorm"

type Link struct {
	Rel   string `json:"rel,omitempty"`
	Href  string `json:"href,omitempty"`
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

type Collection struct {
	gorm.Model

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

type Root struct {
	gorm.Model

	StacVersion string `json:"stac_version,omitempty"`
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Links       []Link `json:"links,omitempty"`
}
