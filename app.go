package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"

	database "go-stac-api-postgres/database"
	router "go-stac-api-postgres/router"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(cache.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(logger.New())

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	router.CollectionRoute(app)
	router.ItemRoute(app)

	app.Listen(":6002")
}
