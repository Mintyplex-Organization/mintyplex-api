package routes

import (
	"mintyplex-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App) {
	route := app.Group("/api/v1/product")

	route.Post("/:id", controllers.AddProduct)
	route.Get("/", controllers.AllProducts)
	route.Get("/:id", controllers.OneProduct)
	route.Put("/:id/:uid", controllers.UpdateProduct)
	// route.Delete("/:id", controllers.DeleteProduct)


}
