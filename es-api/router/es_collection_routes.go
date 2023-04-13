package routes

import (
	"github.com/jonhealy1/goapi-stac/es-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func ESCollectionRoute(app *fiber.App) {
	app.Post("/collections", func(c *fiber.Ctx) error {
		return controllers.CreateESCollection(c, nil)
	})
	app.Get("/collections/:collectionId", controllers.GetESCollection)
	// app.Put("/collections/:collectionId", controllers.EditESCollection)
	app.Put("/collections/:collectionId", func(c *fiber.Ctx) error {
		return controllers.EditESCollection(c, nil)
	})
	app.Delete("/collections/:collectionId", controllers.DeleteESCollection)
	app.Get("/collections", controllers.GetESCollections)
}
