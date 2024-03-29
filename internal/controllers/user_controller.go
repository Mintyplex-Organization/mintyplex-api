package controllers

import (
	"fmt"
	"mime/multipart"
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/repository"
	"mintyplex-api/internal/utils"
	"net/http"
	"strconv"

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

func (uc *UserController) AddImage(ctx *gin.Context){
	id := ctx.Params.ByName("id")
	imagename, ok := ctx.Get("filepath")
	if !ok{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "imagename not found"})
	}

	image, ok := ctx.Get("file")
	if !ok{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "image not found"})
		return
	}

	imageUrl, err := utils.UploadToCloudinary(image.(multipart.File), imagename.(string))
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user *models.ReqResponse
	userId, _ := strconv.Atoi(id)
	fmt.Println(user, imageUrl, userId)

}