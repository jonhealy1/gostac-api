package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go-stac-api-postgres/database"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"

	"go-stac-api-postgres/models"
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

	// Define the mapping for the index
	mapping := `{
		"mappings": {
			"properties": {
			"data": {
				"properties": {
				"extent": {
					"properties": {
					"temporal": {
						"properties": {
						"interval": {
							"type": "text"
						}
						}
					}
					}
				}
				}
			}
			}
		}
	}`

	// Create Elasticsearch index if it doesn't exist
	indexName := "collections"

	// Delete index - just used for testing, comment out in production
	// _, _ = database.ES.Client.DeleteIndex(indexName).Do(context.Background())

	exists, err := database.ES.Client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not contact Elasticsearch"})
		return err
	}
	if !exists {
		_, err := database.ES.Client.CreateIndex(indexName).BodyString(mapping).Do(context.Background())
		if err != nil {
			c.Status(http.StatusInternalServerError).JSON(
				&fiber.Map{"message": "could not create Elasticsearch index"})
			return err
		}
	}

	// Index the collection document in Elasticsearch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
