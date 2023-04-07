package routes

import (
	"github.com/jonhealy1/goapi-stac/controllers"

	"github.com/gofiber/fiber/v2"
)

func ESItemRoute(app *fiber.App) {
	app.Post("/es/collections/:collectionId/items", controllers.ESCreateItem)
	app.Get("/es/collections/:collectionId/items/:itemId", controllers.ESGetItem)
	// app.Get("/es/collections/:collectionId/items", controllers.ESGetItemCollection)
	app.Put("/es/collections/:collectionId/items/:itemId", controllers.ESUpdateItem)
	app.Delete("/es/collections/:collectionId/items/:itemId", controllers.ESDeleteItem)
}
