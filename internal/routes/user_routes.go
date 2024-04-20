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
	// route.Put("/profile", controllers.EditUser)

	// <--- avatar routes --->
	route.Post("/avatar/:id", controllers.UploadUserAvatar)
	route.Get("/avatar/:id", controllers.GetAvatarById)
	// route.Get("/avatar", controllers.GetUserAvatar)
	
	route.Delete("/avatar", middleware.Auth(), controllers.DeleteUserAvatar)

}
