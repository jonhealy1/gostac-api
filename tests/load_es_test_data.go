package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jonhealy1/goapi-stac/database"
	"github.com/jonhealy1/goapi-stac/models"
)

func LoadEsCollection() {
	//database.ConnectES()
	jsonFile, err := os.Open("setup_data/test-collection.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var collection models.Collection
	json.Unmarshal(byteValue, &collection)

	put1, err := database.ES.Client.Index().
		Index("collections").
		Id(collection.Id).
		BodyJson(collection).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed collection %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

func LoadEsItems() {
	jsonFile, err := os.Open("setup_data/sentinel-s2-l2a-cogs_0_100.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var itemCollection models.ItemCollection
	json.Unmarshal(byteValue, &itemCollection)

	for _, item := range itemCollection.Features {
		put, err := database.ES.Client.Index().
			Index("items").
			Id(item.Id).
			BodyJson(item).
			Do(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Indexed item %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	}
}
