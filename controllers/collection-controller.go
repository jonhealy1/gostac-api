package controllers

import (
	"net/http"

	"github.com/jonhealy1/goapi-stac/database"
	"github.com/jonhealy1/goapi-stac/models"

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
		Description: "test catalog for goapistac, please edit",
		Title:       "goapistac",
		Links:       links,
	}

	return c.Status(http.StatusOK).JSON(rootCatalog)
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
	stac_collection := new(models.StacCollection)
	err := c.BodyParser(&stac_collection)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	collection := models.Collection{
		Data: models.JSONB{(&stac_collection)},
		Id:   stac_collection.Id,
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

	c.Status(http.StatusCreated).JSON(&fiber.Map{
		"message":         "success",
		"id":              collection.Id,
		"stac_collection": collection.Data[0],
	})
	return nil
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

	c.Status(http.StatusOK).JSON(collection.Data[0])
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
	results := database.DB.Db.Find(&collections)
	if results.Error != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get collections"})
		return results.Error
	}
	var data []interface{}
	for _, collection := range collections {
		data = append(data, collection.Data[0])
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"count":   results.RowsAffected,
		"results": data,
	})
	return nil
}

// DeleteCollection godoc
// @Summary Delete a Collection
// @Description Delete a collection by ID
// @Tags Collections
// @ID delete-collection-by-id
// @Accept  json
// @Produce  json
// @Param collectionId path string true "Collection ID"
// @Router /collections/{collectionId} [delete]
func DeleteCollection(c *fiber.Ctx) error {
	collection := &models.Collection{}

	id := c.Params("collectionId")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := database.DB.Db.Unscoped().Where("id = ?", id).Delete(&collection).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete collection",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "success",
	})
	return nil
}

// EditCollection godoc
// @Summary Edit a Collection
// @Description Edit a collection by ID
// @Tags Collections
// @ID edit-collection
// @Accept  json
// @Produce  json
// @Param collectionId path string true "Collection ID"
// @Param collection body models.Collection true "STAC Collection json"
// @Router /collections/{collectionId} [put]
// @Success 200 {object} models.Collection
func EditCollection(c *fiber.Ctx) error {

	id := c.Params("collectionId")
	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	collectionModel := &models.Collection{}
	collection := models.StacCollection{}

	err := c.BodyParser(&collection)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	updated := models.Collection{
		Id:   id,
		Data: models.JSONB{(&collection)},
	}

	err = database.DB.Db.Model(collectionModel).Where("id = ?", id).Updates(updated).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not update collection",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "success",
	})
	return nil
}

func Conformance(c *fiber.Ctx) error {
	conformsTo := []string{
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/core",
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/oas30",
		"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/geojson",
	}

	return c.Status(http.StatusOK).JSON(&fiber.Map{
		"conformsTo": conformsTo,
	})
}
