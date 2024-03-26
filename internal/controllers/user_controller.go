package controllers

import (
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserRepository repository.UserRepository
}

func NewUserController(UserRepository repository.UserRepository) UserController {
	return UserController{UserRepository}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.ReqResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}
