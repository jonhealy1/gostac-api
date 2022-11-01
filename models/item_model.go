package models

import (
	"github.com/lib/pq"
)

// type EWKBGeomPolygon geom.Polygon

// func (g *EWKBGeomPolygon) Scan(input interface{}) error {
// 	gt, err := ewkb.Unmarshal(input.([]byte))
// 	if err != nil {
// 		return err
// 	}
// 	g = gt.(*EWKBGeomPolygon)

// 	return nil
// }

// func (g EWKBGeomPolygon) Value() (driver.Value, error) {
// 	b := geom.Polygon(g)
// 	bp := &b
// 	ewkbPt := ewkb.Polygon{Polygon: bp.SetSRID(4326)}
// 	return ewkbPt.Value()
// }

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
	Id         string `json:"id,omitempty"`
	Collection string `json:"collection,omitempty"`
	Data       string `json:"data,omitempty"`
	Geometry   string `json:"geometry,omitempty"`
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
