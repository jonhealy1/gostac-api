package routes

import (
	"go-stac-api-postgres/controllers"

	"github.com/gofiber/fiber/v2"
)

func ItemRoute(app *fiber.App) {
	app.Post("/collections/:collectionId/items", controllers.CreateItem)
	app.Get("/collections/:collectionId/items/:itemId", controllers.GetItem)
	app.Get("/collections/:collectionId/items", controllers.GetItemCollection)
	app.Put("/collections/:collectionId/items/:itemId", controllers.EditItem)
	app.Delete("/collections/:collectionId/items/:itemId", controllers.DeleteItem)
	// app.Post("/search", controllers.PostSearch)
}
