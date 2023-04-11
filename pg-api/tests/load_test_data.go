package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jonhealy1/goapi-stac/pg-api/database"
	"github.com/jonhealy1/goapi-stac/pg-api/models"

	"github.com/go-playground/validator"
)

func LoadCollection() {
	jsonFile, err := os.Open("setup_data/test-collection.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	responseBody := bytes.NewBuffer(byteValue)

	collection := models.Collection{
		Data: models.JSONB{(&responseBody)},
		Id:   "sentinel-s2-l2a-cogs-test",
	}
	validator := validator.New()
	_ = validator.Struct(collection)
	err = database.DB.Db.Create(&collection).Error
}

func LoadItems() {
	jsonFile, err := os.Open("setup_data/sentinel-s2-l2a-cogs_0_100.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	type FeatureCollection struct {
		Type     string        `json:"type"`
		Features []interface{} `json:"features"`
	}
	var fc FeatureCollection
	json.Unmarshal(byteValue, &fc)

	var i int
	for i < (len(fc.Features) - 50) {
		test, _ := json.Marshal(fc.Features[i])
		stac_item := new(models.StacItem)
		json.Unmarshal(test, &stac_item)

		coordinatesString := "[["
		for _, s := range stac_item.Geometry.Coordinates[0] {
			coordinatesString = coordinatesString + fmt.Sprintf("[%f, %f],", s[0], s[1])
		}
		coordinatesString = coordinatesString + "]]"
		rawGeometryJSON := fmt.Sprintf("{'type':'Polygon', 'coordinates':%s}", coordinatesString)
		err = database.DB.Db.Exec(
			`INSERT INTO items (id, collection, data, geometry) 
			VALUES (
				@id, 
				@collection, 
				@data, 
				ST_GeomFromEWKB(ST_GeomFromGeoJSON(@geometry)))`,
			sql.Named("id", stac_item.Id),
			sql.Named("collection", stac_item.Collection),
			sql.Named("data", stac_item),
			sql.Named("geometry", rawGeometryJSON),
		).Error

		i = i + 1
	}
}
