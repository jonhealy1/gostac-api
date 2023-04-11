package routes

import (
	"github.com/jonhealy1/goapi-stac/es-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func ESCollectionRoute(app *fiber.App) {
	app.Post("/collections", controllers.CreateESCollection)
	app.Get("/collections/:collectionId", controllers.GetESCollection)
	app.Put("/collections/:collectionId", controllers.EditESCollection)
	app.Delete("/collections/:collectionId", controllers.DeleteESCollection)
	app.Get("/collections", controllers.GetESCollections)
}
