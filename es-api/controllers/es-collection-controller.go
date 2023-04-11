package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jonhealy1/goapi-stac/es-api/database"
	"github.com/olivere/elastic"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"

	"github.com/jonhealy1/goapi-stac/es-api/models"
)

func CreateESCollection(c *fiber.Ctx) error {
	stac_collection := new(models.StacCollection)
	err := c.BodyParser(&stac_collection)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	now := time.Now()
	collection := models.Collection{
		Data:      models.JSONB{(&stac_collection)},
		Id:        stac_collection.Id,
		CreatedAt: &now,
	}
	validator := validator.New()
	err = validator.Struct(collection)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
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
		c.Status(http.StatusConflict).JSON(
			&fiber.Map{"message": fmt.Sprintf("Collection %s already exists", collection.Id)})
		return err
	}

	doc, err := json.Marshal(collection)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not marshal collection"})
		return err
	}

	resp, err := database.ES.Client.Index().
		Index(indexName).
		Id(collection.Id).
		BodyString(string(doc)).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not index collection"})
		return err
	}

	c.Status(http.StatusCreated).JSON(&fiber.Map{
		"message":         "success",
		"id":              resp.Id,
		"stac_collection": collection.Data[0],
	})
	return nil
}

func GetESCollection(c *fiber.Ctx) error {
	collectionID := c.Params("collectionId")

	indexName := "collections"

	// Retrieve the collection document from Elasticsearch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := database.ES.Client.Get().
		Index(indexName).
		Id(collectionID).
		Do(ctx)

	if err != nil {
		if elastic.IsNotFound(err) {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": fmt.Sprintf("%s does not exist", collectionID)})
		}
		if elasticErr, ok := err.(*elastic.Error); ok {
			return c.Status(http.StatusInternalServerError).JSON(
				&fiber.Map{"message": fmt.Sprintf("Error retrieving the collection: %v, Status: %v", err, elasticErr.Status)})
		}
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprintf("Error retrieving the collection: %v", err)})
	}

	// Unmarshal the Elasticsearch document source into a models.Collection
	var collection models.Collection
	if err := json.Unmarshal(resp.Source, &collection); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Error unmarshalling the collection"})
	}

	// Return the stac_collection JSON
	return c.JSON(collection.Data[0])
}

func GetESCollections(c *fiber.Ctx) error {
	indexName := "collections"

	// Retrieve all collection documents from Elasticsearch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	searchResult, err := database.ES.Client.Search().
		Index(indexName).
		Size(1000). // Adjust this value based on the expected number of collections
		Sort("CreatedAt", true).
		Do(ctx)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Error retrieving the collections"})
	}

	// Unmarshal the Elasticsearch documents sources into a list of models.Collection
	collections := make([]models.Collection, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		var collection models.Collection
		if err := json.Unmarshal(hit.Source, &collection); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				&fiber.Map{"message": "Error unmarshalling the collection"})
		}
		collections = append(collections, collection)
	}

	// Extract and return the collection.Data from each collection
	collectionDataList := make([]interface{}, len(collections))
	for i, collection := range collections {
		collectionDataList[i] = collection.Data[0]
	}
	return c.JSON(collectionDataList)
}

func EditESCollection(c *fiber.Ctx) error {
	id := c.Params("collectionId")
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "missing id parameter"})
		return fmt.Errorf("missing id parameter")
	}

	stac_collection := new(models.StacCollection)
	err := c.BodyParser(&stac_collection)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	now := time.Now()
	collection := models.Collection{
		Data:      models.JSONB{(&stac_collection)},
		Id:        stac_collection.Id,
		UpdatedAt: &now,
	}
	validator := validator.New()
	err = validator.Struct(collection)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	// Update the collection document in Elasticsearch
	indexName := "collections"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.ES.Client.Get().
		Index(indexName).
		Id(id).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": fmt.Sprintf("Collection %s not found", id)})
		return err
	}

	doc, err := json.Marshal(collection)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not marshal collection"})
		return err
	}

	// Unmarshal JSON string back into a map
	var docMap map[string]interface{}
	err = json.Unmarshal(doc, &docMap)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not unmarshal collection"})
		return err
	}

	resp, err := database.ES.Client.Update().
		Index(indexName).
		Id(id).
		Doc(docMap).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update collection"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message":         "success",
		"id":              resp.Id,
		"stac_collection": collection.Data[0],
	})
	return nil
}

func DeleteESCollection(c *fiber.Ctx) error {
	id := c.Params("collectionId")
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "missing id parameter"})
		return fmt.Errorf("missing id parameter")
	}

	indexName := "collections"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the collection exists
	_, err := database.ES.Client.Get().
		Index(indexName).
		Id(id).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": fmt.Sprintf("Collection %s not found", id)})
		return err
	}

	// Delete the collection document from Elasticsearch
	resp, err := database.ES.Client.Delete().
		Index(indexName).
		Id(id).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete collection"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "success",
		"id":      resp.Id,
	})
	return nil
}
