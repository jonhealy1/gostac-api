package controllers

import (
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// CreateItem godoc
// @Summary Create a STAC item
// @Description Create an item with an ID
// @Tags Items
// @ID post-item
// @Accept  json
// @Produce  json
// @Param collectionId path string true "Collection ID"
// @Param item body models.Item true "STAC Item json"
// @Router /collections/{collectionId}/items [post]
func CreateItem(c *fiber.Ctx) error {
	stac_item := new(models.StacItem)

	collection_id := c.Params("collectionId")
	if collection_id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "collection id cannot be empty",
		})
		return nil
	}

	err := c.BodyParser(&stac_item)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	item := models.Item{
		Data:       models.JSONB{(&stac_item)},
		Id:         stac_item.Id,
		Collection: collection_id,
	}

	validator := validator.New()
	err = validator.Struct(item)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	err = database.DB.Db.Create(&item).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create item"})
		return err
	}

	c.Status(http.StatusCreated).JSON(&fiber.Map{
		"message":    "success",
		"id":         item.Id,
		"collection": item.Collection,
		"stac_item":  item.Data[0],
	})
	return nil
}

// DeleteItem godoc
// @Summary Delete an Item
// @Description Delete an Item by ID is a specified collection
// @Tags Items
// @ID delete-item-by-id
// @Accept  json
// @Produce  json
// @Param itemId path string true "Item ID"
// @Param collectionId path string true "Collection ID"
// @Router /collections/{collectionId}/items/{itemId} [delete]
func DeleteItem(c *fiber.Ctx) error {
	item := &models.Item{}

	id := c.Params("itemId")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	collection_id := c.Params("collectionId")
	if collection_id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "collection id cannot be empty",
		})
		return nil
	}

	err := database.DB.Db.Unscoped().Where("id = ? AND collection = ?", id, collection_id).Delete(&item).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete item",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "item has been successfully deleted",
	})
	return nil
}

// EditItem godoc
// @Summary Edit an Item
// @Description Edit a stac item by ID
// @Tags Collections
// @ID edit-item
// @Accept  json
// @Produce  json
// @Param collectionId path string true "Collection ID"
// @Param itemId path string true "Item ID"
// @Param item body models.Item true "STAC Collection json"
// @Router /collections/{collectionId}/items/{itemId} [put]
// @Success 200 {object} models.Item
func EditItem(c *fiber.Ctx) error {
	id := c.Params("itemId")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	collection_id := c.Params("collectionId")
	if collection_id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "collection id cannot be empty",
		})
		return nil
	}

	itemModel := &models.Item{}
	item := models.StacItem{}

	err := c.BodyParser(&item)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	updated := models.Item{
		Id:         id,
		Collection: collection_id,
		Data:       models.JSONB{(&item)},
	}

	err = database.DB.Db.Model(itemModel).Where("id = ? AND collection = ?", id, collection_id).Updates(updated).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not update item",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "item has been successfully updated",
	})
	return nil
}

// GetItem godoc
// @Summary Get an item
// @Description Get an item by its ID
// @Tags Items
// @ID get-item-by-id
// @Accept  json
// @Produce  json
// @Param itemId path string true "Item ID"
// @Param collectionId path string true "Collection ID"
// @Router /collections/{collectionId}/items/{itemId} [get]
// @Success 200 {object} models.Item
func GetItem(c *fiber.Ctx) error {
	item := &models.Item{}

	item_id := c.Params("itemId")
	if item_id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	collection_id := c.Params("collectionId")
	if collection_id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "collection id cannot be empty",
		})
		return nil
	}

	err := database.DB.Db.Where("id = ? AND collection = ?", item_id, collection_id).First(item).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get item"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "item retrieved successfully",
		"id":      item.Id,
		// "collection": collecton_id,
		"stac_item": item.Data[0],
	})
	return nil
}

// GetItemCollection godoc
// @Summary Get all Items from a Collection
// @Description Get all Items with a Collection ID
// @Tags ItemCollection
// @ID get-item-collection
// @Accept  json
// @Produce  json
// @Param collectionId path string true "Collection ID"
// @Router /collections/{collectionId}/items [get]
// @Success 200 {object} models.ItemCollection
func GetItemCollection(c *fiber.Ctx) error {
	var items []models.Item
	collection_id := c.Params("collectionId")

	if collection_id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	limit := 100

	err := database.DB.Db.Where("collection = ?", collection_id).Find(&items).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get collections"})
		return err
	}

	context := models.Context{
		Returned: len(items),
		Limit:    limit,
	}

	var stac_items []interface{}
	for _, a_item := range items {
		stac_items = append(stac_items, a_item.Data)
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message":     "item collection retrieved successfully",
		"collection_": collection_id,
		"context":     context,
		"type":        "FeatureCollection",
		"features":    stac_items,
	})

	return nil
}
