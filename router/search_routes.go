package routes

import (
	"go-stac-api-postgres/controllers"

	"github.com/gofiber/fiber/v2"
)

func SearchRoute(app *fiber.App) {
	app.Post("/search", controllers.PostSearch)
	app.Get("/search", controllers.GetSearch)
}
