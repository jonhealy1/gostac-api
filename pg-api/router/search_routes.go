package routes

import (
	"github.com/jonhealy1/goapi-stac/pg-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func SearchRoute(app *fiber.App) {
	app.Post("/search", controllers.PostSearch)
	app.Get("/search", controllers.GetSearch)
}
