package routes

import (
	"github.com/jonhealy1/goapi-stac/es-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func ESItemRoute(app *fiber.App) {
	app.Post("/collections/:collectionId/items", controllers.ESCreateItem)
	app.Get("/collections/:collectionId/items/:itemId", controllers.ESGetItem)
	app.Get("/collections/:collectionId/items", controllers.ESGetItemCollection)
	app.Put("/collections/:collectionId/items/:itemId", controllers.ESUpdateItem)
	app.Delete("/collections/:collectionId/items/:itemId", controllers.ESDeleteItem)
}
