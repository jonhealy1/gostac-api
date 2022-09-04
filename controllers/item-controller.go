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
	item := new(models.Item)
	err := c.BodyParser(&item)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
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
		"message": "item has been successfully added",
		"data":    item.Data[0],
	})
	return nil
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var item models.Item
	// defer cancel()

	// //validate the request body
	// if err := c.BodyParser(&item); err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.ItemResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	// }

	// //use the validator library to validate required fields
	// if validationErr := validate_item.Struct(&item); validationErr != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.ItemResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	// }

	// newItem := models.Item{
	// 	Id:             item.Id,
	// 	Type:           item.Type,
	// 	StacVersion:    item.StacVersion,
	// 	Collection:     item.Collection,
	// 	StacExtensions: item.StacExtensions,
	// 	Bbox:           item.Bbox,
	// 	Geometry:       item.Geometry,
	// 	Properties:     item.Properties,
	// 	Assets:         item.Assets,
	// 	Links:          item.Links,
	// }

	// result, err := stacItem.InsertOne(ctx, newItem)
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(responses.ItemResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	// }

	// return c.Status(http.StatusCreated).JSON(responses.ItemResponse{Status: http.StatusCreated, Message: "success", Data: result})
}
