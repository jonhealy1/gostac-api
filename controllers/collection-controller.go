package controllers

import (
	"fmt"
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
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
	collection := models.Collection{}
	match := new(models.Collection)
	if err := c.BodyParser(match); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.DB.Db.Where("Id = ?", match.Id).Find(&collection)
	return c.Status(200).JSON(collection)
	// return c.Status(http.StatusOK)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// collectionId := c.Params("collectionId")
	// var collection models.Collection
	// defer cancel()

	// err := stacCollection.FindOne(ctx, bson.M{"id": collectionId}).Decode(&collection)
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(responses.CollectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	// }

	// return c.Status(http.StatusOK).JSON(responses.CollectionResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Not Implemented"}})
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
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var collection models.Collection
	// defer cancel()

	collection := new(models.Collection)
	if err := c.BodyParser(collection); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// type newCollection struct {
	// 	Data interface{}
	// }

	// insert := newCollection{Data: collection.Data[0]}

	database.DB.Db.Create(&collection)

	return c.Status(200).JSON(collection.Data[0])

	// //validate the request body
	// if err := c.BodyParser(&collection); err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.CollectionResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	// }

	// //use the validator library to validate required fields
	// if validationErr := validate.Struct(&collection); validationErr != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.CollectionResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	// }

	// newCollection := models.Collection{
	// 	Id:          collection.Id,
	// 	StacVersion: collection.StacVersion,
	// 	Description: collection.Description,
	// 	Title:       collection.Title,
	// 	Links:       collection.Links,
	// 	Extent:      collection.Extent,
	// 	Providers:   collection.Providers,
	// }

	// result, err := stacCollection.InsertOne(ctx, newCollection)
	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(responses.CollectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	// }

	// return c.Status(http.StatusCreated).JSON(responses.CollectionResponse{Status: http.StatusCreated, Message: "success", Data: result})
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
	database.DB.Db.Find(&collections)

	return c.Status(200).JSON(collections)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// var collections []models.Collection
	// defer cancel()

	// results, err := stacCollection.Find(ctx, bson.M{})

	// if err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(responses.CollectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	// }

	// //reading from the db in an optimal way
	// defer results.Close(ctx)
	// for results.Next(ctx) {
	// 	var singleCollection models.Collection
	// 	if err = results.Decode(&singleCollection); err != nil {
	// 		return c.Status(http.StatusInternalServerError).JSON(responses.CollectionResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	// 	}

	// 	collections = append(collections, singleCollection)
	// }

	// return c.Status(http.StatusOK).JSON(
	// 	responses.CollectionResponse{Status: http.StatusOK, Message: "success", Data: collections},
	// )
}
