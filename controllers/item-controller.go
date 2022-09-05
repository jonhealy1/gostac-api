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

	err := c.BodyParser(&stac_item)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	item := models.Item{
		Data:       models.JSONB{(&stac_item)},
		Id:         stac_item.Id,
		Collection: stac_item.Collection,
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

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message":    "item has been successfully added",
		"id":         item.Id,
		"collection": item.Collection,
		"stac_item":  item.Data[0],
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
	// collecton_id := c.Params("collectionId")
	item_id := c.Params("itemId")
	item := &models.Item{}
	if item_id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := database.DB.Db.Where("id = ?", item_id).First(item).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get collection"})
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
