package controllers

import (
	"fmt"
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// PostSearch godoc
// @Summary POST Search request
// @Description Search for STAC items via the Search endpoint
// @Tags Search
// @ID post-search
// @Accept  json
// @Produce  json
// @Param search body models.Search true "Search body json"
// @Router /search [post]
func PostSearch(c *fiber.Ctx) error {
	var search models.Search
	var items []models.Item

	if err := c.BodyParser(&search); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Status":  http.StatusBadRequest,
			"Message": "error",
			"Data":    err.Error(),
		})
	}
	limit := 100
	if search.Limit > 0 {
		limit = search.Limit
	}

	tx1 := database.DB.Db.Limit(limit)
	tx2 := database.DB.Db.Limit(limit)
	if len(search.Collections) > 0 {
		tx1 = database.DB.Db.Limit(limit).Where("collection IN ?", search.Collections)
		tx2 = tx1.Limit(limit)
		fmt.Println(tx1)
	}

	if len(search.Ids) > 0 {
		tx2 = tx1.Limit(limit).Where("id IN ?", search.Ids)
		fmt.Println(tx1)
	}

	err := tx2.Find(&items).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get items"})
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
		"message":  "item collection retrieved successfully",
		"context":  context,
		"type":     "FeatureCollection",
		"features": stac_items,
	})

	return nil
}
