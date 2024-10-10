package routes

import (
	"mintyplex-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

type SessionRouteController struct{
	authController controllers.AuthController
}

func NewSessionRouteController(authController controllers.AuthController) SessionRouteController{
	return SessionRouteController{authController}
}

func (rc *SessionRouteController) SessionRoute(app fiber.Router){
	router := app.Group("/sessions/oauth")

	router.Get("/google", rc.authController.GoogleAuth)
}