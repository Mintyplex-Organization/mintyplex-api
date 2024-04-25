package routes

import (
	"mintyplex-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func Reserve(app *fiber.App){
	route := app.Group("/api/v1/reserve")

	route.Post("/", controllers.ReserveUsername)
}