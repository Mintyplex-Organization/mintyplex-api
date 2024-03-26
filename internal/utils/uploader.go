package utils

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file multipart.File, filepath string) (string, error) {
	ctx := context.Background()
	cld, err := SetupCloudinary()
	if err != nil {
		return "", err
	}

	uploadParams := uploader.UploadParams{
		PublicID: filepath,
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}
	imageUrl := result.SecureURL
	return imageUrl, nil
}
