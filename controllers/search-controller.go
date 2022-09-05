package controllers

import "github.com/gofiber/fiber/v2"

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
	return nil
}
