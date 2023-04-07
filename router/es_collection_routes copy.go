package routes

import (
	"github.com/jonhealy1/goapi-stac/controllers"

	"github.com/gofiber/fiber/v2"
)

func ESCollectionRoute(app *fiber.App) {
	app.Post("/es/collections", controllers.CreateESCollection)
	app.Get("/es/collections/:collectionId", controllers.GetESCollection)
	app.Put("/es/collections/:collectionId", controllers.EditESCollection)
	//app.Delete("/es/collections/:collectionId", controllers.DeleteESCollection)
	app.Get("/es/collections", controllers.GetESCollections)
}
