package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func LoadCollection() {
	jsonFile, err := os.Open("setup_data/collection.json")

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

// func LoadItems() {
// 	postBody, _ := json.Marshal(map[string]string{
// 		"id":           "sentinel-s2-l2a-cogs",
// 		"collection":   "sentinel-s2-l2a-cogs-testing",
// 		"stac_version": "1.0.0",
// 	})
// 	responseBody := bytes.NewBuffer(postBody)

// 	resp, err := http.Post("http://localhost:6002/collections", "application/json", responseBody)

// 	if err != nil {
// 		log.Fatalf("An Error Occured %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	sb := string(body)
// 	log.Printf(sb)
// }
