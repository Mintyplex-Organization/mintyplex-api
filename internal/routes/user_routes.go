package routes

import (
	"mintyplex-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	route := app.Group("/api/v1/user")

	// <--- user profile --->
	route.Post("/profile/", controllers.DoTier1)
	route.Get("/profile/:id", controllers.UserProfile)
	route.Put("/profile/:id", controllers.UpdateUserProfile)
	route.Get("/users", controllers.GetUsers)

	// <--- avatar routes --->
	// route.Put("/avatar/:id", controllers.UpdateUserAvatar)
	// route.Get("/avatar/:id", controllers.GetAvatarById)

}
