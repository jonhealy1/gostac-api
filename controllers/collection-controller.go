package controllers

import (
	"fmt"
	"go-stac-api-postgres/models"
	"go-stac-api-postgres/responses"
	"net/http"

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
	fmt.Println("Not Implemented")
	// return c.Status(http.StatusOK)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// collectionId := c.Params("collectionId")
	// var collection models.Collection
	// defer cancel()

	// err := stacCollection.FindOne(ctx, bson.M{"id": collectionId}).Decode(&collection)
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(responses.CollectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	// }

	return c.Status(http.StatusOK).JSON(responses.CollectionResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Not Implemented"}})
}
