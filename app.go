package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	database "go-stac-api-postgres/database"
	router "go-stac-api-postgres/router"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	router.CollectionRoute(app)
	router.ItemRoute(app)

	app.Listen(":6002")
}
