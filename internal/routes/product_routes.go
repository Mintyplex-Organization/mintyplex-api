package routes

import (
	"mintyplex-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App) {
	route := app.Group("/api/v1/product")

	

	route.Put("/upload", controllers.ItemUpload)
	route.Post("/:id", controllers.AddProduct)
	route.Get("/", controllers.AllProducts)
	// route.Get("/", controllers.GetProducts)
	route.Get("/:id", controllers.OneProduct)
}
