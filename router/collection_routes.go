package routes

import (
	"github.com/jonhealy1/goapi-stac/controllers"

	"github.com/gofiber/fiber/v2"
)

func CollectionRoute(app *fiber.App) {
	app.Get("/", controllers.Root)
	app.Get("/conformance", controllers.Conformance)
	app.Post("/collections", controllers.CreateCollection)
	app.Get("/collections/:collectionId", controllers.GetCollection)
	app.Put("/collections/:collectionId", controllers.EditCollection)
	app.Delete("/collections/:collectionId", controllers.DeleteCollection)
	app.Get("/collections", controllers.GetCollections)
}
