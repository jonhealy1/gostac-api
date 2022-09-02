package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zeimedee/go-postgres/database"

	"go-stac-api-postgres/router"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	setUpRoutes(app)

	app.Use(cors.New())

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	router.CollectionRoute(app)
	router.ItemRoute(app)

	log.Fatal(app.Listen(":6000"))
}
