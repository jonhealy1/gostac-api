package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	controllers "github.com/jonhealy1/goapi-stac/es-api/controllers"
	database "github.com/jonhealy1/goapi-stac/es-api/database"
	router "github.com/jonhealy1/goapi-stac/es-api/router"
)

func main() {
	app := Setup()

	value, exists := os.LookupEnv("API_PORT")
	api_port := 6003
	if exists {
		api_port, _ = strconv.Atoi(value)
	}

	// Listen on api port
	log.Fatal(app.Listen(fmt.Sprintf(":%d", api_port)))
}

func Setup() *fiber.App {
	// connect to database: elastic search
	database.ConnectES()

	// create new fiber app
	app := fiber.New()

	// register middleware
	app.Use(cors.New())
	app.Use(compress.New())
	//app.Use(cache.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(limiter.New(limiter.Config{
		Max: 1000,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(&fiber.Map{
				"status":  "fail",
				"message": "You have requested too many in a single time-frame! Please wait another minute!",
			})
		},
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	// app.Use(cache.New(cache.Config{
	// 	Next: func(c *fiber.Ctx) bool {
	// 		return c.Query("refresh") == "true"
	// 	},
	// 	Expiration:   30 * time.Minute,
	// 	CacheControl: true,
	// }))

	router.ESCollectionRoute(app)
	router.ESItemRoute(app)

	app.All("*", func(c *fiber.Ctx) error {
		errorMessage := fmt.Sprintf("Route '%s' does not exist in this API!", c.OriginalURL())

		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"status":  "fail",
			"message": errorMessage,
		})
	})

	// Start Kafka consumers
	go database.StartConsumer("new-postgres-collection", controllers.CreateCollectionFromMessage)
	go database.StartConsumer("update-postgres-collection", controllers.UpdateCollectionFromMessage)

	return app
}
