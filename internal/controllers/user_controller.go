package controllers

import (
	"log"
	"mintyplex-api/internal/models"
	"mintyplex-api/internal/repository"
	"mintyplex-api/internal/utils"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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

func (uc *UserController) UploadAvatar(ctx *gin.Context) {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	cld, _ := cloudinary.NewFromURL(config.CloudinaryUri)

	imageName := ctx.PostForm("image")
	image, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	result, err := cld.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID: imageName,

	})
	if err != nil{
		ctx.String(http.StatusConflict, "failed to upload")
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "uploaded successfully",
		"secureURL": result.SecureURL,
		"publicURL": result.URL,
	})

}

// func UploadUserAvatar(ctx *gin.Context) {
// 	user := ctx.MustGet("id").(*models.UserResponse)
// 	database := ctx.MustGet("minty").(*mongo.Database)

// 	file, err := ctx.FormFile("avatar")
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
// 		return
// 	}

// 	if (file.Size) > 1024*1024 {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "file size too large, not >1mb"})
// 		return
// 	}

// 	fileExtension := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
// 	if fileExtension != ".jpg" && fileExtension != ".jpeg" && fileExtension != ".png" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "only .jpg. Invalid file type"})
// 		return

// 	}

// 	openFile, err := file.Open()
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "server error, let's try again"})
// 		return
// 	}

// 	content, err := io.ReadAll(openFile)
// 	if err != nil{
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "server error, let's try again in 3, 2, 1"})
// 		return
// 	}

// 	bucket, err :=

// }
