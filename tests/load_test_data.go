package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func LoadCollection() {
	jsonFile, err := os.Open("setup_data/test-collection.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	responseBody := bytes.NewBuffer(byteValue)

	resp, err := http.Post("http://localhost:6002/collections", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
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
