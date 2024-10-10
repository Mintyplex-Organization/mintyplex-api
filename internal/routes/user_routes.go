package routes

import (
	"mintyplex-api/internal/controllers"
	"mintyplex-api/internal/middleware"
	"mintyplex-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(app *fiber.App, userService services.UserService) {
	router := app.Group("/auth")

	// router.Post("/register", rc.authController.SignUpUser)
	router.Post("/register", rc.authController.SignUpUser)
	// router.Post("/login", rc.authController.SignInUser)
	// router.Get("/refresh", rc.authController.RefreshAccessToken)
	// router.Get("/logout", middleware.DeserializeUser(userService), rc.authController.LogoutUser)
}

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoutes(app *fiber.App, userService services.UserService) {
	router := app.Group("/api/v1/user")
	router.Use(middleware.DeserializeUser(userService))

	router.Get("/me", uc.userController.UserProfile)
}

// func UserRoutes(app *fiber.App) {
// 	route := app.Group("/api/v1/user")

// 	// <--- user profile --->
// 	route.Post("/profile/", controllers.DoTier1)
// 	route.Get("/profile/:id", controllers.UserProfile)
// 	route.Put("/profile/:id", controllers.UpdateUserProfile)
// 	route.Get("/users", controllers.GetUsers)

// 	// <--- avatar routes --->
// 	// route.Put("/avatar/:id", controllers.UpdateUserAvatar)
// 	// route.Get("/avatar/:id", controllers.GetAvatarById)

// }
