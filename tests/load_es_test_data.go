package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/jonhealy1/goapi-stac/database"
	"github.com/jonhealy1/goapi-stac/models"
)

func LoadEsCollection() error {
	jsonFile, err := os.Open("setup_data/test-collection.json")
	if err != nil {
		return err
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var stacCollection models.StacCollection
	json.Unmarshal(byteValue, &stacCollection)

	now := time.Now()
	collection := models.Collection{
		Data:      models.JSONB{(&stacCollection)},
		Id:        stacCollection.Id,
		CreatedAt: &now,
	}
	validator := validator.New()
	err = validator.Struct(collection)

	if err != nil {
		return err
	}

	indexName := "collections"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.ES.Client.Get().
		Index(indexName).
		Id(collection.Id).
		Do(ctx)

	if err == nil {
		return fmt.Errorf("Collection %s already exists", collection.Id)
	}

	doc, err := json.Marshal(collection)
	if err != nil {
		return fmt.Errorf("Could not marshal collection: %v", err)
	}

	_, err = database.ES.Client.Index().
		Index(indexName).
		Id(collection.Id).
		BodyString(string(doc)).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("Could not index collection: %v", err)
	}

	return nil
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
