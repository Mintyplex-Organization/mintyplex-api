package routes

import (
	"mintyplex-api/internal/controllers"
	"mintyplex-api/internal/middleware"
	"mintyplex-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup, UserRepository repository.UserRepository) {

	router := rg.Group("users")
	router.Use(middleware.DeserializeUser(UserRepository))
	router.GET("/myprofile", uc.userController.GetMe)
}
