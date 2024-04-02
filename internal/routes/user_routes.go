package routes

import (
	"mintyplex-api/internal/controllers"
	"mintyplex-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	route := app.Group("/api/v1/user")

	route.Post("/profile", controllers.DoTier1)	
	route.Get("/profile/:id", controllers.UserProfile)
	// route.Put("/profile", controllers.EditUser)
	route.Post("/avatar", middleware.Auth(), controllers.UploadUserAvatar)
	route.Get("/avatar", middleware.Auth(), controllers.GetUserAvatar)
	route.Get("/avatar/:id", controllers.GetAvatarById)
	route.Delete("/avatar", middleware.Auth(), controllers.DeleteUserAvatar)
	
}
