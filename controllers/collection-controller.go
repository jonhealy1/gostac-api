package controllers

import (
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func Root(c *fiber.Ctx) error {
	links := []models.Link{
		{
			Rel:   "self",
			Type:  "application/json",
			Href:  "/",
			Title: "root catalog",
		},
		{
			Rel:   "children",
			Type:  "application/json",
			Href:  "/collections",
			Title: "stac child collections",
		},
	}

	rootCatalog := models.Root{
		Id:          "test-catalog",
		StacVersion: "1.0.0",
		Description: "test catalog for go-stac-api, please edit",
		Title:       "go-stac-api",
		Links:       links,
	}

	return c.Status(http.StatusOK).JSON(rootCatalog)
}

// GetCollection godoc
// @Summary Get a Collection
// @Description Get a collection by ID
// @Tags Collections
// @ID get-collection-by-id
// @Accept  json
// @Produce  json
// @Param collectionId path string true "Collection ID"
// @Router /collections/{collectionId} [get]
// @Success 200 {object} models.Collection
func GetCollection(c *fiber.Ctx) error {
	// collection := models.Collection{}
	// match := new(models.Collection)
	// if err := c.BodyParser(match); err != nil {
	// 	return c.Status(400).JSON(err.Error())
	// }
	// database.DB.Db.Where("Id = ?", match.Id).Find(&collection)
	// return c.Status(200).JSON(collection)
	id := c.Params("collectionId")
	collection := &models.Collection{}
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := database.DB.Db.Where("id = ?", id).First(collection).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get collection"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "collection retrieved successfully",
		"data":    collection,
	})
	return nil
}

// CreateCollection godoc
// @Summary Create a STAC collection
// @Description Create a collection with a unique ID
// @Tags Collections
// @ID post-collection
// @Accept  json
// @Produce  json
// @Param collection body models.Collection true "STAC Collection json"
// @Router /collections [post]
func CreateCollection(c *fiber.Ctx) error {
	collection := new(models.Collection)
	err := c.BodyParser(&collection)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	validator := validator.New()
	err = validator.Struct(collection)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": err},
		)
		return err
	}

	err = database.DB.Db.Create(&collection).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create collection"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "collection has been successfully added",
		"data":    collection.Data[0],
	})
	return nil

}

// GetCollections godoc
// @Summary Get all Collections
// @Description Get all Collections
// @Tags Collections
// @ID get-all-collections
// @Accept  json
// @Produce  json
// @Router /collections [get]
// @Success 200 {object} []models.Collection
func GetCollections(c *fiber.Ctx) error {
	collections := []models.Collection{}
	err := database.DB.Db.Find(&collections).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get collections"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "collections received successfully",
		"data":    collections,
	})
	return nil
}
