package controllers

import (
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
	limit := 100
	if search.Limit > 0 {
		limit = search.Limit
	}

	if err := c.BodyParser(&search); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Status":  http.StatusBadRequest,
			"Message": "error",
			"Data":    err.Error(),
		})
	}

	if len(search.Collections) > 0 {
		err := database.DB.Db.Where("collection = ?", search.Collections).Find(&items).Error

		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "could not get collections"})
			return err
		}
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
