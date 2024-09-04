package main

import (
	"mintyplex-api/internal/database"
	"mintyplex-api/internal/middleware"
	"mintyplex-api/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// fmt.Println("this is the region: ", os.Getenv("AWS_S3_REGION"))

	app := fiber.New(fiber.Config{
		BodyLimit: 200 * 1024 * 1024, // this is the default limit of 4MB
	})

	middleware.CorsMiddleware(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Plexer SZN")
	})

	db := database.MongoClient()
	app.Use(middleware.IngestDb(db))

	routes.UserRoutes(app)
	routes.ProductRoutes(app)
	routes.Reserve(app)

	app.Listen(":8081")
}
