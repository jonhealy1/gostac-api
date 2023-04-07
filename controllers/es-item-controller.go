package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jonhealy1/goapi-stac/database"

	"github.com/go-playground/validator"

	"github.com/gofiber/fiber/v2"

	"github.com/jonhealy1/goapi-stac/models"
)

func checkCollectionExists(collectionId string) (bool, error) {
	indexName := "collections"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if the collection exists
	exists, err := database.ES.Client.Exists().
		Index(indexName).
		Id(collectionId).
		Do(ctx)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func ESItemExists(itemId string) (bool, error) {
	indexName := "items"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := database.ES.Client.Exists().
		Index(indexName).
		Id(itemId).
		Do(ctx)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func ESCreateItem(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collectionId := c.Params("collectionId")
	if collectionId == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "collection id cannot be empty"})
		return fmt.Errorf("missing collectionId parameter")
	}

	exists, err := checkCollectionExists(collectionId)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprintf("Error checking collection %s: %v", collectionId, err)})
		return err
	}

	if !exists {
		c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": fmt.Sprintf("Collection %s not found", collectionId)})
		return fmt.Errorf("collection %s not found", collectionId)
	}

	stac_item := new(models.StacItem)
	err = c.BodyParser(&stac_item)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	itemId := stac_item.Id
	if itemId == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "missing itemId in item"})
		return fmt.Errorf("missing itemId in item")
	}

	validator := validator.New()
	err = validator.Struct(stac_item)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	indexName := "items"

	// Check if the item already exists
	_, err = database.ES.Client.Get().
		Index(indexName).
		Id(itemId).
		Do(ctx)

	if err == nil {
		c.Status(http.StatusConflict).JSON(
			&fiber.Map{"message": fmt.Sprintf("Item %s already exists", itemId)})
		return err
	}

	doc, err := json.Marshal(stac_item)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not marshal item"})
		return err
	}

	resp, err := database.ES.Client.Index().
		Index(indexName).
		Id(itemId).
		BodyString(string(doc)).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not index item"})
		return err
	}

	c.Status(http.StatusCreated).JSON(&fiber.Map{
		"message":    "success",
		"id":         resp.Id,
		"collection": collectionId,
		"stac_item":  stac_item,
	})
	return nil
}

func ESDeleteItem(c *fiber.Ctx) error {
	itemId := c.Params("itemId")
	if itemId == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "item id cannot be empty"})
		return fmt.Errorf("missing itemId parameter")
	}

	exists, err := ESItemExists(itemId)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not check if item exists"})
		return err
	}

	if !exists {
		c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": fmt.Sprintf("Item %s not found", itemId)})
		return c.Status(http.StatusNotFound).SendString(fmt.Sprintf("Item %s not found", itemId))
	}

	// Proceed with the deletion if the item exists
	indexName := "items"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.ES.Client.Delete().
		Index(indexName).
		Id(itemId).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not delete item"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": fmt.Sprintf("Item %s deleted successfully", itemId),
	})
	return nil
}

func ESUpdateItem(c *fiber.Ctx) error {
	collectionId := c.Params("collectionId")
	itemId := c.Params("itemId")

	if collectionId == "" || itemId == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "collection id and item id cannot be empty"})
		return fmt.Errorf("missing collectionId or itemId parameter")
	}

	exists, err := ESIItemExists(itemId)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "error checking item existence"})
		return err
	}

	if !exists {
		c.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": fmt.Sprintf("Item %s not found", itemId)})
		return fmt.Errorf("item not found")
	}

	updatedStacItem := new(models.StacItem)
	err = c.BodyParser(&updatedStacItem)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	validator := validator.New()
	err = validator.Struct(updatedStacItem)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	indexName := "items"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := json.Marshal(updatedStacItem)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not marshal item"})
		return err
	}

	_, err = database.ES.Client.Update().
		Index(indexName).
		Id(itemId).
		Doc(string(doc)).
		DocAsUpsert(true).
		Do(ctx)

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update item"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message":    "success",
		"id":         itemId,
		"collection": collectionId,
		"stac_item":  updatedStacItem,
	})
	return nil
}
