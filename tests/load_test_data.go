package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	print(len(fc.Features))

	var i int
	for i < (len(fc.Features) - 50) {
		test, _ := json.Marshal(fc.Features[i])
		responseBody := bytes.NewBuffer(test)
		resp, err := http.Post("http://localhost:6002/collections/sentinel-s2-l2a-cogs-test/items", "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		i = i + 1
	}
}
