package controllers

import (
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

	return c.Status(http.StatusOK).JSON(&fiber.Map{"data": rootCatalog})
}
