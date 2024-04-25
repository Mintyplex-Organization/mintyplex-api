package routes

import (
	"mintyplex-api/internal/controllers"
	"mintyplex-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	route := app.Group("/api/v1/user")

	// <--- user profile --->
	route.Post("/profile", controllers.DoTier1)
	route.Get("/profile/:id", controllers.UserProfile)
	route.Put("/profile/:id", controllers.UpdateUserProfile)

	// <--- avatar routes --->
	route.Post("/avatar/:id", controllers.UploadUserAvatar)
	// route.Get("/avatar", controllers.UpdateUserAvatar)
	route.Get("/avatar/:id", controllers.GetAvatarById)
	route.Delete("/avatar/:id", controllers.DeleteUserAvatar)
	// route.Get("/avatar", controllers.GetUserAvatar)
	
	route.Delete("/avatar", middleware.Auth(), controllers.DeleteUserAvatar)

}
