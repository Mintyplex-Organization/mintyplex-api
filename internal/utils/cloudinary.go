package utils

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	config, _ := LoadConfig(".")

	cldSecret := config.CloudinarySecret
	cldName := config.CloudName
	cldKey := config.CloudinaryKey

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}
	return cld, nil

}
